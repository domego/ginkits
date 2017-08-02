package sessionkits

import (
	"time"

	"github.com/domego/ginkits/redis"
	"github.com/domego/gokits/log"
	redis "gopkg.in/redis.v3"
)

type redisStore struct {
	rds *redis.Client
}

func NewRedisStore(cfg *rediskits.RedisConfig) *redisStore {
	return &redisStore{
		rds: rediskits.NewClient(cfg),
	}
}

func (s *redisStore) GetSession(sessionID, key string) []byte {
	r, err := s.rds.Get(sessionID + ":" + key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Errorf("GetSession: %s", err)
		}
		return nil
	}

	if r == "" {
		return nil
	}
	return []byte(r)
}

func (s *redisStore) SetSession(sessionID, key string, data []byte) {
	err := s.rds.Set(sessionID+":"+key, data, time.Duration(SessionTimeout)*time.Minute).Err()
	if err != nil {
		log.Errorf("SetSession: %s", err)
	}
}

func (s *redisStore) ClearSession(sessionID, key string) {
	err := s.rds.Del(sessionID + ":" + key).Err()
	if err != nil {
		log.Errorf("ClearSession: %s", err)
	}
}
