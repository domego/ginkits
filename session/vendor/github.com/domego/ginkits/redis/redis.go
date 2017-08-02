package rediskits

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/domego/ginkits/cache"
	"github.com/domego/gokits/log"
	redis "gopkg.in/redis.v3"
)

type RedisConfig struct {
	Address     string `yaml:"address"`
	Password    string `yaml:"password"`
	Timeout     int    `yaml:"timeout"`
	MaxIdle     int    `yaml:"max_idle"`
	IdleTimeout int    `yaml:"idle_timeout"`
}

func NewClient(c *RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         c.Address,
		Password:     c.Password,
		PoolSize:     c.MaxIdle,
		IdleTimeout:  time.Duration(c.IdleTimeout) * time.Second,
		DialTimeout:  time.Duration(c.Timeout) * time.Second,
		ReadTimeout:  time.Duration(c.Timeout) * time.Second,
		WriteTimeout: time.Duration(c.Timeout) * time.Second,
	})
}

// DeleteCache delete cache
func DeleteCache(rc *redis.Client, key string, retryTimes int) error {
	if rc == nil {
		return nil
	}
	log.Tracef("[cache: %s]: delete cache", key)
	var err error
	for i := 0; i < retryTimes; i++ {
		_, err = rc.Del(key).Result()
		if err == nil {
			return nil
		}
		log.Errorf("delete cache, retry %d, %s", i, err)
	}
	return err
}

// SetModelToCache save model to cache
func SetModelToCache(rc *redis.Client, key string, model interface{}, ttl int) error {
	if rc == nil {
		return nil
	}
	log.Tracef("[cache: set_model_to_cache]: key=%s", key)
	gziped, data, err := cachekits.ToGzipJSON(model)
	if err != nil {
		return err
	}
	cacheFlag := cachekits.CacheFormatJSON
	if gziped {
		cacheFlag = cachekits.CacheFormatJSONGzip
	}
	bs, err := json.Marshal(cachekits.ModelCacheItem{
		Data: data,
		Flag: uint32(cacheFlag),
	})
	if err != nil {
		return err
	}
	_, err = rc.Set(key, string(bs), time.Duration(ttl)*time.Second).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetCacheToModel get cache to model
func GetCacheToModel(rc *redis.Client, key string, model interface{}) (notFound bool) {
	if rc == nil {
		return true
	}
	log.Tracef("[cache: get_cache_to_model], key=%s", key)
	var err error
	it, err := rc.Get(key).Result()
	if it == "" || err == redis.Nil {
		return true
	}
	cacheItem := cachekits.ModelCacheItem{}
	if err = json.Unmarshal([]byte(it), &cacheItem); err != nil {
		log.Errorf("[cache:%s] Unmarshal value error, %s", key, err)
		return true
	}
	switch cacheItem.Flag {
	case cachekits.CacheFormatJSON:
		err = json.Unmarshal(cacheItem.Data, model)
	case cachekits.CacheFormatJSONGzip:
		err = cachekits.FromGzipJSON(cacheItem.Data, model)
	default:
		err = fmt.Errorf("invalid cache formate %d", cacheItem.Flag)
	}
	if err != nil {
		log.Errorf("[cache:%s] %s", key, err)
		return true
	}
	return false
}
