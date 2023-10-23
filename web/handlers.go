package web

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/AngelVI13/foos/game"
	"github.com/AngelVI13/foos/log"
	"github.com/AngelVI13/foos/routes"
	"github.com/AngelVI13/foos/views"
	"github.com/gin-gonic/gin"
)

func errorHandler(c *gin.Context, msg string) {
	c.HTML(http.StatusOK, "", views.Page(views.Error(msg)))
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "", views.Page(views.UsersForm()))
}

func parsePlayersInput(playersInput string) []string {
	playersSplit := strings.Split(playersInput, "\r\n")
	var players []string
	for _, p := range playersSplit {
		// remove X) preffix where X is a number
		// i.e. 1) Ugnius
		if strings.Count(p, ")") == 1 {
			index := strings.Index(p, ")")
			p = p[index+1:]
		}

		if strings.HasPrefix(p, "-") {
			p = strings.ReplaceAll(p, "-", "")
		}

		player := strings.TrimSpace(p)
		if player == "" {
			continue
		}
		players = append(players, player)
	}

	return players
}

func findDuplicates(values []string) []string {
	seen := map[string]bool{}
	var duplicates []string
	for _, v := range values {
		if _, found := seen[v]; found {
			duplicates = append(duplicates, v)
		} else {
			seen[v] = true
		}
	}
	return duplicates
}

func usersListHandler(c *gin.Context) {
	deletePrevTeams := c.PostForm("deletePrevTeams")
	if deletePrevTeams == "on" {
		log.L.Info("deleting teams")
		os.Remove(game.TeamsFile1)
		os.Remove(game.TeamsFile2)
	}

	resetSeasonStats := c.PostForm("resetSeasonStats")
	if resetSeasonStats == "on" {
		log.L.Info("deleting season stats")
		os.Remove(game.StatsFile)
	}

	playersRawInput := c.PostForm("playersListInput")
	players := parsePlayersInput(playersRawInput)

	if !(len(players) == 8 || len(players) == 12) {
		errorHandler(c, fmt.Sprintf("Expected 8 or 12 players but got %d", len(players)))
		return
	}

	duplicates := findDuplicates(players)
	if len(duplicates) > 0 {
		errorHandler(c, fmt.Sprintf("Found duplicate players: %v", duplicates))
		return
	}

	var err error
	state, err = NewGlobalState(players)
	if err != nil {
		errorHandler(c, fmt.Sprintf("Failed to generate teams: %v", err))
		return
	}

	c.Redirect(http.StatusFound, routes.TournamentTableUrl)
}

func tournamentBracketHandler(c *gin.Context) {
	var teamsCopy []*game.Team
	teamsCopy = append(teamsCopy, state.Rounds.Teams...)

	sort.Slice(teamsCopy, func(i, j int) bool {
		return teamsCopy[i].AllScores() > teamsCopy[j].AllScores()
	})

	var statsCopy []*game.Stats
	for _, s := range state.Stats {
		statsCopy = append(statsCopy, s)
	}

	sort.Slice(statsCopy, func(i, j int) bool {
		if statsCopy[i].Score == statsCopy[j].Score {
			return statsCopy[i].Won > statsCopy[j].Won
		}
		return statsCopy[i].Score > statsCopy[j].Score
	})

	c.HTML(
		http.StatusOK,
		"",
		views.Page(views.Rounds(state.Rounds, teamsCopy, statsCopy)),
	)
}

func getMatchInfoFromRequest(c *gin.Context) (*game.Match, *game.Team, int, error) {
	match_id := c.Param("match_id")
	team := c.Param("team")

	if match_id == "" || team == "" {
		return nil, nil, -1, fmt.Errorf("missing match_id or team: %s %s", match_id, team)
	}

	teamNum, err := strconv.Atoi(team)
	if err != nil {
		return nil, nil, -1, fmt.Errorf(
			"failed to convert team to int: %s %v",
			team,
			err,
		)
	}

	match := state.Rounds.FindMatchById(match_id)
	teamPtr := match.Teams()[teamNum-1]
	return match, teamPtr, teamNum, nil
}

func tournamentBracketMatchUpdateHandler(c *gin.Context) {
	match, teamPtr, teamNum, err := getMatchInfoFromRequest(c)
	if err != nil {
		errorHandler(c, err.Error())
		return
	}
	url := routes.MakeMatchShowUrl(match, teamNum)

	c.HTML(
		http.StatusOK,
		"",
		views.TeamRowUpdate(teamPtr.String(), url, teamPtr.Score()),
	)
}

func tournamentBracketMatchShowHandler(c *gin.Context) {
	match, teamPtr, teamNum, err := getMatchInfoFromRequest(c)
	if err != nil {
		errorHandler(c, err.Error())
		return
	}

	if c.Request.Method == http.MethodPost {
		scoreStr := c.PostForm("score")
		score, err := strconv.Atoi(scoreStr)
		if err != nil {
			errorHandler(c, fmt.Sprintf("Score is not an integer: %q: %v", scoreStr, err))
			return
		}
		teamPtr.SetScore(score)
	}

	editUrl := routes.MakeMatchUpdateUrl(match, teamNum)

	c.HTML(http.StatusOK, "", views.TeamRow(teamPtr, editUrl, state.Rounds.CurrentRound))
}

func tournamentBracketEndRoundHandler(c *gin.Context) {
	currentRound := state.Rounds.All[state.Rounds.CurrentRound]

	// Check if all matches ended i.e. one team has more points than the other
	endedMatches := 0
	for _, match := range currentRound.Matches {
		teams := match.Teams()
		team1 := teams[0]
		team2 := teams[1]

		if team1.Score() == team2.Score() {
			continue
		}

		endedMatches++
	}

	if endedMatches != len(currentRound.Matches) {
		errorHandler(c, fmt.Sprintf("%d matches not finished", endedMatches))
		return
	}

	// Update all time stats
	for _, match := range currentRound.Matches {
		stats := match.End()

		for player := range stats {
			if _, found := state.Stats[player]; found {
				state.Stats[player].Score += stats[player].Score
				state.Stats[player].Won += stats[player].Won
				state.Stats[player].Lost += stats[player].Lost
			} else {
				state.Stats[player] = stats[player]
			}
		}
	}
	err := game.SaveStats(state.Stats)
	log.L.Error("failed to update all time stats", "err", err)

	state.Rounds.NextRound()

	c.Redirect(http.StatusFound, routes.TournamentTableUrl)
}

func newTournamentHandler(c *gin.Context) {
	var err error
	state, err = NewGlobalState(state.Players)
	if err != nil {
		errorHandler(c, fmt.Sprintf("Failed to generate teams: %v", err))
		return
	}

	c.Redirect(http.StatusFound, routes.TournamentTableUrl)
}
