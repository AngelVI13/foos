package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/AngelVI13/foos/game"
	"github.com/AngelVI13/foos/views"
	"github.com/gin-gonic/gin"
)

type GlobalState struct {
	Players []string
	Teams   []game.Team
	Matches []game.Match
	Round   int
}

func (s *GlobalState) GenerateMatches() {
	var matches []game.Match
	if s.Round == 0 {
		for i := 2; i <= len(s.Teams); i += 2 {
			matches = append(
				matches,
				game.NewMatch(s.Teams[i-2], s.Teams[i-1]),
			)
		}
	} else {
		panic("not implemented")
	}

	s.Matches = matches
}

func NewEmptyGlobalState() GlobalState {
	return GlobalState{
		Players: []string{},
		Teams:   []game.Team{},
		Matches: []game.Match{},
		Round:   0,
	}
}

func NewGlobalState(players []string, teams []game.Team) GlobalState {
	return GlobalState{
		Players: players,
		Teams:   teams,
		Matches: []game.Match{},
		Round:   0,
	}
}

var state = NewEmptyGlobalState()

func errorHandler(c *gin.Context, msg string) {
	c.HTML(http.StatusOK, "", views.Page(views.Error(msg)))
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "", views.Page(views.UsersForm()))
}

func usersListHandler(c *gin.Context) {
	playersRawInput := c.PostForm("playersListInput")
	playersSplit := strings.Split(playersRawInput, "\r\n")
	var players []string
	for _, p := range playersSplit {
		player := strings.TrimSpace(p)
		if player == "" {
			continue
		}
		players = append(players, player)
	}

	if !(len(players) == 8 || len(players) == 12) {
		errorHandler(c, fmt.Sprintf("Expected 8 or 12 playes but got %d", len(players)))
		return
	}

	teams, err := game.GenerateTeams(players, log)
	if err != nil {
		errorHandler(c, fmt.Sprintf("Failed to generate teams: %v", err))
		return
	}
	log.Info("", "teams", teams)
	state = NewGlobalState(players, teams)

	c.Redirect(http.StatusFound, "/tournament_table")
}

func tournamentBracketHandler(c *gin.Context) {
	state.GenerateMatches()
	c.HTML(http.StatusOK, "", views.Page(views.Bracket(state.Matches)))
}

var log = slog.New(
	slog.NewTextHandler(os.Stdout, nil))

// slog.NewJSONHandler(os.Stdout, nil))

func main() {
	r := gin.Default()
	r.HTMLRender = &views.TemplRender{}

	r.StaticFS("/static", http.FS(views.EmbedFs))

	r.GET("/", indexHandler)
	r.POST("/users_list", usersListHandler)
	r.GET("/tournament_table", tournamentBracketHandler)

	log.Info("Running")

	r.Run(":5555")
}
