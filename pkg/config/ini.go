package config

import (
	"gopkg.in/ini.v1"
)

type Config struct {
	value *ini.File
}

var config *Config

func NewConfig(filename string) *Config {
	cfg, err := ini.Load(filename)

	if err != nil {
		panic("加载配置文件失败" + filename + ", 原因: " + err.Error())
	}
	config = &Config{
		value: cfg,
	}

	return config
}

func Get(section, key string) *ini.Key {
	return config.value.Section(section).Key(key)
}

func GetUrlByCorpid(corpid string, productid string, key string) string {
	if corpid == "" {
		corpid = "0"
	}
	if productid == "" {
		productid = "0"
	}
	value := config.value.Section("").Key(key).String()
	urlPrefix := config.value.Section("").Key("URL_PREFIX").String()
	return urlPrefix + corpid + "/" + productid + value
}
