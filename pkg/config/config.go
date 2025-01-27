package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

/*
	type Config struct {
		ServerConfig  ServerConfig  `json:"server" yaml:"server"`
		DBConfig      DBConfig      `json:"db"`
		LoggingConfig LoggingConfig `json:"logging" yaml:"logging"`
		TraceConfig   TraceConfig   `json:"trace" yaml:"trace"`
	}

	type LoggingConfig struct {
		Level           int    `json:"level"`
		Encoding        string `json:"encoding"`
		Development     bool   `json:"development"`
		InfoFilename    string `json:"info_filename"`
		InfoMaxSize     int    `json:"info_max_size"`
		InfoMaxBackups  int    `json:"info_max_backups"`
		InfoMaxAge      int    `json:"info_max_age"`
		InfoCompress    bool   `json:"info_compress"`
		ErrorFilename   string `json:"error_filename"`
		ErrorMaxSize    int    `json:"error_max_size"`
		ErrorMaxBackups int    `json:"error_max_backups"`
		ErrorMaxAge     int    `json:"error_max_age"`
		ErrorCompress   bool   `json:"error_compress"`
	}

	type ServerConfig struct {
		Port             int           `json:"port"`
		ReadTimeout      time.Duration `json:"readTimeout"`
		WriteTimeout     time.Duration `json:"writeTimeout"`
		GracefulShutdown time.Duration `json:"gracefulShutdown"`
		Host             string        `json:"host"`
	}

	type DBConfig struct {
		DataSourceName string `json:"dataSourceName"`
		LogLevel       int    `json:"logLevel"`
		Pool           struct {
			MaxOpen     int           `json:"maxOpen"`
			MaxIdle     int           `json:"maxIdle"`
			MaxLifetime time.Duration `json:"maxLifetime"`
		} `json:"pool"`
	}

	type TraceConfig struct {
		IsTraceEnabled    bool   `json:"is_enabled"`
		Url               string `json:"trace_url"`
		ServiceName       string `json:"trace_service_name"`
		IsHttpBodyEnabled bool   `json:"trace_is_http_body_enabled"`
	}

	func Load() (*Config, error) {
		k := koanf.New(".")

		path, err := filepath.Abs("./config/config.local.yaml")
		if err != nil {
			log.Printf("failed to get absoulute config path. configPath:%s, err: %v", "./config/config.local.yaml", err)
			return nil, err
		}

		log.Printf("load config file from %s", path)
		if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
			log.Printf("failed to load config from file. err: %v", err)
			return nil, err
		}

		var cfg Config
		if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "json", FlatPaths: false}); err != nil {
			log.Printf("failed to unmarshal with conf. err: %v", err)
			return nil, err
		}
		return &cfg, err
	}
*/
type Config struct {
	DBHost            string
	DBPort            int
	DBUser            string
	DBPassword        string
	DBName            string
	ServeHost         string
	ServerPort        int
	ServerReadTimeout int

	LogLevel           int
	LogEncoding        string
	LogInfoFilename    string
	LogInfoMaxSize     int
	LogInfoMaxBackups  int
	LogInfoMaxAge      int
	LogInfoCompress    bool
	LogErrorFilename   string
	LogErrorMaxSize    int
	LogErrorMaxBackups int
	LogErrorMaxAge     int
	LogErrorCompress   bool
	DataSourceName     string
	Listen             string
	SecretKey          string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error loading config file: %v", err)
		return nil, err
	}

	return &Config{
		DBHost:            viper.GetString("DB_HOST"),
		DBPort:            viper.GetInt("DB_PORT"),
		DBUser:            viper.GetString("DB_USER"),
		DBPassword:        viper.GetString("DB_PASSWORD"),
		DBName:            viper.GetString("DB_NAME"),
		ServeHost:         viper.GetString("SERVER_HOST"),
		ServerPort:        viper.GetInt("SERVER_PORT"),
		ServerReadTimeout: viper.GetInt("SERVER_READ_TIMEOUT"),
		LogLevel:          viper.GetInt("level"),
		LogEncoding:       viper.GetString("encoding"),
		LogInfoFilename:   viper.GetString("info_filename"),
		LogInfoMaxSize:    viper.GetInt("info_max_size"),
		LogInfoMaxBackups: viper.GetInt("info_max_backups"),
		LogInfoMaxAge:     viper.GetInt("info_max_age"),
		LogInfoCompress:   viper.GetBool("info_compress"),
		LogErrorFilename:  viper.GetString("error_filename"),

		LogErrorMaxSize:    viper.GetInt("error_max_size"),
		LogErrorMaxBackups: viper.GetInt("error_max_backups"),
		LogErrorMaxAge:     viper.GetInt("error_max_age"),
		LogErrorCompress:   viper.GetBool("error_compress"),
		DataSourceName:     fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Europe/Moscow", viper.GetString("SERVER_HOST"), viper.GetString("DB_USER"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"), viper.GetInt("DB_PORT")),
		Listen:             fmt.Sprintf("%v:%v", viper.GetString("SERVER_HOST"), viper.GetInt("SERVER_PORT")),
		SecretKey:          viper.GetString("SECRET_KEY"),
	}, nil
}
