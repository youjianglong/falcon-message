package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config 配置
type Config struct {
	Addr     string   `json:"addr"`
	DingTalk DingTalk `json:"dingTalk"`

	Weixin Weixin `json:"weixin"`
	
	QCloudSms QCloudSms `json:"qcloud_sms"`
}

// Weixin 微信配置
type Weixin struct {
	Enable         bool
	CorpID         string `json:"corpID"`
	AgentID        string `json:"agentId"`
	Secret         string `json:"secret"`
	EncodingAESKey string `json:"encodingAESKey"`
}

// DingTalk 钉钉配置
type DingTalk struct {
	Enable bool `json:"enable"`
	// Level 等级， 只发送level 及其以下的消息

	Level uint `json:"level"`
}

// DingTalk 钉钉配置
type QCloudSms struct {
	Enable bool `json:"enable"`
	// Level 等级， 只发送level 及其以下的消息

	Level uint `json:"level"`

	//是否开启debug
	Debug bool `json:"debug"`

	//连接超时时间
	Timeout int `json:"timeout"`
	//短信签名
	SIGN string `json:"sign"`
	//腾讯云短信appid
	APPID string `json:"appid"`
	//腾讯云短信appkey
	APPKEY string `json:"appkey"`
	//模板ID
	TplID int `json:"tpl_id"`
	//发送方式(语音还是短信)
	Way string `json:"way"`
	//接收人
	Receivers []string `json:"receivers"`
}

// Read 读取配置
func Read() Config {
	bytes, err := ioutil.ReadFile("cfg.json")
	if err != nil {
		log.Fatalln("need file: cfg.json")
	}
	var cfg Config
	if err = json.Unmarshal(bytes, &cfg); err != nil {
		log.Fatalln("config file error", err.Error())
	}

	return cfg
}
