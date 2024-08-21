package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/VikaPaz/matchmaker/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisRepo struct {
	conn *redis.Client
	log  *logrus.Logger
}

func NewRepo(conn *redis.Client, logger *logrus.Logger) *RedisRepo {
	return &RedisRepo{
		conn: conn,
		log:  logger,
	}
}

func (r *RedisRepo) QueryMatching(clusterID string, count int64) ([]models.Player, error) {
	matchedPlayers := make([]models.Player, 0)

	playerIDs, err := r.conn.ZRange(context.Background(), clusterID, 0, count-1).Result()
	if err != nil {
		r.log.Error("Error retrieving player IDs from sorted set (ZRange)")
		return nil, err
	}

	r.log.Debugf("ZRange cluster: %v IDs: %v", clusterID, playerIDs)

	for _, playerID := range playerIDs {
		playerJSON, err := r.conn.HGet(context.Background(), "players", playerID).Result()
		if err != nil {
			r.log.Error("Error retrieving player data from hash (HGet)")
			return nil, err
		}

		var player models.Player
		err = json.Unmarshal([]byte(playerJSON), &player)
		if err != nil {
			r.log.Error("Error unmarshalling player JSON")
			return nil, err
		}

		r.log.Debugf("HGet players ID: %v date: %v", playerID, player)

		matchedPlayers = append(matchedPlayers, player)
	}

	return matchedPlayers, nil
}

func (r *RedisRepo) QueryDel(clusterKey string, playerID uint) error {
	playerIDStr := fmt.Sprintf("%d", playerID)

	err := r.conn.ZRem(context.Background(), clusterKey, playerID).Err()
	if err != nil {
		r.log.Error("Error removing player from sorted set (ZREM)")
		return err
	}

	r.log.Debugf("ZRem cluster: %v, id: %v", clusterKey, playerIDStr)

	err = r.conn.HDel(context.Background(), "players", playerIDStr).Err()
	if err != nil {
		r.log.Error("Error removing player from hash (HDEL)")
		return err
	}

	r.log.Debugf("HDel players id: %v", playerIDStr)

	return nil
}

func (r *RedisRepo) QueryAdd(player models.Player, cluster string, score float64) error {
	var err error
	playerJSON, err := json.Marshal(player)
	if err != nil {
		r.log.Error("Error serializing player to JSON")
		return err
	}
	err = r.conn.HSet(context.Background(), "players", fmt.Sprintf("%d", player.ID), playerJSON).Err()
	if err != nil {
		r.log.Error("Error HSET")
		return err
	}

	r.log.Debugf("Hset players id: %v player: %v", player.ID, player)

	z := redis.Z{
		Score:  score,
		Member: player.ID,
	}
	err = r.conn.ZAdd(context.Background(), cluster, z).Err()
	if err != nil {
		r.log.Error("Error ZADD")
		return err
	}

	r.log.Debugf("Zadd cluster: %v, score: %v, id: %v", cluster, z.Score, z.Member)

	return nil
}
