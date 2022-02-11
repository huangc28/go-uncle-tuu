package config

import (
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type AppConf struct {
	APIPort      uint   `mapstructure:"API_PORT"`
	APIJWTSecret string `mapstructure:"API_JWT_SECRET"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     uint   `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBDbname   string `mapstructure:"DB_NAME"`
	DBTimeZone string `mapstructure:"DB_TIMEZONE"`

	ImportFailedFileDirPath string `mapstructure:"IMPORT_FAILED_FILE_DIR_PATH"`
}

var appConf AppConf

// GetProjRootPath gets project root directory relative to `config/config.go`
func GetProjRootPath() string {
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(b)
	)

	log.Printf("project root dir %v", basepath)

	return filepath.Join(basepath, "..")
}

// List of environments.
// Note: not in use at the moment
type Env string

var (
	Production  Env = "production"
	Staging     Env = "staging"
	Development Env = "development"
	Test        Env = "test"
)

// reads config from .env
func InitConfig() {
	viper.SetConfigType("env")

	viper.SetConfigName(".env")

	// Config path can be at project root directory
	cwd, err := os.Getwd()
	// retrieve executable path

	if err != nil {
		log.Fatalf("failed to get current working directory: %s", err.Error())
	}

	log.Infof("search .env in path... %s", cwd)

	viper.AddConfigPath(cwd)
	viper.AddConfigPath(".")
	viper.AddConfigPath(GetProjRootPath())
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("configuration file not found: %s", err.Error())
		} else {
			log.Fatalf("error occurred when read in config file: %s", err.Error())
		}
	}

	if err = viper.Unmarshal(&appConf); err != nil {
		log.Fatalf("failed to unmarshal app config to struct %s", err.Error())
	}
}

func GetAppConf() *AppConf {
	return &appConf
}
