package config

import (
	"demo-curd/util/constant"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Env      string
	Database struct {
		Host            string `yaml:"host"`
		Port            string `yaml:"port"`
		Username        string `yaml:"username"`
		Password        string `yaml:"password"`
		Dbname          string `yaml:"dbname"`
		MaxIdleConns    int    `yaml:"maxIdleConns"`
		MaxOpenConns    int    `yaml:"maxOpenConns"`
		ConnMaxLifetime string `yaml:"connMaxLifetime"`
	} `yaml:"database"`

	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	RabbitMQ struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Producer *struct {
			Exchange string `yaml:"exchange"`
			Bindings map[string]struct {
				Queue      string  `yaml:"queue"`
				RoutingKey *string `yaml:"routingKey,omitempty"`
			} `yaml:"bindings"`
		} `yaml:"producer,omitempty"`
		Consumer *struct {
			Queue string `yaml:"queue"`
		} `yaml:"consumer,omitempty"`
	} `yaml:"rabbitmq"`

	Jwt struct {
		Realm              string `yaml:"realm"`
		SigningAlg         string `yaml:"signAlg"`
		Secret             string `yaml:"secret"`
		ExpiredTime        string `yaml:"expiredTime"`
		RefreshExpTime     string `yaml:"refreshExpTime"`
		LongRefreshExpTime string `yaml:"longRefreshExpTime"`
	} `yaml:"jwt"`

	I18n struct {
		Langs []string `yaml:"langs"`
	} `yaml:"i18n"`

	CORS struct {
		AllowOrigins     []string `yaml:"allowOrigins"`
		AllowMethods     []string `yaml:"allowMethods"`
		AllowHeaders     []string `yaml:"allowHeaders"`
		ExposeHeaders    []string `yaml:"exposeHeaders"`
		AllowCredentials bool     `yaml:"allowCredentials"`
		MaxAge           string   `yaml:"maxAge"`
	} `yaml:"cors"`

	HostUrl map[string]string `yaml:"hostUrl"`

	Security struct {
		AuthorizedRequests []ConfigAuthorizedRequests `yaml:"authorizedRequests"`
	} `yaml:"security"`

	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`

	Swagger struct {
		Url string `yaml:"url"`
	} `yaml:"swagger"`
}

type ConfigAuthorizedRequests struct {
	Urls        []string                `yaml:"urls"`
	Access      constant.SecurityAccess `yaml:"access"`
	Roles       []string                `yaml:"roles"`
	Permissions []string                `yaml:"permissions"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (c Config, err error) {
	env := extractEnv()
	defer func() {
		if err != nil {
			c = Config{
				Env: env,
			}
		}
	}()
	// get current path
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	path, err := getAbsPath(pwd)
	if err != nil {
		return
	}

	// load config from config directory
	if path == "/" {
		viper.AddConfigPath("/config")
	} else {
		viper.AddConfigPath(fmt.Sprintf("%v/config", path))
	}
	viper.SetConfigName(fmt.Sprintf("app-%v", strings.ToLower(env)))
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(&c); err != nil {
		return
	}
	return c, nil
}

func extractEnv() string {
	env := os.Getenv("ENVIRONMENT")
	if len(env) == 0 {
		env = os.Getenv("ENV")
	}
	if len(env) == 0 {
		env = constant.DefaultEnv
	}
	return env
}

func getAbsPath(dir string) (string, error) {
	path, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}
	return path, nil
}
