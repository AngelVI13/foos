package game

import (
	"fmt"

	"github.com/google/uuid"
)

type Result string

const (
	Win   Result = "win"
	Loss  Result = "loss"
	Empty Result = "n/a"
)

const TBD = "tbd"

type Team struct {
	Player1  string
	Player2  string
	matchIdx int
	scores   []int
	results  []Result
}

func NewTeam(p1, p2 string) Team {
	return Team{
		Player1:  p1,
		Player2:  p2,
		matchIdx: 0,
		scores:   []int{0},
		results:  []Result{Empty},
	}
}

func (m *Team) Score() int {
	return m.scores[m.matchIdx]
}

func (m *Team) ScoreForMatch(matchIdx int) int {
	if matchIdx > m.matchIdx {
		return 0
	}
	return m.scores[matchIdx]
}

func (m *Team) CurrentMatch() int {
	return m.matchIdx
}

func (m *Team) AllScores() int {
	sum := 0
	for _, score := range m.scores {
		sum += score
	}
	return sum
}

func (m *Team) SetScore(v int) {
	m.scores[m.matchIdx] = v
}

func (m *Team) SetResult(v Result) {
	m.results[m.matchIdx] = v
}

func (m *Team) MatchDone() {
	m.matchIdx++
	m.scores = append(m.scores, 0)
	m.results = append(m.results, Empty)
}

func (m *Team) Result(matchIdx int) Result {
	if matchIdx >= len(m.results) {
		return Empty
	}
	return m.results[matchIdx]
}

func (m Team) String() string {
	return fmt.Sprintf("%s & %s", m.Player1, m.Player2)
}

type Match struct {
	Id    string
	team1 *Team
	team2 *Team
}

func NewMatch(t1, t2 *Team) Match {
	id := uuid.New()
	return Match{
		Id:    id.String(),
		team1: t1,
		team2: t2,
	}
}

func NewTbdMatch() Match {
	emptyTeam := NewTeam(TBD, TBD)
	return NewMatch(&emptyTeam, &emptyTeam)
}

func (m *Match) Teams() []*Team {
	return []*Team{m.team1, m.team2}
}

func (m *Match) End() {
	if m.team1.Score() > m.team2.Score() {
		m.team1.SetResult(Win)
		m.team2.SetResult(Loss)
	} else {
		m.team1.SetResult(Loss)
		m.team2.SetResult(Win)
	}

	m.team1.MatchDone()
	m.team2.MatchDone()
}
