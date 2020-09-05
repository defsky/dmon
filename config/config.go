package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/file"
)

const defaultConfigRootNode string = "app"

var (
	appConfig ApplicationConfig
)

type redisConfig struct {
	Host string
	Port string
	Keys map[string]string
}

type mysqlConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	Charset  string
	Encrypt  string
}

type mssqlConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	Charset  string
	Encrypt  string
}

type databaseConfig struct {
	Redis *redisConfig
	Mssql map[string]mssqlConfig
	Mysql map[string]mysqlConfig
}

// ApplicationConfig ...
type ApplicationConfig struct {
	DB     *databaseConfig
	Target string
}

// GetConfig ...
func GetConfig() *ApplicationConfig {
	return &appConfig
}

// Init ...
func Init() {
	cfgFile, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	cfgFile = filepath.Join(cfgFile, "config.yml")

	log.Printf("Load config file: %s\n", cfgFile)

	fs := file.NewSource(file.WithPath(cfgFile))
	err = config.Load(fs)

	if err != nil {
		panic(err)
	}

	appConfig.Target = config.Get(defaultConfigRootNode, "target").String("")

	dbTypes := config.Get(defaultConfigRootNode, "database").StringMap(map[string]string{})
	if len(dbTypes) <= 0 {
		return
	}
	appConfig.DB = &databaseConfig{}

	for dbType := range dbTypes {
		switch dbType {
		case "mssql":
			appConfig.DB.Mssql = make(map[string]mssqlConfig)
			cfgs := config.Get(defaultConfigRootNode, "database", "mssql").StringMap(map[string]string{})
			for k := range cfgs {
				cfg := &mssqlConfig{}
				if err := config.Get(defaultConfigRootNode, "database", "mssql", k).Scan(cfg); err != nil {
					panic(err)
				}
				appConfig.DB.Mssql[k] = *cfg
			}
		case "mysql":
			appConfig.DB.Mysql = make(map[string]mysqlConfig)
			cfgs := config.Get(defaultConfigRootNode, "database", "mysql").StringMap(map[string]string{})
			for k := range cfgs {
				cfg := &mysqlConfig{}
				if err := config.Get(defaultConfigRootNode, "database", "mysql", k).Scan(cfg); err != nil {
					panic(err)
				}
				appConfig.DB.Mysql[k] = *cfg
			}
		case "redis":
			appConfig.DB.Redis = &redisConfig{}
			if err := config.Get(defaultConfigRootNode, "database", "redis").Scan(appConfig.DB.Redis); err != nil {
				panic(err)
			}
		}
	}
}
