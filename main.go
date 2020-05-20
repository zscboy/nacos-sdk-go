package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/nacos_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/http_agent"
	"github.com/nacos-group/nacos-sdk-go/config"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var clientConfigTest = constant.ClientConfig{
	TimeoutMs:           10 * 1000,
	BeatInterval:        5 * 1000,
	ListenInterval:      300 * 1000,
	NotLoadCacheAtStart: true,
	// Username:            "nacos",
	// Password:            "nacos",
}

// var serverConfig = constant.ServerConfig{
// 	IpAddr:      "localhost",
// 	Port:        8848,
// 	ContextPath: "/nacos",
// }

func cretateConfigClient() config_client.ConfigClient {
	serverConfig := constant.ServerConfig{
		IpAddr:      config.GetConfigServerIP(),
		Port:        config.GetConfigServerPort(),
		ContextPath: config.GetConfigServerContentPath(),
	}

	nc := nacos_client.NacosClient{}
	nc.SetServerConfig([]constant.ServerConfig{serverConfig})
	nc.SetClientConfig(clientConfigTest)
	nc.SetHttpAgent(&http_agent.HttpAgent{})
	client, err := config_client.NewConfigClient(&nc)
	if err != nil {
		log.Println("cretateConfigClient error:", err)
	}
	return client
}

var (
	cfgFilepath = ""
)

func init() {
	flag.StringVar(&cfgFilepath, "c", "", "specify the config file path name")
}

func getConfigFilePath(dataid string) string {
	configs := config.GetConfigServerConfigParams()
	for _, config := range configs {
		if config.DataID == dataid {
			return config.ConfigFilePath
		}
	}

	return ""
}

func main() {
	flag.Parse()

	if len(cfgFilepath) == 0 {
		log.Println("Need config file path")
		return
	}

	config.LoadConfig(cfgFilepath)

	log.Println("config server ip:", config.GetConfigServerIP())

	client := cretateConfigClient()

	configs := config.GetConfigServerConfigParams()
	for _, config := range configs {
		client.ListenConfig(vo.ConfigParam{
			DataId: config.DataID,
			Group:  config.Group,
			OnChange: func(namespace, group, dataId, data string) {
				configFilePath := getConfigFilePath(dataId)
				err := ioutil.WriteFile(configFilePath, []byte(data), 0644)
				if err != nil {
					panic(err)
				}
				fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", data:" + data)
				// fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", data:" + data + " configFilePath:" + configFilePath)
			},
		})
	}

	select {}

}
