package tickets

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type provider interface {
	GetTotalTickets(destination string) (int, error)
	GetMornings(time string) (int, error)
	AverageDestination(destination string, total int) (int, error)
	LoadTicketsByCSV(filename string) error
}

type Stats struct {
	tickets []ticket
}

func NewStats() Stats {
	return Stats{}
}

type ticket struct {
	ID       int64
	FullName string
	Email    string
	Country  string
	Time     string
	Amount   int64
}

type interval struct {
	Start int
	End   int
}

// GetTotalTickets how many people travel to the same destiny
func (s *Stats) GetTotalTickets(destination string) (int, error) {

	dest := make(map[string]int)
	for _, t := range s.tickets {
		total, found := dest[t.Country]
		if !found {
			dest[t.Country] = 1
			continue
		}

		dest[t.Country] = total + 1
	}

	totalTickets, found := dest[destination]
	if !found {
		return 0, fmt.Errorf("destination not found")
	}
	return totalTickets, nil
}

// GetMornings how are distributed the tickets during the day
func (s *Stats) GetMornings(time string) (int, error) {
	count := 0
	timeInterval, err := parseInterval(time)
	if err != nil {
		return count, fmt.Errorf("error parsing interval %s", err)
	}
	for _, t := range s.tickets {
		ticketInterval, fnErr := parseInterval(t.Time)
		if fnErr != nil {
			return count, fnErr
		}

		if belongsToInterval(timeInterval, ticketInterval) {
			count++
		}
	}

	return count, err
}

// AverageDestination how many people travel to a particular country
func (s *Stats) AverageDestination(destination string, total int) (float64, error) {
	if total <= 0 {
		return 0, fmt.Errorf("invalid parameter total cannot be zero or less - total:%d", total)
	}

	dest := make(map[string]float64)
	for _, t := range s.tickets {
		tTickets, found := dest[t.Country]
		if !found {
			dest[t.Country] = 1
			continue
		}

		dest[t.Country] = tTickets + 1
	}

	totalTickets, found := dest[destination]
	if !found {
		return 0, fmt.Errorf("destination not found")
	}
	return totalTickets / float64(total) * 100, nil
}

// LoadTicketsByCSV es una funcion que lee un archivo
func (s *Stats) LoadTicketsByCSV(filename string) error {
	var res []ticket

	inputFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening file", err)
		return err
	}

	reader := csv.NewReader(bufio.NewReader(inputFile))
	// Skipp header
	_, err = reader.Read()
	if err != nil {
		log.Fatal("Error reading header", err)
	}
	for {
		line, rErr := reader.Read()
		if rErr == io.EOF {
			break
		} else if err != nil {
			log.Fatal(rErr)
		}

		id, fnErr := strconv.ParseInt(line[0], 10, 64)
		if fnErr != nil {
			log.Fatal("Id cannot be converted", err)
		}

		amount, fnErr := strconv.ParseInt(line[5], 10, 64)
		if fnErr != nil {
			log.Fatal("Ticket amount cannot be converted", err)
		}

		p := ticket{
			ID:       id,
			FullName: line[1],
			Email:    line[2],
			Country:  line[3],
			Time:     line[4],
			Amount:   amount,
		}

		res = append(res, p)
	}

	s.tickets = res
	/**/
	return err
}

func belongsToInterval(inside interval, outside interval) bool {
	return outside.Start >= inside.Start && outside.Start <= inside.End
}

func parseInterval(time string) (interval, error) {
	parts := strings.Split(time, ":")
	if len(parts) != 2 {
		return interval{}, fmt.Errorf("invalid interval format")
	}

	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return interval{}, err
	}

	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return interval{}, err
	}

	return interval{Start: start, End: end}, nil

}
