package main

import (
	"errors"
	"flag"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jaegertracing/jaeger/plugin/storage/grpc"
	"github.com/jaegertracing/jaeger/plugin/storage/grpc/shared"
	"github.com/qiansheng91/jaeger-sls/sls_store"
	"github.com/spf13/viper"
)

const (
	DefaultLookBack = 7
)

var configPath string

type Configuration struct {
	Endpoint     string                 `yaml:"endpoint"`
	AccessKeyID  string                 `yaml:"accessKeyId"`
	AccessSecret string                 `yaml:"accessSecret"`
	Project      string                 `yaml:"project"`
	Instance     string                 `yaml:"instance"`
	Optional     *OptionalConfiguration `yaml:"optionals"`
}

type OptionalConfiguration struct {
	MaxLookBack int64 `yaml:maxlookback`
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Info,
		Name:       "aliyun-log-jaeger-plugin",
		JSONFormat: true,
	})

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
		7*24*time.Hour,
		logger,
	)

	grpc.Serve(&shared.PluginServices{
		Store: plugin,
	})

	logger.Info("SLS jaeger plugin initialized Successfully")
}

func initialParameters(configPath string, logger hclog.Logger) (*Configuration, error) {
	v := viper.New()
	v.SetEnvPrefix("sls_storage")
	v.AutomaticEnv()
	if configPath != "" {
		v.SetConfigFile(configPath)
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	configuration := &Configuration{}
	if err := configuration.InitFromViper(v); err != nil {
		return nil, err
	} else {
		return configuration, nil
	}
}

func (c *Configuration) InitFromViper(v *viper.Viper) error {
	c.AccessSecret = v.GetString("sls_storage.access_secret")
	if c.AccessSecret == "" {
		return errors.New("The AccessSecret can't be empty")
	}

	c.AccessKeyID = v.GetString("sls_storage.access_key")

	if c.AccessKeyID == "" {
		return errors.New("The access key id can't be empty")
	}

	c.Project = v.GetString("sls_storage.project")
	if c.Project == "" {
		return errors.New("The project name can't be empty")
	}

	c.Endpoint = v.GetString("sls_storage.endpoint")
	if c.Endpoint == "" {
		return errors.New("The endpoint can't be empty")
	}

	c.Instance = v.GetString("sls_storage.instance")
	if c.Instance == "" {
		return errors.New("The instance can't be empty")
	}

	c.Optional = &OptionalConfiguration{}

	c.Optional.MaxLookBack = v.GetInt64("sls_storage.optional.max_look_back")
	if c.Optional.MaxLookBack == 0 {
		c.Optional.MaxLookBack = DefaultLookBack
	}

	return nil
}
