package game

import "fmt"

type Result string

const (
	Win   Result = "win"
	Loss  Result = "loss"
	Empty Result = "n/a"
)

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

func (m *Team) SetScore(v int) {
	m.matchIdx = len(m.scores)
	m.scores = append(m.scores, v)
}

func (m *Team) SetResult(v Result) {
	m.results = append(m.results, v)
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
	team1 *Team
	team2 *Team
}

func NewMatch(t1, t2 *Team) Match {
	return Match{
		team1: t1,
		team2: t2,
	}
}

func NewTbdMatch() Match {
	emptyTeam := NewTeam("tbd", "tbd")
	return Match{
		team1: &emptyTeam,
		team2: &emptyTeam,
	}
}

func (m *Match) Teams() []*Team {
	return []*Team{m.team1, m.team2}
}

func (m *Match) Result() (winner *Team, loser *Team) {
	if m.team1.Score() > m.team2.Score() {
		m.team1.SetResult(Win)
		m.team2.SetResult(Loss)
		return m.team1, m.team2
	}
	m.team1.SetResult(Loss)
	m.team2.SetResult(Win)
	return m.team2, m.team1
}

func (m *Match) AddScores(team1Score, team2Score int) {
	m.team1.SetScore(team1Score)
	m.team2.SetScore(team2Score)
}
