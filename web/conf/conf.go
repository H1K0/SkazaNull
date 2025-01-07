package conf

import (
	"os"

	"gopkg.in/yaml.v3"
)

const CONF_FILE = "/etc/skazanull/web.conf.yml"

type Conf struct {
	DBConfig      string `yaml:"db_config"`
	Interface     string `yaml:"interface"`
	EncryptionKey string `yaml:"encryption_key"`
}

func GetConf() (conf Conf, err error) {
	confFile, err := os.ReadFile(CONF_FILE)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(confFile, &conf)
	return
}
