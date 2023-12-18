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
	Player1     string
	Player2     string
	Player1Rank int
	Player2Rank int
	matchIdx    int
	scores      []int
	results     []Result
}

func NewTeam(p1, p2 string, p1Rank, p2Rank int) Team {
	return Team{
		Player1:     p1,
		Player2:     p2,
		Player1Rank: p1Rank,
		Player2Rank: p2Rank,
		matchIdx:    0,
		scores:      []int{0},
		results:     []Result{Empty},
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

func (m *Team) Wins() int {
	num := 0
	for _, result := range m.results {
		if result == Win {
			num++
		}
	}
	return num
}

func (m *Team) Losses() int {
	num := 0
	for _, result := range m.results {
		if result == Loss {
			num++
		}
	}
	return num
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

func (m *Team) Player1RankClass() string {
	return rankClass(m.Player1Rank)
}

func (m *Team) Player2RankClass() string {
	return rankClass(m.Player2Rank)
}

func (m Team) String() string {
	return fmt.Sprintf("%s & %s", m.Player1, m.Player2)
}

func rankClass(rank int) string {
	switch rank {
	case 1:
		return "gold"
	case 2:
		return "silver"
	case 3:
		return "bronze"
	case 4:
		return "fourth"
	default:
		return ""
	}
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
	emptyTeam := NewTeam(TBD, TBD, -1, -1)
	return NewMatch(&emptyTeam, &emptyTeam)
}

func (m *Match) Teams() []*Team {
	return []*Team{m.team1, m.team2}
}

func (m *Match) End() map[string]*Stats {
	stats := map[string]*Stats{}

	if m.team1.Score() > m.team2.Score() {
		m.team1.SetResult(Win)
		m.team2.SetResult(Loss)
	} else {
		m.team1.SetResult(Loss)
		m.team2.SetResult(Win)
	}

	for _, team := range []*Team{m.team1, m.team2} {
		for _, player := range []string{team.Player1, team.Player2} {
			stats[player] = &Stats{}
			stats[player].Player = player
			stats[player].Score = team.Score()

			if team.Result(team.CurrentMatch()) == Win {
				stats[player].Won++
			} else {
				stats[player].Lost++
			}
		}
	}

	m.team1.MatchDone()
	m.team2.MatchDone()

	return stats
}
