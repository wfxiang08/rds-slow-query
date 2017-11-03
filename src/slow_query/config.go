package slow_query

import (
	"github.com/BurntSushi/toml"
	"github.com/wfxiang08/cyutils/utils/errors"
	"io/ioutil"
)

type DatabaseConfig struct {
	DatabasesOnline  []string `toml:"dbs_online"`
	DatabasesOffline []string `toml:"dbs_offline"`
	AwsKey           string   `toml:"aws_key"`
	AwsSecret        string   `toml:"aws_secret"`
	AwsRegion        string   `toml:"aws_region"`
	EmailSender      string   `toml:"email_sender"`
	EmailReceivers   []string `toml:"email_receivers"`
}

func NewConfigWithFile(name string) (*DatabaseConfig, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return NewConfig(string(data))
}

func NewConfig(data string) (*DatabaseConfig, error) {
	var c DatabaseConfig
	_, err := toml.Decode(data, &c)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &c, nil
}
