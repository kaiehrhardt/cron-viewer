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

func (r *Runner) stateHandler(cfg Config, cs []*cron.Cron) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" && req.Method != "GET" {
			http.Error(w, "Method is not supported.", http.StatusNotFound)
			return
		}

		if req.Method == "GET" {
			_, err := io.WriteString(w, strconv.FormatBool(r.Active))
			if err != nil {
				log.Fatal(err)
			}
		} else if req.Method == "POST" {
			r.Switch()
			log.Printf("New Runner State - Active: %+v", r.Active)

			if len(cs) >= 0 && r.Active {
				for i, _ := range cs {
					cfg.StartCron(i)
				}
			} else if !r.Active {
				for _, c := range cs {
					cfg.StopCron(c)
				}
			}
		}
	}
}
