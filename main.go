package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

// QueryResult struct represents the JSON structure of the API response
type QueryResult struct {
	QueryResult struct {
		Data struct {
			Rows []interface{} `json:"rows"`
		} `json:"data"`
	} `json:"query_result"`
}

// DataProcessor is a type that represents a function to process the data from the API response
type DataProcessor func([]interface{}) (interface{}, error)

// init function to load environment variables from a .env file
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	nowFlag := flag.Bool("now", false, "Run the task immediately, ignoring the schedule")
	flag.Parse()

	// Load the configurations from environment variables
	redashBaseURL := getEnv("REDASH_BASE_URL", "")
	redashAPIKey := getEnv("REDASH_API_KEY", "")
	redashQueryID := getEnv("REDASH_QUERY_ID", "")
	googleChatWebhookURL := getEnv("GOOGLE_CHAT_WEBHOOK_URL", "")
	hour := getEnv("SCHEDULE_HOUR", "8")      // Default to 8
	minute := getEnv("SCHEDULE_MINUTE", "30") // Default to 30
	timezone := getEnv("TIMEZONE", "America/Toronto") // Default to America/Toronto

	var processor DataProcessor = countMembersProcessor

	// If the --now flag is provided, run the task immediately and exit
	if *nowFlag {
		runScheduledTask(redashBaseURL, redashAPIKey, redashQueryID, googleChatWebhookURL, processor)
		return
	}

	// Create a new cron scheduler
	c := cron.New(cron.WithLocation(timezone), cron.WithSeconds())

	// Constructing the schedule using the retrieved hour and minute
	schedule := fmt.Sprintf("0 %s %s * * *", minute, hour)

	// Scheduling the task to run at the specified time
	_, err := c.AddFunc(schedule, func() {
		runScheduledTask(redashBaseURL, redashAPIKey, redashQueryID, googleChatWebhookURL, processor)
	})
	if err != nil {
		log.Fatal("Could not schedule task: ", err)
	}
	c.Start()

	// Use a WaitGroup to keep the application running indefinitely
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

// runScheduledTask is a function to run the scheduled task for fetching and processing the data
func runScheduledTask(redashBaseURL, redashAPIKey, redashQueryID, googleChatWebhookURL string, processor DataProcessor) {
	// Construct the URLs for refreshing and fetching results
	refreshURL := fmt.Sprintf("%s/api/queries/%s/refresh", redashBaseURL, redashQueryID)
	resultsURL := fmt.Sprintf("%s/api/queries/%s/results.json?api_key=%s", redashBaseURL, redashQueryID, redashAPIKey)

	// Create a custom HTTP client
	client := &http.Client{}

	// Refresh the query with Authorization header
	req, err := http.NewRequest("POST", refreshURL, nil)
	if err != nil {
		log.Println("Error creating refresh request: ", err)
		return
	}
	req.Header.Add("Authorization", "Key "+redashAPIKey)
	_, err = client.Do(req)
	if err != nil {
		log.Println("Error refreshing query: ", err)
		return
	}

	// Sleep or poll until the results are ready
	time.Sleep(10 * time.Second) // Adjust as needed

	// Fetch the results using API key in the URL
	resp, err := http.Get(resultsURL)
	if err != nil {
		log.Println("Error getting results: ", err)
		return
	}
	defer resp.Body.Close()

	// Decode the results into the QueryResult struct
	var result QueryResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("Error decoding results: ", err)
		return
	}

	// Process the data using the provided processor function
	processedData, err := processor(result.QueryResult.Data.Rows)
	if err != nil {
		log.Printf("Error processing data: %v", err)
		return
	}

	// Convert the processed data to an integer
	count, ok := processedData.(int)
	if !ok {
		log.Printf("Error: Processed data is not an integer")
		return
	}

	// Send the processed data to Google Chat
	err = sendMessageToGoogleChat(googleChatWebhookURL, count)
	if err != nil {
		log.Printf("Error sending message to Google Chat: %v", err)
	}
}

// countMembersProcessor is a function to process the data and return the number of rows as the processed data
func countMembersProcessor(rows []interface{}) (interface{}, error) {
	return len(rows), nil
}

// sendMessageToGoogleChat is a function to send the count message to Google Chat Webhook
func sendMessageToGoogleChat(webhookURL string, count int) error {
	// Format the current date to "Month Day, Year" format
	currentDate := time.Now().Format("January 2, 2006")

	// Construct the message to be sent
	message := map[string]interface{}{
		"text": fmt.Sprintf("Total member count for %s: *%d*", currentDate, count),
	}

	// Marshal the message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Send the message to the provided webhook URL
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(messageBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send message to Google Chat, status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

// getEnv is a function to fetch the value of the environment variable identified by key,
// returns fallback if the environment variable is not set.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}