package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var (
	Conf struct {
		RunSpider     bool   `yaml:"RunSpider"`
		FilterMin     int    `yaml:"FilterMin"`
		MaxCrawlIndex int    `yaml:"MaxCrawlIndex"`
		RefreshLimit  int    `yaml:"RefreshLimit"`
		WsAddr        string `yaml:"WsAddr"`
		MongoURI      string `yaml:"MongoURI"`
		CertFile      string `yaml:"CertFile"`
		KeyFile       string `yaml:"KeyFile"`
		Qq            map[string]interface{}
		Iqiyi         []string `yaml:"iqiyi"`
		Youku         []string `yaml:"youku"`
		Encrypt       bool     `yaml:"Encrypt"`
		FilterMaxMin  int      `yaml:"FilterMaxMin"`
	}
)

func init() {
	conf, err := ioutil.ReadFile("conf/conf.yaml")
	if err != nil {
		log.Fatalf("read conf.yaml: %v", err)
	}

	if err = yaml.Unmarshal(conf, &Conf); err != nil {
		log.Fatalf("Unmarshal conf.yaml: %v", err)
	}
}
