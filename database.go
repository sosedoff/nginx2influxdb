package main

import (
	"log"
	"net/http"
	"strings"
)

type Database struct {
	Url string
}

func NewDatabase(url, name string) Database {
	return Database{url + "/write?db=" + name}
}

func (db Database) Write(requests []Request) error {
	log.Println("writing", len(requests), "points")

	lines := make([]string, len(requests))
	for i, req := range requests {
		lines[i] = req.InfluxString()
	}

	reader := strings.NewReader(strings.Join(lines, "\n"))

	resp, err := http.Post(db.Url, "application/x-www-form-urlencoded", reader)
	if err == nil {
		resp.Body.Close()
	}

	log.Println("got response:", resp.Status)
	return err
}
