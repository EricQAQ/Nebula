package config

import (
	"sync"

	"github.com/EricQAQ/Traed/err"

	"github.com/BurntSushi/toml"
	"github.com/juju/errors"
)

const (
	ConfigErrCode = 1001
)

var ConfigErr = err.CreateTraedError(ConfigErrCode, "Invalid config: %s", nil)

type TraedConfig struct {
	KlineInterval []int                      `toml:"k-line-interval" json:"k_line_interval"`
	Http          *HttpConfig                `toml:"http" json:"http"`
	ExchangeMap   map[string]*ExchangeConfig `toml:"exchange" json:"exchange"`
	Websocket     *WebsocketConfig           `toml:"websocket" json:"websocket"`
	Storage       StorageConfig              `toml:"storage" json:"storage"`

	// Log settings
	Log *LogConfig `toml:"log" json:"log"`
}

type HttpConfig struct {
	Proxy         string `toml:"proxy" json:"proxy"`
	Timeout       int    `toml:"timeout" json:"timeout"`
	RetryCount    int    `toml:"retry-count" json:"retry_count"`
	RetryInterval int    `toml:"retry-interval" json:"retry_interval"`
}

type ExchangeConfig struct {
	Address   string   `toml:"address" json:"address"`
	APIKey    string   `toml:"api-key" json:"api_key"`
	APISecret string   `toml:"api-secret" json:"api_secret"`
	Symbols   []string `toml:"symbols" json:"symbols"`
	Topic     []string `toml:"topic" json:"topic"`
	HttpUrl   string   `toml:"http-url" json:"http_url"`
}

type WebsocketConfig struct {
	HeartbeatDuration int `toml:"heartbeat-duration" json:"heartbeat_duration"`
	WriteWait         int `toml:"write-wait" json:"write_wait"`
	ReadWait          int `toml:"read-wait" json:"read_wait"`
	RetryCount        int `toml:"retry" json:"retry"`
}

type StorageConfig struct {
	StorageType string    `toml:"storage-type" json:"storage_type"`
	Csv         CsvConfig `toml:"csv" json:"csv"`
}

type CsvConfig struct {
	DataDir string `toml:"data-dir" json:"data_dir"`
}

type LogConfig struct {
	// log level
	Level string `toml:"log-level" json:"log_level"`
	// log format. One of json, text, or console
	Format string `toml:"log-format" json:"log_format"`
	// Log file
	File string `toml:"log-file" json:"log_file"`
	// Is log rotate enabled.
	LogRotate bool `toml:"log-rotate" json:"log_rotate"`
	// Max size for a single file, in MB.
	MaxSize uint `toml:"max-size" json:"max-size"`
	// Max log keep days, default is never deleting.
	MaxDays uint `toml:"max-days" json:"max-days"`
	// Maximum number of old log files to retain.
	MaxBackups uint `toml:"max-backups" json:"max-backups"`
}

var (
	once         sync.Once
	globalConfig *TraedConfig
)

func GetTraedConfig() *TraedConfig {
	once.Do(func() {
		globalConfig = new(TraedConfig)
		globalConfig.KlineInterval = []int{
			60, 180, 300, 900, 1800, 3600, 14400, 86400, 604800}
		globalConfig.Http = &HttpConfig{
			Proxy:   "",
			Timeout: 1000,
			RetryCount: 5,
			RetryInterval: 500,
		}
		globalConfig.ExchangeMap = make(map[string]*ExchangeConfig)
		globalConfig.Websocket = &WebsocketConfig{
			ReadWait:          10,
			WriteWait:         10,
			HeartbeatDuration: 10,
		}
		globalConfig.Log = &LogConfig{
			Level:      "info",
			Format:     "text",
			File:       "",
			LogRotate:  true,
			MaxSize:    500,
			MaxDays:    14,
			MaxBackups: 64,
		}
	})
	return globalConfig
}

func (c *TraedConfig) LoadFromToml(configFile string) error {
	_, err := toml.DecodeFile(configFile, c)
	if len(c.ExchangeMap) <= 0 {
		return ConfigErr.FastGen("account")
	}
	return errors.Trace(err)
}
