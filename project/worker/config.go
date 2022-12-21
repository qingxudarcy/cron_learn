package worker

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	EtcdEndpoints []string `json:"etcdEndpoints"`
	EtcdDialTimeout int `json:"etcdEtcdDialTimeout"`
	MongoUri string `json:"mongoUri"`
}

var (
	G_config *Config
)

func InitConfig(fileName string) (err error) {
	var (
		content []byte
		conf Config
	)
	if content, err = ioutil.ReadFile(fileName); err != nil {
		return
	}

	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	G_config = &conf

	return

}