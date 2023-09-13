package rediscache_test

import (
	"errors"
	"github.com/krakendio/krakend-ce/v2/ext/rediscache"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"

	httpcache2 "github.com/gregjones/httpcache"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRedisCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := rediscache.NewMockClient(ctrl)
	ttl := 10 * time.Millisecond

	c := rediscache.NewRedisCache(client, ttl)

	t.Run("Calls client to set new value", func(t *testing.T) {
		k := "setkey"
		v := []byte("aresponse")

		client.EXPECT().Set(gomock.Any(), k, v, ttl).Times(1)

		c.Set(k, v)
	})

	t.Run("Calls client to delete existing value", func(t *testing.T) {
		k := "delkey"

		client.EXPECT().Del(gomock.Any(), k)

		c.Delete(k)
	})

	t.Run("Get returns ko when cant find in cache", func(t *testing.T) {
		k := "getko"

		res := redis.NewStringResult("", errors.New(""))
		client.EXPECT().Get(gomock.Any(), k).Times(1).Return(res)

		_, ok := c.Get(k)

		assert.False(t, ok)
	})

	t.Run("Get returns value when found in cache", func(t *testing.T) {
		k := "getok"
		v := "foundincache"
		res := redis.NewStringResult(v, nil)
		client.EXPECT().Get(gomock.Any(), k).Times(1).Return(res)

		r, ok := c.Get(k)

		assert.True(t, ok)
		assert.Equal(t, v, string(r))
	})
}

func TestNewRedis(t *testing.T) {
	rc := rediscache.NewRedis(rediscache.RedisConfig{})

	assert.IsType(t, &redis.Client{}, rc)
}

func TestNewRedisCluster(t *testing.T) {
	rc := rediscache.NewRedisCluster(rediscache.RedisConfig{})

	assert.IsType(t, &redis.ClusterClient{}, rc)
}

func TestRedisCacheTransport_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rc := rediscache.NewMockCache(ctrl)

	rct := rediscache.NewRedisCacheTransport(rc)

	assert.IsType(t, &httpcache2.Transport{}, rct)
}
