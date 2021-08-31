package main

import (
	"fmt"
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Content-Type", "application/json")

		if req.Method != "POST" && req.Method != "GET" && req.Method != "OPTIONS" {
			http.Error(w, "Method is not supported.", http.StatusNotFound)
			return
		}

		if req.Method == "POST" {
			r.Switch()

			log.Printf("New Runner State - Active: %+v", r.Active)

			if c != nil && r.Active {
				c = cfg.Start()
			} else if !r.Active {
				cfg.Stop(c)
			}
		}

		w.WriteHeader(http.StatusOK)
		data := fmt.Sprintf("{\"Active\":%t}\n", r.Active)
		response := []byte(data)
		_, err := w.Write(response)
		if err != nil {
			log.Fatal(err)
		}
	}
}
