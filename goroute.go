package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

func fetchData(startTime, endTime time.Time, wg *sync.WaitGroup, ch chan struct{}) {
	defer wg.Done()

	fmt.Printf("Fetching data from %s to %s...\n", startTime.Format("15:04"), endTime.Format("15:04"))
	// Simulating data fetching process
	time.Sleep(2 * time.Second)
	fmt.Printf("Data fetched from %s to %s.\n", startTime.Format("15:04"), endTime.Format("15:04"))

	// Signal that fetching is done for this interval
	ch <- struct{}{}
}

func compareData(startTime, endTime time.Time, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Comparing data from %s to %s...\n", startTime.Format("15:04"), endTime.Format("15:04"))
	// Simulating data comparison process
	time.Sleep(2 * time.Second)
	fmt.Printf("Data compared from %s to %s.\n", startTime.Format("15:04"), endTime.Format("15:04"))
}

func main() {

	os.Setenv("START_TIME", "08:00")
	os.Setenv("END_TIME", "18:00")
	os.Setenv("INTERVAL", "2")

	var wg sync.WaitGroup

	// Get start time, end time, and interval from environment variables
	startTimeStr := os.Getenv("START_TIME")
	endTimeStr := os.Getenv("END_TIME")
	intervalStr := os.Getenv("INTERVAL")

	startTime, err := time.Parse("15:04", startTimeStr)
	if err != nil {
		fmt.Println("Error parsing start time:", err)
		return
	}

	endTime, err := time.Parse("15:04", endTimeStr)
	if err != nil {
		fmt.Println("Error parsing end time:", err)
		return
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		fmt.Println("Error parsing interval:", err)
		return
	}

	ch := make(chan struct{}, 1) // Channel to signal completion of fetching

	// Schedule data fetching and comparison routines
	for t := startTime; t.Before(endTime); t = t.Add(time.Duration(interval) * time.Hour) {
		intervalEndTime := t.Add(time.Duration(interval) * time.Hour)

		// Fetch data in a goroutine
		wg.Add(1)
		go fetchData(t, intervalEndTime, &wg, ch)

		// Wait for fetching to complete before starting comparison
		<-ch

		// Compare data after fetching
		wg.Add(1)
		go compareData(t, intervalEndTime, &wg)
	}

	wg.Wait()
	fmt.Println("All data fetching and comparison completed.")
}
