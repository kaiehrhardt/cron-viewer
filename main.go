package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Endpoints []string `yaml:"endpoints"`
	Cronjobs  []string `yaml:"cronjobs"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func ParseFlags() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")
	flag.Parse()

	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	return configPath, nil
}

func (cfg *Config) Toggle() {
	for _, endpoint := range cfg.Endpoints {
		resp, err := http.Get(endpoint)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		bodyString := string(body)
		log.Println(bodyString)
	}
}

func (cfg *Config) Run() {
	c := cron.New()
	for _, job := range cfg.Cronjobs {
		_, err := c.AddFunc(job, cfg.Toggle)
		if err != nil {
			log.Fatal(err)
		}
	}
	c.Start()

	log.Println(c.Entries())

	select {}
}

func main() {
	cfgPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config Path: ", cfgPath)

	cfg, err := NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config: ", cfg)

	cfg.Run()
}
