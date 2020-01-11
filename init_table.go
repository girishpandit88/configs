package main

import (
	"flag"
	"github.com/girishpandit88/configs/postgres"
	"log"
	"os"
)

func main() {
	var url string
	flag.StringVar(&url, "db", "", "DB url for string")
	flag.Parse()
	if url == "" {
		log.Fatal("Usage: go run init_table -db <>")
	}
	env, present := os.LookupEnv("DB_URL")
	if !present {
		log.Println("DB url not set")
		os.Exit(1)
	}
	dba, err := postgres.NewDBA(env)
	if err != nil {
		log.Fatal(err)
	}
	err = dba.CreateTableFromFile("configs.sql")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table and indexes created!")
}
