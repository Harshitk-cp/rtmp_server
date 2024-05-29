package config

import (
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
	DefaultHTTPRelayPort int = 9090
)

type Config struct {
	*ServiceConfig  `yaml:",inline"`
	*InternalConfig `yaml:",inline"`
}

type ServiceConfig struct {
	Redis *redis.RedisConfig `yaml:"redis"`

	HttpPort         int           `yaml:"http_port"`
	DebugHandlerPort int           `yaml:"debug_handler_port"`
	RTMPPort         int           `yaml:"rtmp_port"`
	HTTPRelayPort    int           `yaml:"http_relay_port"`
	Logging          logger.Config `yaml:"logging"`
	Development      bool          `yaml:"development"`
}

type InternalConfig struct {
	ServiceName string `yaml:"service_name"`
	NodeID      string `yaml:"node_id"`
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
	return nil
}

func (conf *Config) Init() error {
	conf.NodeID = utils.NewGuid("NE_")
	err := conf.InitDefaults()
	if err != nil {
		return err
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

	conf.InternalConfig = &InternalConfig{}

	if err := conf.Init(); err != nil {
		return nil, err
	}

	return &conf, nil
}
