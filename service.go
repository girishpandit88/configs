package main

import (
	"github.com/girishpandit88/configs/postgres"
	"log"
	"os"
)

func main() {
	env, present := os.LookupEnv("DB_URL")
	if !present {
		log.Fatal("DB url not set")
		os.Exit(1)
	}
	dba, err := postgres.NewDBA(env)
	if err != nil {
		log.Fatal(err)
	}
	err = dba.Save("config", map[string]string{"url": "http://apple.com", "address": "1 Apple Park"})
	if err != nil {
		log.Fatal(err)
	}

	err = dba.Save("config", map[string]string{"url": "http://google.com", "address": "1600 Amphitheatre Parkway"})
	if err != nil {
		log.Fatal(err)
	}

	getObj, err := dba.GetConfig("google.com")
	if err != nil {
		log.Fatalf("No matching config found for key: google.com %s", err)
	}
	log.Println(getObj)
	matchingRows, err := dba.GetConfigByProperty("url", "e")
	if err != nil {
		log.Fatalf("No matching config found for key: google.com %s", err)
	}
	log.Println("Matching rows")
	for _, o := range *matchingRows {
		log.Println(o.String())
	}
}
