package config

import (
	"errors"
	"os"

	"github.com/livekit/protocol/logger"
	"github.com/livekit/protocol/redis"
	"github.com/livekit/protocol/utils"
	lksdk "github.com/livekit/server-sdk-go/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	DefaultRTMPPort      int = 1935
	DefaultWHIPPort      int = 8080
	DefaultHTTPRelayPort int = 9090
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

func (conf *ServiceConfig) InitDefaults() error {
	if conf.RTMPPort == 0 {
		conf.RTMPPort = DefaultRTMPPort
	}
	if conf.HTTPRelayPort == 0 {
		conf.HTTPRelayPort = DefaultHTTPRelayPort
	}
	if conf.WHIPPort == 0 {
		conf.WHIPPort = DefaultWHIPPort
	}
	return nil
}

func (conf *Config) Init() error {
	conf.NodeID = utils.NewGuid("NE_")
	err := conf.InitDefaults()
	if err != nil {
		return err
	}

	if conf.InternalConfig == nil {
		return errors.New("InternalConfig is nil")
	}

	if err := conf.InitLogger(); err != nil {
		return err
	}

	return nil
}

func (c *Config) GetLoggerFields() logrus.Fields {
	fields := logrus.Fields{
		"logger": c.ServiceName,
	}
	v := c.getLoggerValues()
	for i := 0; i < len(v); i += 2 {
		fields[v[i].(string)] = v[i+1]
	}

	return fields
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

func LoadFromFile(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var conf Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	// Initialize the InternalConfig field
	conf.InternalConfig = &InternalConfig{}

	// Initialize the configuration
	if err := conf.Init(); err != nil {
		return nil, err
	}

	return &conf, nil
}
