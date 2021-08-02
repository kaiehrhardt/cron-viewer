package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/robfig/cron/v3"
)

func toggle() {
	resp, err := http.Get("http://api.iobroker.homelab.local/toggle/hmip.0.devices.3014F711A00001DD8993E517.channels.1.on")
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

func printCronEntries(cronEntries []cron.Entry) {
	log.Printf("Cron Info: %+v\n", cronEntries)
}

func main() {
	log.Println("Starting...")
	c := cron.New()
	c.AddFunc("0 4 * * *", toggle)
	c.AddFunc("0 6 * * *", toggle)
	c.AddFunc("0 9 * * *", toggle)
	c.AddFunc("0 22 * * *", toggle)
	c.Start()
	printCronEntries(c.Entries())
	select {}
}
