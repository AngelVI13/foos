package game

var RoundNames = map[int]string{
	0: "Semifinal",
	1: "Final",
}

type Round struct {
	Id      int
	Name    string
	Matches []Match
}

func NewRound(id int, teams []*Team) Round {
	var matches []Match
	if id == 0 {
		for i := 2; i <= len(teams); i += 2 {
			matches = append(matches, NewMatch(teams[i-2], teams[i-1]))
		}
	} else if id != 0 && teams[0].Result(id-1) == Empty {
		for i := 2; i <= len(teams); i += 2 {
			matches = append(matches, NewTbdMatch())
		}
	} else {
		// Winners from previous round play with other winners &
		// losers play with losers
		var winners []*Team
		var losers []*Team

		for _, team := range teams {
			switch team.Result(id - 1) {
			case Win:
				winners = append(winners, team)
			case Loss:
				losers = append(losers, team)
			default:
				panic("something went terribly wrong")
			}
		}

		for i := 2; i <= len(winners); i += 2 {
			matches = append(matches, NewMatch(winners[i-2], winners[i-1]))
		}

		for i := 2; i <= len(losers); i += 2 {
			matches = append(matches, NewMatch(losers[i-2], losers[i-1]))
		}
	}

	return Round{
		Id:      id,
		Name:    RoundNames[id],
		Matches: matches,
	}
}

type Rounds struct {
	Teams        []*Team
	All          []Round
	CurrentRound int
}

func NewRounds(teams []Team) Rounds {
	var teamPtrs []*Team
	for i := range teams {
		teamPtrs = append(teamPtrs, &teams[i])
	}

	var allRounds []Round
	if len(teams) == 4 {
		for i := 0; i <= len(teams)/4; i++ {
			allRounds = append(allRounds, NewRound(i, teamPtrs))
		}
	} else {
		allRounds = append(allRounds, NewRound(0, teamPtrs))
	}

	return Rounds{
		Teams:        teamPtrs,
		All:          allRounds,
		CurrentRound: 0,
	}
}

func (r *Rounds) NextRound() {
	var allRounds []Round
	for i := 0; i <= len(r.Teams)/4; i++ {
		allRounds = append(allRounds, NewRound(i, r.Teams))
	}

	r.All = allRounds
	r.CurrentRound++
}

func (r *Rounds) ResultsTable() []*Team {
	return nil
}

func (r *Rounds) FindMatchById(id string) *Match {
	for _, round := range r.All {
		for _, match := range round.Matches {
			if id == match.Id {
				return &match
			}
		}
	}
	return nil
}
