package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {

	cfgPath, backend, err := ParseFlags()
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

	if !backend {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				t, _ := template.ParseFiles("/www/index.html")
				err = t.Execute(w, nil)
				if err != nil {
					log.Fatal(err)
				}
			}
		})
	}

	http.HandleFunc("/changeState", r.stateHandler(cfg, c))

	log.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
