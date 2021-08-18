package main

import (
	"io"
	"log"
	"net/http"
	"strconv"

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
		if req.Method != "POST" && req.Method != "GET" {
			http.Error(w, "Method is not supported.", http.StatusNotFound)
			return
		}

		if req.Method == "GET" {
			io.WriteString(w, strconv.FormatBool(r.Active))
		} else if req.Method == "POST" {
			r.Switch()
			log.Printf("New Runner State - Active: %+v", r.Active)

			if c != nil && r.Active {
				c = cfg.Start()
			} else if !r.Active {
				cfg.Stop(c)
			}
		}
	}
}
