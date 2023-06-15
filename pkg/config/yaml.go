package config

import (
	"time"

	"github.com/spf13/viper"
)

var vipconfig *VipConfig

type VipConfig struct {
	*viper.Viper
}

func NewVipConfig(filename string) {
	vipcfg := viper.New()
	vipcfg.SetConfigFile(filename)
	if err := vipcfg.ReadInConfig(); err != nil {
		panic("加载配置文件失败" + filename + ", 原因: " + err.Error())
	}
	vipconfig = &VipConfig{vipcfg}
}

func VipCfg() *VipConfig {
	return vipconfig
}

func (vcfg *VipConfig) GetStringWithDefault(key, defaul string) string {
	vcfg.SetDefault(key, defaul)
	return vcfg.GetString(key)
}

func (vcfg *VipConfig) GetIntWithDefault(key string, defaul int) int {
	vcfg.SetDefault(key, defaul)
	return vcfg.GetInt(key)
}

func (vcfg *VipConfig) GetBoolWithDefault(key string, defaul bool) bool {
	vcfg.SetDefault(key, defaul)
	return vcfg.GetBool(key)
}

func (vcfg *VipConfig) GetDurationWithDefault(key string, defaul time.Duration) time.Duration {
	vcfg.SetDefault(key, defaul)
	return vcfg.GetDuration(key)
}

func (vcfg *VipConfig) GetStringSliceWithDefault(key string, defaul []string) []string {
	vcfg.SetDefault(key, defaul)
	return vcfg.GetStringSlice(key)
}
