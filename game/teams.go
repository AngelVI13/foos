package game

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/mroth/weightedrand/v2"
)

func saveTeamsToFile(teams map[string]string, newTeamsName, lastTeamsName string) error {
	// delete existing teams
	os.Rename(newTeamsName, lastTeamsName)

	b, err := json.MarshalIndent(teams, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(newTeamsName, b, 0o666)
}

func loadTeamsFromFile(name string) map[string]string {
	b, err := os.ReadFile(name)
	if err != nil {
		return map[string]string{}
	}

	var teams map[string]string
	err = json.Unmarshal(b, &teams)
	if err != nil {
		return map[string]string{}
	}

	return teams
}

func generateTeams(
	players []string,
	prevTeams1, prevTeams2 map[string]string,
	log *slog.Logger,
) map[string]string {
	teams := map[string]string{}

	var player string
	for len(players) > 1 {
		player, players = playersPop(players)
		prevPartner1 := prevTeams1[player]
		prevPartner2 := prevTeams2[player]
		log.Info(fmt.Sprintf("Selecting partner for %s\n", player))
		log.Info(fmt.Sprintf("Last partners: %s %s\n", prevPartner1, prevPartner2))

		var choices []weightedrand.Choice[string, int]
		for i, p := range players {
			if p == prevPartner1 || p == prevPartner2 {
				// no probability to be selected
				choices = append(choices, weightedrand.NewChoice(p, 0))
			} else {
				choices = append(choices, weightedrand.NewChoice(p, i+1))
			}
		}

		log.Info("", "Probabilities:", choices)
		chooser, err := weightedrand.NewChooser(choices...)
		if err != nil {
			panic(err)
		}
		result := chooser.Pick()
		players = playersRemove(players, result)

		log.Info("", player, result)
		teams[player] = result
	}
	return teams
}

const (
	playersFile     = "players.txt"
	teamsFile1      = "teams1.json"
	teamsFile2      = "teams2.json"
	playersErrorTxt = `

    couldn't process players file.
    make sure you create a file 'players.txt' in the same directory as the script
    with each players name on a separate line starting from the strongest player
    (at the top of the file) and finishing with the least strong player on the last
    line.
    `
)

func GenerateTeams(players []string, log *slog.Logger) ([]Team, error) {
	prevTeams1 := loadTeamsFromFile(teamsFile1)
	prevTeams2 := loadTeamsFromFile(teamsFile2)
	log.Info("", "Previous Teams1:", prevTeams1)
	log.Info("", "Previous Teams2:", prevTeams2)

	teamsMap := generateTeams(players, prevTeams1, prevTeams2, log)
	log.Info("", "Teams:", teamsMap)

	err := saveTeamsToFile(teamsMap, teamsFile1, teamsFile2)
	if err != nil {
		return nil, err
	}

	var teams []Team
	for p1, p2 := range teamsMap {
		teams = append(teams, NewTeam(p1, p2))
	}

	return teams, nil
}
