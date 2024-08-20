package service

import (
	"fmt"
	"math"
	"strconv"
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

type MatchService struct {
	repo Repo
	log  *logrus.Logger
}

type Repo interface {
	QueryMatching(clusterID uint, count uint) ([]models.Player, error)
	QueryDel(uint, uint) error
	QueryAdd(player models.Player, cluster string, score float64) error
	// OuerySumSkill(uint) (float64, error)
	// OuerySumLatency(uint) (float64, error)
	// OueryCountPlayers(uint) (float64, error)
}

func NewService(r Repo, maxSkill float64, maxLatency float64, logger *logrus.Logger) *MatchService {
	superPlayer.Skill = maxSkill
	superPlayer.Latency = maxLatency
	return &MatchService{
		repo: r,
		log:  logger,
	}
}

func (s *MatchService) Matching(groupSize int) error {
	var countGroup uint
	for _, cluster := range clusters {
		if cluster.Center == (models.Player{}) {
			continue
		}
		s.log.Debug("Matching...")
		players, err := s.repo.QueryMatching(cluster.ID, uint(groupSize))
		if err != nil {
			s.log.Errorf("Error create group: %v", err)
			return err
		}
		if len(players) < groupSize {
			continue
		}

		for _, player := range players {
			err := s.repo.QueryDel(cluster.ID, player.ID)
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
		// TODO AddedS

	}
	clusters[id] = cluster
	fmt.Println(clusters[id])
	return nil
}

func showResp(resp models.MatchResponse) {
	fmt.Println(resp.Group)
	fmt.Println(resp.Skill)
	fmt.Println(resp.Latency)
	fmt.Println(resp.Added)
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
			resp.Latency[0] = player.Latency
		}
		if player.Added.UTC().Unix() > resp.Added[1].UTC().Unix() {
			resp.Latency[1] = player.Latency
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

	err = s.repo.QueryAdd(player, strconv.FormatUint(uint64(clasterID), 10), score)
	if err != nil {
		s.log.Error("Error add to Cluster")
		return models.Player{}, err
	}

	err = updateCenters(player, int(clasterID))
	if err != nil {
		s.log.Error("Error update centers of Clusters")
		return models.Player{}, err
	}

	fmt.Println(clasterID, player)

	fmt.Println(clusters)

	return player, nil
}

func updateCenters(player models.Player, id int) error {
	count := float64(clusters[id].Count)
	clusters[id].Center.Skill = (clusters[id].Center.Skill*count + player.Skill) / (count + 1)
	clusters[id].Center.Latency = (clusters[id].Center.Latency*count + player.Latency) / (count + 1)
	// TODO Added
	clusters[id].Count += 1
	return nil
}

// clusters forom Repo

// func updateCenters(s *MatchService, id uint) error {
// 	sumSkill, err := s.repo.OuerySumSkill(id)
// 	if err != nil {
// 		s.log.Error("Error query sum of skill")
// 		return err
// 	}
// 	sumLatency, err := s.repo.OuerySumLatency(id)
// 	if err != nil {
// 		s.log.Error("Error query sum of latency")
// 		return err
// 	}
// 	count, err := s.repo.OueryCountPlayers(id)
// 	if err != nil {
// 		s.log.Error("Error query count of players")
// 		return err
// 	}
// 	clusters[id].Center.Skill = sumSkill / count
// 	clusters[id].Center.Latency = sumLatency / count
// 	return nil
// }

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

		fmt.Println(dist)
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
