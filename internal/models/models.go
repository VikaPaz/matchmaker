package models

import (
	"time"
)

type Player struct {
	ID      uint
	Name    string
	Skill   float64
	Latency float64
	Added   time.Time
}

type MatchResponse struct {
	Group   uint
	Skill   []float64
	Latency []float64
	Added   []time.Time
	Players []Player
}

type AddRequest struct {
	Name    string  `json:"name"`
	Skill   float64 `json:"skill"`
	Latency float64 `json:"latency"`
}

type Cluster struct {
	Center Player
	ID     uint
	Count  int
}
