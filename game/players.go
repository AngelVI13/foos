package game

import (
	"fmt"
	"log"
	"os"
	"strings"
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
