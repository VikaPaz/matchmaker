package service

import (
	"fmt"
	"math"
	"time"

	"github.com/VikaPaz/matchmaker/internal/models"
	"github.com/sirupsen/logrus"
)

var superPlayer = models.Player{
	Name:    "Super",
	Skill:   5000,
	Latency: 300,
	Added:   time.Now(),
}

var clusters = make([]models.Cluster, 3)

var id uint

var wPostgre bool

type MatchService struct {
	repo    Repo
	postgre Postgre
	log     *logrus.Logger
}

type Repo interface {
	QueryMatching(clusterID string, count int64) ([]models.Player, error)
	QueryDel(string, uint) error
	QueryAdd(player models.Player, cluster string, score float64) error
}

type Postgre interface {
	Create(models.Player) error
	Delete(uint) error
}

func NewService(r Repo, p Postgre, maxSkill float64, maxLatency float64, logger *logrus.Logger, writePostgres bool) *MatchService {
	superPlayer.Skill = maxSkill
	superPlayer.Latency = maxLatency
	wPostgre = writePostgres
	return &MatchService{
		repo:    r,
		postgre: p,
		log:     logger,
	}
}

func (s *MatchService) Matching(groupSize int) error {
	var countGroup uint
	for _, cluster := range clusters {
		if cluster.Center == (models.Player{}) {
			continue
		}
		s.log.Debug("Matching...")
		players, err := s.repo.QueryMatching(fmt.Sprintf("%d", cluster.ID), int64(groupSize))
		if err != nil {
			s.log.Errorf("Error create group: %v", err)
			return err
		}
		if len(players) < groupSize {
			continue
		}

		for _, player := range players {
			err := s.repo.QueryDel(fmt.Sprintf("%d", cluster.ID), player.ID)
			if err != nil {
				return err
			}
			if !wPostgre {
				continue
			}
			err = s.postgre.Delete(player.ID)
			if err != nil {
				return err
			}
		}

		err = updateCentersDel(players, int(cluster.ID))
		if err != nil {
			return err
		}

		countGroup++

		resp := creatResponse(players, countGroup)

		showResp(resp)

	}
	return nil
}

func updateCentersDel(players []models.Player, id int) error {
	cluster := clusters[id]
	for _, player := range players {
		cluster.Center.Skill = (cluster.Center.Skill*float64(cluster.Count) - player.Skill) / (float64(cluster.Count) - 1)
		cluster.Center.Latency = (cluster.Center.Latency*float64(cluster.Count) - player.Latency) / (float64(cluster.Count) - 1)
	}
	clusters[id] = cluster
	return nil
}

func showResp(resp models.MatchResponse) {
	fmt.Printf("Group: %v \n", resp.Group)
	fmt.Printf("Skill min/max/avg: %v \n", resp.Skill)
	fmt.Printf("Latency min/max/avg: %v \n", resp.Latency)
	fmt.Printf("Added min/max/avg: %v \n", resp.Added)
	fmt.Printf("Payers: ")
	for _, player := range resp.Players {
		fmt.Print(player.Name, " ")
	}
	fmt.Println()

}

func creatResponse(players []models.Player, group uint) models.MatchResponse {
	resp := models.MatchResponse{
		Group:   group,
		Skill:   make([]float64, 3),
		Latency: make([]float64, 3),
		Added:   make([]time.Time, 3),
		Players: players,
	}
	var sumSkill, sumLatency float64
	var sumAdded int64
	resp.Added[0] = players[0].Added
	for _, player := range players {
		sumSkill += player.Skill
		sumLatency += player.Latency
		sumAdded += player.Added.UTC().Unix()
		if player.Skill < resp.Skill[0] {
			resp.Skill[0] = player.Skill
		}
		if player.Skill > resp.Skill[1] {
			resp.Skill[1] = player.Skill
		}
		if player.Latency < resp.Latency[0] {
			resp.Latency[0] = player.Latency
		}
		if player.Latency > resp.Latency[1] {
			resp.Latency[1] = player.Latency
		}
		if player.Added.UTC().Unix() < resp.Added[0].UTC().Unix() {
			resp.Added[0] = player.Added
		}
		if player.Added.UTC().Unix() > resp.Added[1].UTC().Unix() {
			resp.Added[1] = player.Added
		}
	}
	countPlayers := float64(len(players))
	resp.Skill[2] = sumSkill / countPlayers
	resp.Latency[2] = sumLatency / countPlayers
	resp.Added[2] = time.Unix(sumAdded/int64(countPlayers), 0)

	return resp
}

func (s *MatchService) AddPlayer(req models.AddRequest) (models.Player, error) {
	playerID := assignID()

	player := models.Player{
		ID:      playerID,
		Name:    req.Name,
		Skill:   req.Skill,
		Latency: req.Latency,
		Added:   time.Now(),
	}

	clasterID, err := assignCluster(player)
	if err != nil {
		s.log.Error("Error assign Cluster")
		return models.Player{}, err
	}

	score := euclideanDistance(player, superPlayer)

	err = s.repo.QueryAdd(player, fmt.Sprintf("%d", clasterID), score)
	if err != nil {
		s.log.Error("Error add to Cluster")
		return models.Player{}, err
	}

	if wPostgre {
		err = s.postgre.Create(player)
		if err != nil {
			s.log.Error("Error add to Postgres")
			return models.Player{}, err
		}
	}

	err = updateCenters(player, int(clasterID))
	if err != nil {
		s.log.Error("Error update centers of Clusters")
		return models.Player{}, err
	}
	return player, nil
}

func updateCenters(player models.Player, id int) error {
	count := float64(clusters[id].Count)
	clusters[id].Center.Skill = (clusters[id].Center.Skill*count + player.Skill) / (count + 1)
	clusters[id].Center.Latency = (clusters[id].Center.Latency*count + player.Latency) / (count + 1)
	clusters[id].Count += 1
	return nil
}

func assignCluster(player models.Player) (uint, error) {
	minDistance := math.Inf(1)
	var clusterID uint
	for i, cluster := range clusters {
		if cluster.Center == (models.Player{}) {
			clusters[i].ID = uint(i)
			return uint(i), nil
		}
		dist := euclideanDistance(player, cluster.Center)
		if dist < minDistance {
			minDistance = dist
			clusterID = cluster.ID
		}
	}
	return clusterID, nil
}

func euclideanDistance(p1, p2 models.Player) float64 {
	return math.Sqrt(math.Pow(p1.Skill-p2.Skill, 2) + math.Pow(p1.Latency-p2.Latency, 2))
}

func assignID() uint {
	id++
	return id
}
