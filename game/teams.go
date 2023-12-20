package game

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/AngelVI13/foos/log"
	"github.com/mroth/weightedrand/v2"
)

type rawTeams [][2]string

func (t rawTeams) Partner(player string) string {
	for _, team := range t {
		for i, p := range team {
			if p == player {
				if i == 0 {
					return team[1]
				} else {
					return team[0]
				}
			}
		}
	}
	return ""
}

func saveTeamsToFile(teams rawTeams, newTeamsName, lastTeamsName string) error {
	// delete existing teams
	os.Rename(newTeamsName, lastTeamsName)

	b, err := json.MarshalIndent(teams, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(newTeamsName, b, 0o666)
}

func loadTeamsFromFile(name string) rawTeams {
	b, err := os.ReadFile(name)
	if err != nil {
		return rawTeams{}
	}

	var teams rawTeams
	err = json.Unmarshal(b, &teams)
	if err != nil {
		return rawTeams{}
	}

	return teams
}

func numChoices(choices []weightedrand.Choice[string, int]) int {
	choiceNum := 0
	for _, c := range choices {
		if c.Weight > 0 {
			choiceNum++
		}
	}
	return choiceNum
}

func generateJudgementDayTeams(
	players []string,
	prevTeams1, prevTeams2 rawTeams,
) [][2]string {
	teams := [][2]string{} // use this in order to keep teams order
	for len(teams) < len(players)/2 {
		teams = append(teams, [2]string{})
	}
	log.L.Info("", "teams", teams)
	currentTeamIndex := 0
	totalPlayers := len(players)

	var player string
	for len(players) > 1 {
		player, players = playersPop(players)
		prevPartner1 := prevTeams1.Partner(player)
		prevPartner2 := prevTeams2.Partner(player)
		log.L.Info(fmt.Sprintf("Selecting partner for %s\n", player))
		log.L.Info(fmt.Sprintf("Last partners: %s %s\n", prevPartner1, prevPartner2))

		var choices []weightedrand.Choice[string, int]
		for i, p := range players {
			if i+totalPlayers-len(players) < (totalPlayers/2)+currentTeamIndex {
				// in case of judgement day -> no top half players can be paired
				// with each other in a team
				choices = append(choices, weightedrand.NewChoice(p, 0))
			} else if (i == len(players)-1) && numChoices(choices) == 0 {
				// if we are at the last element and currently no choices are added
				// to list of players -> add the last player with 100%
				choices = append(choices, weightedrand.NewChoice(p, i+1))
			} else if p == prevPartner1 || p == prevPartner2 {
				// skip previous partners
				choices = append(choices, weightedrand.NewChoice(p, 0))
			} else {
				// for any other case add players with increasing weighted choice probability
				choices = append(choices, weightedrand.NewChoice(p, i+1))
			}
		}

		log.L.Info("", "Probabilities:", choices)
		chooser, err := weightedrand.NewChooser(choices...)
		if err != nil {
			panic(err)
		}
		result := chooser.Pick()
		players = playersRemove(players, result)

		teams[currentTeamIndex][0] = player
		teams[currentTeamIndex][1] = result
		currentTeamIndex++
	}
	return teams
}

func generateTeams(
	players []string,
	prevTeams1, prevTeams2 rawTeams,
) [][2]string {
	teams := [][2]string{} // use this in order to keep teams order
	for len(teams) < len(players)/2 {
		teams = append(teams, [2]string{})
	}
	log.L.Info("", "teams", teams)
	currentTeamIndex := 0

	var player string
	for len(players) > 1 {
		player, players = playersPop(players)
		prevPartner1 := prevTeams1.Partner(player)
		prevPartner2 := prevTeams2.Partner(player)
		log.L.Info(fmt.Sprintf("Selecting partner for %s\n", player))
		log.L.Info(fmt.Sprintf("Last partners: %s %s\n", prevPartner1, prevPartner2))

		var choices []weightedrand.Choice[string, int]
		for i, p := range players {
			if (i == len(players)-1) && numChoices(choices) == 0 {
				// if we are at the last element and currently no choices are added
				// to list of players -> add the last player with 100%
				choices = append(choices, weightedrand.NewChoice(p, i+1))
			} else if p == prevPartner1 || p == prevPartner2 {
				// no probability to be selected
				choices = append(choices, weightedrand.NewChoice(p, 0))
			} else {
				choices = append(choices, weightedrand.NewChoice(p, i+1))
			}
		}

		log.L.Info("", "Probabilities:", choices)
		chooser, err := weightedrand.NewChooser(choices...)
		if err != nil {
			panic(err)
		}
		result := chooser.Pick()
		players = playersRemove(players, result)

		teams[currentTeamIndex][0] = player
		teams[currentTeamIndex][1] = result
		currentTeamIndex++
	}
	return teams
}

const (
	playersFile     = "players.txt"
	TeamsFile1      = "teams1.json"
	TeamsFile2      = "teams2.json"
	playersErrorTxt = `

    couldn't process players file.
    make sure you create a file 'players.txt' in the same directory as the script
    with each players name on a separate line starting from the strongest player
    (at the top of the file) and finishing with the least strong player on the last
    line.
    `
)

func PlayersByStats(players []string, stats map[string]*Stats) []string {
	if len(stats) == 0 {
		return players
	}
	sort.Slice(players, func(i, j int) bool {
		player1Stat, player1Found := stats[players[i]]
		player2Stat, player2Found := stats[players[j]]
		if !player1Found && !player2Found {
			return true
		} else if player1Found && !player2Found {
			return true
		} else if !player1Found && player2Found {
			return false
		}

		return player1Stat.SuccessRate() > player2Stat.SuccessRate()
	})
	return players
}

func PlayersRankings(stats map[string]*Stats) map[string]int {
	var allPlayers []string

	for player := range stats {
		allPlayers = append(allPlayers, player)
	}

	sort.Slice(allPlayers, func(i, j int) bool {
		player1Stat, player1Found := stats[allPlayers[i]]
		player2Stat, player2Found := stats[allPlayers[j]]
		if !player1Found && !player2Found {
			return true
		} else if player1Found && !player2Found {
			return true
		} else if !player1Found && player2Found {
			return false
		}

		return player1Stat.SuccessRate() > player2Stat.SuccessRate()
	})

	rankings := map[string]int{}
	for i, player := range allPlayers {
		rankings[player] = i + 1
	}
	return rankings
}

// TeamsByStats Reorder teams so that 1st player does not play with 2nd player (based on overall stats)
func TeamsByStats(teams []Team) []Team {
	var newTeams []Team
	numMatches := len(teams) / 2
	i := 0
	for i < numMatches {
		newTeams = append(newTeams, teams[i])
		newTeams = append(newTeams, teams[i+numMatches])
		i++
	}

	return newTeams
}

// GenerateTeams Generates weighted random teams based on order of playes (best to worst)
func GenerateTeams(
	players []string,
	judgementDay bool,
	playersRankings map[string]int,
) ([]Team, error) {
	prevTeams1 := loadTeamsFromFile(TeamsFile1)
	prevTeams2 := loadTeamsFromFile(TeamsFile2)
	log.L.Info("", "Previous Teams1:", prevTeams1)
	log.L.Info("", "Previous Teams2:", prevTeams2)

	var teamsSlice [][2]string
	if judgementDay {
		teamsSlice = generateJudgementDayTeams(players, prevTeams1, prevTeams2)
	} else {
		teamsSlice = generateTeams(players, prevTeams1, prevTeams2)
	}
	log.L.Info("", "Teams:", teamsSlice)

	err := saveTeamsToFile(teamsSlice, TeamsFile1, TeamsFile2)
	if err != nil {
		return nil, err
	}

	var teams []Team
	lastRank := len(playersRankings) + 1
	for _, players := range teamsSlice {
		player1Rank, found := playersRankings[players[0]]
		if !found {
			player1Rank = lastRank
			lastRank++
		}
		player2Rank, found := playersRankings[players[1]]
		if !found {
			player2Rank = lastRank
			lastRank++
		}
		teams = append(teams, NewTeam(players[0], players[1], player1Rank, player2Rank))
	}

	return teams, nil
}
