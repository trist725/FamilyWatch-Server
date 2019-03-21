package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var (
	Conf struct {
		RunSpider bool     `yaml:"RunSpider"`
		FilteMin  int      `yaml:"FilteMin"`
		MongoURI  string   `yaml:"MongoURI"`
		Qq        []string `yaml:"qq"`
		Iqiyi     []string `yaml:"iqiyi"`
		Youku     []string `yaml:"youku"`
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
