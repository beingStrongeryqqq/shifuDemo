package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	pollingInterval = 10 * time.Second
	measurements    []float64
	mu              sync.Mutex
)

func fetchMeasurement() {
	url := "http://deviceshifu-plate-reader.deviceshifu.svc.cluster.local/get_measurement"

	for {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error making request: %v", err)
			time.Sleep(pollingInterval)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Received non-200 response: %d", resp.StatusCode)
			time.Sleep(pollingInterval)
			continue
		}

		var measurement float64 // 根据实际返回值类型调整
		if err := json.NewDecoder(resp.Body).Decode(&measurement); err != nil {
			log.Printf("Error decoding response: %v", err)
			time.Sleep(pollingInterval)
			continue
		}

		mu.Lock()
		measurements = append(measurements, measurement)
		mu.Unlock()

		avg := calculateAverage()
		log.Printf("Average Measurement: %.2f", avg)

		time.Sleep(pollingInterval)
	}
}

func calculateAverage() float64 {
	mu.Lock()
	defer mu.Unlock()

	if len(measurements) == 0 {
		return 0
	}

	sum := 0.0
	for _, m := range measurements {
		sum += m
	}
	return sum / float64(len(measurements))
}

func main() {
	fetchMeasurement()
}
