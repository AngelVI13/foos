package web

import "github.com/AngelVI13/foos/game"

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

func NewGlobalState(players []string, teams []game.Team) GlobalState {
	return GlobalState{
		Players: players,
		Rounds:  game.NewRounds(teams),
	}
}

// global state for the app
var state = NewEmptyGlobalState()
