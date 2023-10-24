package web

import (
	"github.com/AngelVI13/foos/game"
	"github.com/AngelVI13/foos/log"
)

type GlobalState struct {
	Players          []string
	Rounds           game.Rounds
	Stats            map[string]*game.Stats
	CurrentStandings map[string]*game.Stats
}

func NewEmptyGlobalState() GlobalState {
	return GlobalState{
		Players:          []string{},
		Rounds:           game.Rounds{},
		Stats:            make(map[string]*game.Stats),
		CurrentStandings: make(map[string]*game.Stats),
	}
}

func NewGlobalState(
	players []string,
	currentStandings map[string]*game.Stats,
) (GlobalState, error) {
	stats := game.LoadStats()
	log.L.Info("stats", "len", len(stats))
	for _, s := range stats {
		log.L.Info("", "p", s.Player, "s", s.SuccessRate())
	}
	log.L.Info("", "players before", players)
	players = game.PlayersByStats(players, stats)
	log.L.Info("", "players after", players)

	teams, err := game.GenerateTeams(players)
	if err != nil {
		return NewEmptyGlobalState(), err
	}

	return GlobalState{
		Players:          players,
		Rounds:           game.NewRounds(teams),
		Stats:            stats,
		CurrentStandings: currentStandings,
	}, nil
}

// global state for the app
var state = NewEmptyGlobalState()
