package main

import (
	"log"
	"net/http"

	"github.com/robfig/cron/v3"
)

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
		if req.Method != "POST" {
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
