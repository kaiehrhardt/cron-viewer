package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/robfig/cron/v3"
)

func toggle() {
	resp, err := http.Get("http://api.iobroker.homelab.local/toggle/0_userdata.0.test")
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

func main() {
	c := cron.New()
	c.AddFunc("0 0 4 * * 0-5", toggle)
	c.AddFunc("0 0 6 * * 0-5", toggle)
	c.AddFunc("0 0 9 * * 0-5", toggle)
	c.AddFunc("0 0 22 * * 0-5", toggle)
	c.Start()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
