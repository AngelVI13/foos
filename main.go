package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mroth/weightedrand/v2"
)

func getPlayersFromFile(name string) ([]string, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var players []string
	txt := string(b)
	for _, line := range strings.Split(txt, "\n") {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "#") {
			continue
		}
		if l != "" {
			players = append(players, l)
		}
	}

	if len(players)%2 != 0 {
		return nil, fmt.Errorf("uneven number of players: %d: %v", len(players), players)
	}
	return players, nil
}

// playersPop remove first element of players and return it along with remaining elements
func playersPop(players []string) (string, []string) {
	if len(players) < 1 {
		log.Fatalf("can't pop from players list: size is already 0")
	}

	player := players[0]

	if len(players) == 1 {
		return player, nil
	}
	return player, players[1:]
}

func playersRemove(players []string, player string) []string {
	var out []string
	for _, p := range players {
		if p != player {
			out = append(out, p)
		}
	}

	return out
}

func savePairsToFile(pairs map[string]string, newPairsName, lastPairsName string) error {
	// delete existing pairs
	os.Rename(newPairsName, lastPairsName)

	b, err := json.MarshalIndent(pairs, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(newPairsName, b, 0o666)
}

func loadPairsFromFile(name string) map[string]string {
	b, err := os.ReadFile(name)
	if err != nil {
		return map[string]string{}
	}

	var pairs map[string]string
	err = json.Unmarshal(b, &pairs)
	if err != nil {
		return map[string]string{}
	}

	return pairs
}

func generatePairs(players []string, prevPairs1, prevPairs2 map[string]string) map[string]string {
	pairs := map[string]string{}

	var player string
	for len(players) > 1 {
		player, players = playersPop(players)
		prevPartner1 := prevPairs1[player]
		prevPartner2 := prevPairs2[player]
		fmt.Printf("= Selecting partner for %s\n", player)
		fmt.Printf("\tLast partners: %s %s\n", prevPartner1, prevPartner2)

		var choices []weightedrand.Choice[string, int]
		for i, p := range players {
			if p == prevPartner1 || p == prevPartner2 {
				// no probability to be selected
				choices = append(choices, weightedrand.NewChoice(p, 0))
			} else {
				choices = append(choices, weightedrand.NewChoice(p, i+1))
			}
		}

		fmt.Println("\tProbabilities:", choices)
		chooser, err := weightedrand.NewChooser(choices...)
		if err != nil {
			log.Fatal(err)
		}
		result := chooser.Pick()
		players = playersRemove(players, result)

		fmt.Println("\n>> ", player, result)
		fmt.Println()
		pairs[player] = result
	}
	return pairs
}

const (
	playersFile     = "players.txt"
	pairsFile1      = "teams1.json"
	pairsFile2      = "teams2.json"
	playersErrorTxt = `

    couldn't process players file.
    make sure you create a file 'players.txt' in the same directory as the script
    with each players name on a separate line starting from the strongest player
    (at the top of the file) and finishing with the least strong player on the last
    line.
    `
)

func main() {
	players, err := getPlayersFromFile(playersFile)
	if err != nil {
		log.Println(err)
		log.Fatal(playersErrorTxt)
	}
	if !(len(players) == 8 || len(players) == 12) {
		log.Fatalf("players should be 8 or 12 but got: %d: %v", len(players), players)
	}
	fmt.Println("=== Player List:", players)

	prevPairs1 := loadPairsFromFile(pairsFile1)
	prevPairs2 := loadPairsFromFile(pairsFile2)
	fmt.Println("== Previous Teams1:", prevPairs1)
	fmt.Println("== Previous Teams2:", prevPairs2)
	fmt.Println()

	pairs := generatePairs(players, prevPairs1, prevPairs2)

	fmt.Println("=====================\n\n")
	fmt.Println("Teams:", pairs)

	err = savePairsToFile(pairs, pairsFile1, pairsFile2)
	if err != nil {
		log.Fatal(err)
	}
}
