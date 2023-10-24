package game

import (
	"encoding/json"
	"os"
	"sort"
)

const (
	StatsFile = "stats.json"
)

type Stats struct {
	Player string
	Score  int
	Won    int
	Lost   int
}

// SuccessRate calculate as percentage of score per match
// 10 is maximum number of points per match
func (s Stats) SuccessRate() int {
	return int(100 * (float64(s.Score) / float64(10.0*(s.Won+s.Lost))))
}

func LoadStats() map[string]*Stats {
	b, err := os.ReadFile(StatsFile)
	if err != nil {
		return map[string]*Stats{}
	}

	var stats map[string]*Stats
	err = json.Unmarshal(b, &stats)
	if err != nil {
		return map[string]*Stats{}
	}

	return stats
}

func SaveStats(s map[string]*Stats) error {
	os.Remove(StatsFile)

	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(StatsFile, b, 0o666)
}

func OrderedStatsSlice(stats map[string]*Stats) []*Stats {
	var statsCopy []*Stats
	for _, s := range stats {
		statsCopy = append(statsCopy, s)
	}

	sort.Slice(statsCopy, func(i, j int) bool {
		if statsCopy[i].Score == statsCopy[j].Score {
			return statsCopy[i].Won > statsCopy[j].Won
		}
		return statsCopy[i].Score > statsCopy[j].Score
	})

	return statsCopy
}
