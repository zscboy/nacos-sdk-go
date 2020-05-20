package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	huoysViper *viper.Viper
)

// HuoysConfigParam 配置参数
type HuoysConfigParam struct {
	DataID         string `json:"dataid"`
	Group          string `json:"group"`
	ConfigFilePath string `json:"configFilePath"`
}

// LoadConfig 加载配置
func LoadConfig(configPath string) {
	viper, err := loadServerConfig(configPath)
	if err != nil {
		log.Panic("LoadConfig error:", err)
	}

	huoysViper = viper
}

func loadServerConfig(configPath string) (*viper.Viper, error) {
	log.Println("config dir path:", configPath)
	conf := viper.New()

	conf.SetConfigFile(configPath)
	// conf.AddConfigPath(configPath)
	err := conf.ReadInConfig()
	if err != nil {
		fmt.Println("failed to read config,err:", err)
		return nil, err
	}

	//配置文件修改监听
	conf.WatchConfig()
	conf.OnConfigChange(func(event fsnotify.Event) {
		err1 := conf.ReadInConfig()
		if err1 != nil {
			log.Panic("server config was changed,err:", err1)
		}
		log.Info("config change")
	})

	return conf, nil
}

// GetConfigServerIP 获取配置中心的ip
func GetConfigServerIP() string {
	return huoysViper.GetString("server.ip")
}

// GetConfigServerPort 获取配置中心的端口
func GetConfigServerPort() uint64 {
	return huoysViper.GetUint64("server.port")
}

// GetConfigServerContentPath 获取配置中心的contentPath
func GetConfigServerContentPath() string {
	return huoysViper.GetString("server.contentPath")
}

// GetConfigServerConfigParams 获取配置参数
func GetConfigServerConfigParams() []HuoysConfigParam {
	var cfgs []HuoysConfigParam

	huoysViper.UnmarshalKey("configs", &cfgs)

	return cfgs
}
