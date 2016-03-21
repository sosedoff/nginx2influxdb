package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	period int
	batch  int
	stream bool
	dbUrl  string
	dbName string
)

func init() {
	flag.IntVar(&period, "p", 5, "Number of seconds between writes")
	flag.StringVar(&dbUrl, "h", "", "InfluxDB server url")
	flag.StringVar(&dbName, "d", "", "InfluxDB database name")
	flag.BoolVar(&stream, "s", false, "Stream")
	flag.IntVar(&batch, "b", 1000, "Batch size")
	flag.Parse()

	if dbUrl == "" {
		fmt.Println("Server url is not provided")
		os.Exit(1)
	}

	if dbName == "" {
		fmt.Println("Database name is not provided")
		os.Exit(1)
	}
}

func main() {
	db := NewDatabase(dbUrl, dbName)
	scanner := bufio.NewScanner(os.Stdin)
	requests := []Request{}

	// Write to influx db periodically only when streaming is enabled
	if stream {
		go func() {
			for {
				if len(requests) > 0 {
					go db.Write(requests)
					requests = []Request{}
				}

				time.Sleep(time.Duration(period) * time.Second)
			}
		}()
	}

	for scanner.Scan() {
		line := scanner.Text()

		req, err := NewRequest(line)
		if err != nil {
			log.Println(err)
			continue
		}

		requests = append(requests, req)

		if stream {
			log.Println(req.InfluxString())
		}
	}

	if !stream && len(requests) > 0 {
		num := len(requests)
		pages := num / batch
		offset := 0

		for i := 0; i < pages; i++ {
			db.Write(requests[offset : offset+batch])
			offset += batch
		}

		if offset < num {
			db.Write(requests[offset:num])
		}
	}
}
