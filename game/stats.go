package game

import (
	"encoding/json"
	"os"
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
