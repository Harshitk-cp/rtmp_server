package config

import (
	"github.com/livekit/protocol/logger"
	"github.com/livekit/protocol/redis"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

type Config struct {
	*ServiceConfig  `yaml:",inline"`
	*InternalConfig `yaml:",inline"`
}

type ServiceConfig struct {
	Redis     *redis.RedisConfig `yaml:"redis"`      // required
	ApiKey    string             `yaml:"api_key"`    // required (env LIVEKIT_API_KEY)
	ApiSecret string             `yaml:"api_secret"` // required (env LIVEKIT_API_SECRET)
	WsUrl     string             `yaml:"ws_url"`     // required (env LIVEKIT_WS_URL)

	HealthPort       int           `yaml:"health_port"`
	DebugHandlerPort int           `yaml:"debug_handler_port"`
	PrometheusPort   int           `yaml:"prometheus_port"`
	RTMPPort         int           `yaml:"rtmp_port"` // -1 to disable RTMP
	WHIPPort         int           `yaml:"whip_port"` // -1 to disable WHIP
	HTTPRelayPort    int           `yaml:"http_relay_port"`
	Logging          logger.Config `yaml:"logging"`
	Development      bool          `yaml:"development"`
}

type InternalConfig struct {
	// internal
	ServiceName string `yaml:"service_name"`
	NodeID      string `yaml:"node_id"` // Do not provide, will be overwritten
}

type LoggingConfig struct {
	JSON bool
}

func (c *Config) GetLoggerFields() map[string]interface{} {
	return nil
}

func (c *Config) InitLogger(values ...interface{}) error {
	zl, err := logger.NewZapLogger(&c.Logging)
	if err != nil {
		return err
	}

	values = append(c.getLoggerValues(), values...)
	l := zl.WithValues(values...)
	logger.SetLogger(l, c.ServiceName)
	lksdk.SetLogger(l)

	return nil
}

func (c *Config) getLoggerValues() []interface{} {
	return []interface{}{"nodeID", c.NodeID}
}
