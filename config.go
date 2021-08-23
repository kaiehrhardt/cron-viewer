package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"
)

type Config []struct {
	Endpoints []string `yaml:"endpoints"`
	Cronjobs  []string `yaml:"cronjobs"`
}

func NewConfig(configPath string) (Config, error) {
	config := Config{}

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

func (cfg Config) CallEndpoints() {
	for _, entry := range cfg {
		for _, endpoint := range entry.Endpoints {
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
}

func (cfg Config) StartCron() (c *cron.Cron) {
	c = cron.New()
	for _, entry := range cfg {
		for _, job := range entry.Cronjobs {
			_, err := c.AddFunc(job, cfg.CallEndpoints)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	c.Start()

	log.Println(c.Entries())

	return c
}

func (cfg Config) StopCron(c *cron.Cron) {
	if len(c.Entries()) >= 1 {
		c.Stop()
	}
}
