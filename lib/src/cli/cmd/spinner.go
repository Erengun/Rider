package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"time"
	"net/http"
)

type Spinner struct {
	Interval int      `json:"interval"`
	Frames   []string `json:"frames"`
}

// Global variable to store all spinners
var allSpinners map[string]Spinner

func init() {
	// Seed the random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))

	resp, err := http.Get("https://raw.githubusercontent.com/sindresorhus/cli-spinners/master/spinners.json")
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	n := []byte(body)

	

	// Unmarshal JSON data into the allSpinners map
	err = json.Unmarshal(n, &allSpinners)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
}

// Function to get a random spinner
func getRandomSpinner() Spinner {
	keys := make([]string, 0, len(allSpinners))
	for key := range allSpinners {
		keys = append(keys, key)
	}

	randomKey := keys[rand.Intn(len(keys))]
	return allSpinners[randomKey]
}

func LoadingSpinner(stopChan chan struct{}, msg string) {
	// Get a random
	spinner := getRandomSpinner()

	i := 0
	for {
		select {
		case <-stopChan:
			return
		default:
			fmt.Printf("\r%s \033[32m%s\033[0m", msg, spinner.Frames[i])
			i = (i + 1) % len(spinner.Frames)
			time.Sleep(time.Duration(spinner.Interval) * time.Millisecond)
		}
	}
}
