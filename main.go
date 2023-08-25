package main

import (
	"flag"
	"fmt"
	"github.com/bootcamp-go/desafio-go-bases/internal/tickets"
	"log"
	"sync"
	"time"
)

const (
	earlyMorning = "0:6"
	morning      = "7:12"
	afternoon    = "13:19"
	evening      = "20:23"
)

var times = [4]string{earlyMorning, morning, afternoon, evening}

func main() {

	start := time.Now()
	ticketPtr := flag.String("input", "", "The ticket file to process")
	destination := flag.String("destination", "Brazil", "Destination total tickets")
	totalPeople := flag.Int("total", 1000, "Total people")

	flag.Parse()

	stats := tickets.NewStats()
	err := stats.LoadTicketsByCSV(*ticketPtr)
	if err != nil {
		log.Fatal("Tickets cannot be loaded")
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	chErr := make(chan error)

	go func(destination string, cErr chan<- error) {
		defer wg.Done()
		total, fnErr := stats.GetTotalTickets(destination)
		if fnErr != nil {
			cErr <- fnErr
		}
		fmt.Println(fmt.Sprintf("Total tickets sell in %s %d", destination, total))
	}(*destination, chErr)

	go func(destination string, total int, cErr chan<- error) {
		defer wg.Done()

		avg, fnErr := stats.AverageDestination(destination, total)
		if fnErr != nil {
			cErr <- fnErr
		}

		fmt.Println(fmt.Sprintf("Average people traveling to %s of a total of %d in comparisson with other countries %.2f%%", destination, total, avg))
	}(*destination, *totalPeople, chErr)

	for _, t := range times {
		wg.Add(1)
		go func(time string, cErr chan<- error) {
			wg.Done()
			ticketsByTime, fnErr := stats.GetMornings(time)
			if fnErr != nil {
				cErr <- fnErr
			}
			fmt.Println(fmt.Sprintf("Total tickets by time %s - %d", time, ticketsByTime))
		}(t, chErr)
	}
	log.Println("Waiting for results")

	go func() {
		for fnErr := range chErr {
			log.Println(fmt.Sprintf("An error has be occurred %s", fnErr))
		}
	}()

	wg.Wait()
	close(chErr)

	end := time.Now()
	fmt.Println("start -> ", start.Format("2006-01-02T15:04:05.000Z"))
	fmt.Println("end -> " + end.Format("2006-01-02T15:04:05.000Z"))
	elapsed := end.Sub(start)
	fmt.Println("elapsed time: ", elapsed)
}
