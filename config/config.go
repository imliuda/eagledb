package config

import (
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
)

var G = Config{}

type HttpConfig struct {
	ListenAddr string `toml:"listen_addr"`
	ListenPort int    `toml:"listen_port"`
}

type Config struct {
	Http HttpConfig `toml:"http"`
}

func LoadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(data, &G)
	if err != nil {
		return err
	}

	return nil
}
