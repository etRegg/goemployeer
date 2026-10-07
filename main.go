package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Employee struct {
	ID     int     `json:"ID"`
	Hours  float64 `json:"Hours"`
	Salary float64 `json:"Salary"`
}

func main() {
	// Cada endpoint debe devolver un JSON con una lista de empleados.
	endpoints := []string{
		"https://api.example.com/data1",
		"https://api.example.com/data2",
		"https://api.example.com/data3",
		"https://api.example.com/data4",
		"https://api.example.com/data5",
		"https://api.example.com/data6",
	}

	cacheManager := NewCacheManager()

	// Permite cerrar el programa correctamente con Ctrl+C o SIGTERM.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nGracefully shutting down...")
		os.Exit(0)
	}()

	numWorkers := len(endpoints)
	workers := make(chan struct{}, numWorkers)

	var wg sync.WaitGroup

	for _, endpoint := range endpoints {
		wg.Add(1)
		workers <- struct{}{}

		go func(ep string) {
			defer wg.Done()
			defer func() {
				<-workers
			}()

			if err := fetchData(cacheManager, ep); err != nil {
				fmt.Printf("Endpoint %s failed: %v\n", ep, err)
			}
		}(endpoint)
	}

	wg.Wait()

	employees := cacheManager.GetProcessedEmployees()

	fmt.Printf("Processed employees: %d\n", len(employees))
	for _, employee := range employees {
		fmt.Printf(
			"ID: %d | Hours: %.2f | Salary: %.2f\n",
			employee.ID,
			employee.Hours,
			employee.Salary,
		)
	}
}
