package web

import (
	"github.com/AngelVI13/foos/game"
)

type GlobalState struct {
	Players []string
	Rounds  game.Rounds
	Stats   map[string]*game.Stats
}

func NewEmptyGlobalState() GlobalState {
	return GlobalState{
		Players: []string{},
		Rounds:  game.Rounds{},
	}
}

func NewGlobalState(players []string) (GlobalState, error) {
	teams, err := game.GenerateTeams(players)
	if err != nil {
		return NewEmptyGlobalState(), err
	}

	stats := game.LoadStats()

	return GlobalState{
		Players: players,
		Rounds:  game.NewRounds(teams),
		Stats:   stats,
	}, nil
}

// global state for the app
var state = NewEmptyGlobalState()
