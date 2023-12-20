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
	JudgementDay     bool
}

func NewEmptyGlobalState() GlobalState {
	return GlobalState{
		Players:          []string{},
		Rounds:           game.Rounds{},
		Stats:            make(map[string]*game.Stats),
		CurrentStandings: make(map[string]*game.Stats),
		JudgementDay:     false,
	}
}

func NewGlobalState(
	players []string,
	currentStandings map[string]*game.Stats,
	judgementDay bool,
) (GlobalState, error) {
	stats := game.LoadStats()
	log.L.Info("stats", "len", len(stats))
	for _, s := range stats {
		log.L.Info("", "p", s.Player, "s", s.SuccessRate())
	}
	log.L.Info("", "players before", players)

	if judgementDay {
		players = game.PlayersByRankings(players, stats)
	} else {
		players = game.PlayersBySuccessRate(players, stats)
	}
	log.L.Info("", "players after", players)

	playersRankings := game.PlayersRankings(stats)

	teams, err := game.GenerateTeams(players, judgementDay, playersRankings)
	if err != nil {
		return NewEmptyGlobalState(), err
	}
	log.L.Info("", "team before", teams)
	teams = game.TeamsByStats(teams)
	log.L.Info("", "team after", teams)

	return GlobalState{
		Players:          players,
		Rounds:           game.NewRounds(teams),
		Stats:            stats,
		CurrentStandings: currentStandings,
		JudgementDay:     judgementDay,
	}, nil
}

// global state for the app
var state = NewEmptyGlobalState()
