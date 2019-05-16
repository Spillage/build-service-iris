package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/heiing/logs"
)

var config *Config = loadDefaultConfig()

type MysqlConfig struct {
	Conn string `json:"conn"`
}

type BuilderTemplateConfig struct {
	TypeId string `json:"typeId"`
	Path   string `json:"path"`
}

type CIConfig struct {
	CISZ string `json:"cisz"`
	CISH string `json:"cish"`
}

type Config struct {
	Logs                   *logs.LogsConfig `json:"logs"`
	Mysql                  *MysqlConfig     `json:"mysql"`
	CI                     *CIConfig        `json:"ci"`
	BuilderTemplateConfigs map[string]*BuilderTemplateConfig
}

func loadDefaultConfig() *Config {
	currentPath := logs.GetExecPath()
	configFile := filepath.Join(currentPath, "config.json")

	buf, err := ioutil.ReadFile(configFile)
	if nil != err {
		logs.Error("Load config file Faild [", configFile, "]: ", err)
		return nil
	}

	config := &Config{}
	if err = json.Unmarshal(buf, config); nil != err {
		logs.Error("Parse config file Faild [", configFile, "]: ", err)
		return nil
	}

	logs.SetDefaultLoggerForConfig(config.Logs)

	tcp := filepath.Join(currentPath, "config-type-map.json")
	tcf, err := os.Open(tcp)
	if nil != err {
		logs.Error("Open type map config file Faild [", tcp, "]: ", err)
		return config
	}

	tcs := make([]*BuilderTemplateConfig, 0)
	if err := json.NewDecoder(tcf).Decode(&tcs); nil != err {
		logs.Error("Parse type map config Faild [", tcp, "]: ", err)
		return config
	}

	config.BuilderTemplateConfigs = make(map[string]*BuilderTemplateConfig)
	for _, item := range tcs {
		config.BuilderTemplateConfigs[item.TypeId] = item
	}

	return config
}
