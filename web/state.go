package web

import (
	"github.com/AngelVI13/foos/game"
)

type GlobalState struct {
	Players []string
	Rounds  game.Rounds
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

	return GlobalState{
		Players: players,
		Rounds:  game.NewRounds(teams),
	}, nil
}

// global state for the app
var state = NewEmptyGlobalState()
