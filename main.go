package main

import (
	"errors"
	"flag"
	"github.com/aliyun/aliyun-log-jaeger/sls_store"
	"github.com/hashicorp/go-hclog"
	"github.com/jaegertracing/jaeger/plugin/storage/grpc"
	"github.com/jaegertracing/jaeger/plugin/storage/grpc/shared"
	"github.com/spf13/viper"
	"time"
)

const (
	DefaultLookBack = 6
)

var configPath string

type Configuration struct {
	Endpoint     string        `yaml:"endpoint"`
	AccessKeyID  string        `yaml:"accessKeyId"`
	AccessSecret string        `yaml:"accessSecret"`
	Project      string        `yaml:"project"`
	Instance     string        `yaml:"instance"`
	MaxLookBack  time.Duration `yaml:maxLookBack`
}

var logger = hclog.New(&hclog.LoggerOptions{
	Level:      hclog.Info,
	Name:       "aliyun-log-jaeger-plugin",
	JSONFormat: true,
})

func main() {
	flag.StringVar(&configPath, "config", "", "Path to the alibaba log jaeger plugin's configuration file")
	flag.Parse()

	configuration, err := initialParameters(configPath, logger)
	if err != nil {
		logger.Error("Fatal error config file: %w\n", err)
		return
	}

	var plugin = sls_store.NewSLSStorageForJaegerPlugin(
		configuration.Endpoint,
		configuration.AccessKeyID,
		configuration.AccessSecret,
		configuration.Project,
		configuration.Instance,
		configuration.MaxLookBack,
		logger,
	)

	grpc.Serve(&shared.PluginServices{
		Store: plugin,
	})

	logger.Info("SLS jaeger plugin initialized Successfully")
}

func initialParameters(configPath string, logger hclog.Logger) (*Configuration, error) {
	v := viper.New()
	v.AutomaticEnv()

	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			logger.Error("Failed to read file.", "Exception", err)
			return nil, err
		}
	}

	configuration := &Configuration{}
	if err := configuration.InitFromViper(v); err != nil {
		return nil, err
	} else {
		return configuration, nil
	}
}

func (c *Configuration) InitFromViper(v *viper.Viper) error {
	c.AccessSecret = v.GetString("ACCESS_KEY_SECRET")
	if c.AccessSecret == "" {
		logger.Error("The AccessSecret can't be empty")
		return errors.New("The AccessSecret can't be empty")
	}

	c.AccessKeyID = v.GetString("ACCESS_KEY_ID")
	if c.AccessKeyID == "" {
		logger.Error("The access key id can't be empty")
		return errors.New("The access key id can't be empty")
	}

	c.Project = v.GetString("PROJECT")
	if c.Project == "" {
		logger.Error("The project name can't be empty")
		return errors.New("The project name can't be empty")
	}

	c.Endpoint = v.GetString("ENDPOINT")
	if c.Endpoint == "" {
		logger.Error("The endpoint can't be empty")
		return errors.New("The endpoint can't be empty")
	}

	c.Instance = v.GetString("INSTANCE")
	if c.Instance == "" {
		logger.Error("The instance can't be empty")
		return errors.New("The instance can't be empty")
	}

	lookBack := v.GetInt32("MAX_LOOK_BACK")
	if lookBack == 0 {
		lookBack = DefaultLookBack
	}

	if lookBack > 3*24 {
		logger.Warn("Setting a larger value for MAX_LOOK_BACK will affect the query efficiency.", "MAX_LOOK_BACK", lookBack)
	}

	c.MaxLookBack = time.Duration(lookBack) * time.Hour
	logger.Info("Parameters", "AccessSecret", c.AccessSecret, "AccessKeyID", c.AccessKeyID, "Project", c.Project, "Instance", c.Instance, "Endpoint", c.Endpoint, "MaxLookBack", c.MaxLookBack)
	return nil
}
