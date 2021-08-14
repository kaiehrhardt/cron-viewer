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

func (cfg *Config) Start() (c *cron.Cron) {
	c = cron.New()
	for _, job := range cfg.Cronjobs {
		_, err := c.AddFunc(job, cfg.Toggle)
		if err != nil {
			log.Fatal(err)
		}
	}
	c.Start()

	log.Println(c.Entries())

	return c
}

func (cfg *Config) Stop(c *cron.Cron) {
	if len(c.Entries()) >= 1 {
		c.Stop()
	}
}

type Runner struct {
	Active bool
}

func NewRunner() *Runner {
	return &Runner{
		Active: true,
	}
}

func (r *Runner) Switch() {
	if r.Active {
		r.Active = false
	} else {
		r.Active = true
	}
}

func (r *Runner) stateHandler(cfg *Config, c *cron.Cron) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/changeState" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}

		if req.Method != "GET" {
			http.Error(w, "Method is not supported.", http.StatusNotFound)
			return
		}

		r.Switch()
		log.Printf("New Runner State - Active: %+v", r.Active)

		if c != nil && r.Active {
			c = cfg.Start()
		} else if !r.Active {
			cfg.Stop(c)
		}
	}
}

func main() {

	cfgPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initial Prepair - Config Path: ", cfgPath)

	cfg, err := NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initial Prepair - Config: ", cfg)

	r := NewRunner()
	log.Printf("Initial Prepair - Runner State Active: %+v", r.Active)

	c := cfg.Start()

	http.HandleFunc("/changeState", r.stateHandler(cfg, c))

	log.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
