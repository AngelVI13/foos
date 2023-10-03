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
		log.Fatal(err)
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

func savePairsToFile(pairs map[string]string, name string) error {
	// delete existing pairs
	os.Remove(name)

	b, err := json.MarshalIndent(pairs, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(name, b, 0666)
}

func loadPairsFromFile(name string) (map[string]string, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var pairs map[string]string
	err = json.Unmarshal(b, &pairs)
	if err != nil {
		return nil, err
	}

	return pairs, nil
}

const (
	playersFile = "players.txt"
	pairsFile   = "pairs.json"
)

func main() {
	players, err := getPlayersFromFile(playersFile)
	if err != nil {
		log.Println(err)
	}
	if !(len(players) == 8 || len(players) == 12) {
		log.Fatalf("players should be 8 or 12 but got: %d: %v", len(players), players)
	}
	fmt.Println("=== Player List:", players)

	prevPairs, err := loadPairsFromFile(pairsFile)
	if err != nil {
		fmt.Println("No previous pairs found.")
		prevPairs = map[string]string{}
	}
	fmt.Println("== Previous Teams:", prevPairs)
	fmt.Println()
	pairs := map[string]string{}

	var player string
	for len(players) > 1 {
		player, players = playersPop(players)
		prevPartner, _ := prevPairs[player]
		fmt.Printf("= Selecting partner for %s\n", player)
		fmt.Printf("\tLast partner: %s\n", prevPartner)

		var choices []weightedrand.Choice[string, int]
		for i, p := range players {
			if p == prevPartner {
				// no probability to be selected
				choices = append(choices, weightedrand.NewChoice(p, 0))
			}
			choices = append(choices, weightedrand.NewChoice(p, i+1))
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
	fmt.Println("=====================\n\n")
	fmt.Println("Teams:", pairs)

	err = savePairsToFile(pairs, pairsFile)
	if err != nil {
		log.Fatal(err)
	}
}
