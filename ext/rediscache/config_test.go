package rediscache_test

import (
	"github.com/krakendio/krakend-ce/v2/ext/rediscache"
	"testing"
	"time"

	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/encoding"
	"github.com/stretchr/testify/assert"
)

func TestConfig_ok(t *testing.T) {
	t.Run("With empty config returns default config", func(t *testing.T) {
		cfg := &config.Backend{
			Decoder: encoding.JSONDecoder,
			ExtraConfig: map[string]interface{}{
				rediscache.Namespace: map[string]interface{}{},
			},
		}

		c, err := rediscache.ConfigGetter(cfg)

		assert.NoError(t, err)
		assert.Equal(t, "memory", c.(rediscache.Config).Type)
	})

	t.Run("With type memory is parsed correctly", func(t *testing.T) {
		cfg := &config.Backend{
			Decoder: encoding.JSONDecoder,
			ExtraConfig: map[string]interface{}{
				rediscache.Namespace: map[string]interface{}{
					"type": "memory",
				},
			},
		}

		c, err := rediscache.ConfigGetter(cfg)

		assert.NoError(t, err)
		assert.Equal(t, "memory", c.(rediscache.Config).Type)
		assert.Empty(t, c.(rediscache.Config).RedisConfig)
	})

	t.Run("With redis type loads default values", func(t *testing.T) {
		cfg := &config.Backend{
			Decoder: encoding.JSONDecoder,
			ExtraConfig: map[string]interface{}{
				rediscache.Namespace: map[string]interface{}{
					"type": "redis",
					"redis": map[string]interface{}{
						"address": "localhost:6379",
						"mode":    "redis",
					},
				},
			},
		}

		c, err := rediscache.ConfigGetter(cfg)
		rc := c.(rediscache.Config).RedisConfig
		assert.NoError(t, err)
		assert.Equal(t, "redis", c.(rediscache.Config).Type)
		assert.Equal(t, rediscache.RedisDefaultDialTimeout, rc.DialTimeout)
		assert.Equal(t, rediscache.RedisDefaultReadTimeout, rc.ReadTimeout)
		assert.Equal(t, rediscache.RedisDefaultWriteTimeout, rc.WriteTimeout)
		assert.Equal(t, rediscache.RedisDefaultMaxRetries, rc.MaxRetries)
		assert.Equal(t, rediscache.RedisDefaultIdleTimeout, rc.IdleTimeout)
		assert.Equal(t, rediscache.RedisDefaultIdleCheckFrequency, rc.IdleCheckFrequency)
		assert.Equal(t, rediscache.RedisDefaultPoolSize, rc.PoolSize)
		assert.Equal(t, rediscache.RedisDefaultPoolTimeout, rc.PoolTimeout)
		assert.Equal(t, rediscache.RedisDefaultTtl, rc.Ttl)
	})

	t.Run("With redis type loads values from config", func(t *testing.T) {
		address := "loremipsum:6379"
		mode := "rediscluster"
		dt := "10us"
		rt := "18us"
		wt := "25us"
		mr := 666
		it := "99s"
		icf := "14ms"
		ps := 77
		pt := "878s"
		ttl := "65432s"
		cfg := &config.Backend{
			Decoder: encoding.JSONDecoder,
			ExtraConfig: map[string]interface{}{
				rediscache.Namespace: map[string]interface{}{
					"type": "redis",
					"redis": map[string]interface{}{
						"address":            address,
						"mode":               mode,
						"dialTimeout":        dt,
						"readTimeout":        rt,
						"writeTimeout":       wt,
						"maxRetries":         mr,
						"idleTimeout":        it,
						"idleCheckFrequency": icf,
						"poolSize":           ps,
						"poolTimeout":        pt,
						"ttl":                ttl,
					},
				},
			},
		}

		c, err := rediscache.ConfigGetter(cfg)
		rc := c.(rediscache.Config).RedisConfig

		assert.NoError(t, err)
		assert.Equal(t, "redis", c.(rediscache.Config).Type)
		assert.Equal(t, address, rc.Address)
		assert.Equal(t, mode, rc.Mode)
		exDt, _ := time.ParseDuration(dt)
		assert.Equal(t, exDt, rc.DialTimeout)
		exRt, _ := time.ParseDuration(rt)
		assert.Equal(t, exRt, rc.ReadTimeout)
		exWt, _ := time.ParseDuration(wt)
		assert.Equal(t, exWt, rc.WriteTimeout)
		assert.Equal(t, mr, rc.MaxRetries)
		exIt, _ := time.ParseDuration(it)
		assert.Equal(t, exIt, rc.IdleTimeout)
		exIcf, _ := time.ParseDuration(icf)
		assert.Equal(t, exIcf, rc.IdleCheckFrequency)
		assert.Equal(t, ps, rc.PoolSize)
		exPt, _ := time.ParseDuration(pt)
		assert.Equal(t, exPt, rc.PoolTimeout)
		exTtl, _ := time.ParseDuration(ttl)
		assert.Equal(t, exTtl, rc.Ttl)
	})
}

func TestConfig_ko(t *testing.T) {
	t.Run("Error thrown when missing redis block", func(t *testing.T) {
		cfg := &config.Backend{
			Decoder: encoding.JSONDecoder,
			ExtraConfig: map[string]interface{}{
				rediscache.Namespace: map[string]interface{}{
					"type": "redis",
				},
			},
		}

		c, err := rediscache.ConfigGetter(cfg)

		assert.Error(t, err)
		assert.Nil(t, c)
	})

	t.Run("Error thrown when missing redis mode", func(t *testing.T) {
		cfg := &config.Backend{
			Decoder: encoding.JSONDecoder,
			ExtraConfig: map[string]interface{}{
				rediscache.Namespace: map[string]interface{}{
					"type":  "redis",
					"redis": map[string]interface{}{},
				},
			},
		}

		c, err := rediscache.ConfigGetter(cfg)

		assert.Error(t, err)
		assert.Nil(t, c)
	})

	t.Run("Error thrown when missing redis address", func(t *testing.T) {
		cfg := &config.Backend{
			Decoder: encoding.JSONDecoder,
			ExtraConfig: map[string]interface{}{
				rediscache.Namespace: map[string]interface{}{
					"type": "redis",
					"redis": map[string]interface{}{
						"mode": "redis",
					},
				},
			},
		}

		c, err := rediscache.ConfigGetter(cfg)

		assert.Error(t, err)
		assert.Nil(t, c)
	})

	t.Run("Error thrown when missing invalid redis mode", func(t *testing.T) {
		cfg := &config.Backend{
			Decoder: encoding.JSONDecoder,
			ExtraConfig: map[string]interface{}{
				rediscache.Namespace: map[string]interface{}{
					"type": "redis",
					"redis": map[string]interface{}{
						"mode": "whatever",
					},
				},
			},
		}

		c, err := rediscache.ConfigGetter(cfg)

		assert.Error(t, err)
		assert.Nil(t, c)
	})
}
