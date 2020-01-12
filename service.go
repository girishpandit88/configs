package main

import (
	"github.com/girishpandit88/configs/postgres"
	"log"
	"os"
)

func main() {
	env, present := os.LookupEnv("DB_URL")
	if !present {
		os.Exit(1)
	}
	dba, err := postgres.NewDBA(env)
	if err != nil {
		log.Fatal(err)
	}
	//simple config
	err = dba.Save("config", &map[string]interface{}{"url": "http://apple.com", "address": "1 Apple Park"})
	if err != nil {
		log.Fatal(err)
	}
	//simple config
	err = dba.Save("config", &map[string]interface{}{"url": "http://google.com", "address": "1600 Amphitheatre Parkway"})
	if err != nil {
		log.Fatal(err)
	}
	//complex config
	err = dba.Save("db", &map[string]interface{}{
		"url":  "localhost",
		"port": 5432,
		"config": map[string]interface{}{
			"maxConnection":  10,
			"batchSize":      1000,
			"connectionType": "pgxpool.Conn",
		},
		"driverContext": "postgres",
	})

	getObj, err := dba.GetConfig("db")
	if err != nil {
		log.Fatalf("No matching config found for key: db %v", err)
	}
	if getObj != nil {
		log.Println(getObj)
	}
	matchingRows, err := dba.GetConfigByProperty("url", "e")
	if err != nil {
		log.Fatalf("No matching config found for key: google.com %s", err)
	}
	log.Println("Matching rows")
	for _, o := range *matchingRows {
		log.Println(o.String())
	}
}
