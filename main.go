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
	"github.com/sendgrid/sendgrid-go"
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
	senderEmail := getEnv("SENDER_EMAIL", "")
	emailSubject := getEnv("EMAIL_SUBJECT", "")
	recipientEmails := getEnv("RECIPIENT_EMAILS", "")
	hour := getEnv("SCHEDULE_HOUR", "8")      // Default to 8
	minute := getEnv("SCHEDULE_MINUTE", "30") // Default to 30
	timezoneStr := getEnv("TIMEZONE", "America/Toronto") // Default to America/Toronto
}

	var processor DataProcessor = countMembersProcessor

	// If the --now flag is provided, run the task immediately and exit
	if *nowFlag {
		runScheduledTask(redashBaseURL, redashAPIKey, redashQueryID, googleChatWebhookURL, processor)
		return
	}

	// Load the location from the timezone string
	location, err := time.LoadLocation(timezoneStr)
	if err != nil {
	    log.Fatal("Invalid timezone: ", err)
	}

	// Create a new cron scheduler with the loaded location
	c := cron.New(cron.WithSeconds(), cron.WithLocation(location))

	// Constructing the schedule using the retrieved hour and minute
	schedule := fmt.Sprintf("0 %s %s * * *", minute, hour)

	// Scheduling the task to run at the specified time
	_, err = c.AddFunc(schedule, func() {
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
func runScheduledTask(redashBaseURL, redashAPIKey, redashQueryID, googleChatWebhookURL, senderEmail, emailSubject, recipientEmails string, processor DataProcessor) {
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

	// Send the processed data via email
	err = sendEmail(senderEmail, emailSubject, recipientEmails, count)
	if err != nil {
		log.Printf("Error sending email: %v", err)
	}
}

// countMembersProcessor is a function to process the data and return the number of rows as the processed data
func countMembersProcessor(rows []interface{}) (interface{}, error) {
	return len(rows), nil
}

// sendEmail is a function to send the count message via email
func sendEmail(senderEmail, emailSubject, recipientEmails string, count int) error {
	// Format the current date to "Month Day, Year" format
	currentDate := time.Now().Format("January 2, 2006")

	// Construct the message to be sent
	message := map[string]interface{}{
		"from": map[string]string{
			"email": senderEmail,
		},
		"subject": emailSubject,
		"content": []map[string]string{
			{
				"type":  "text/plain",
				"value": fmt.Sprintf("Total member count for %s: %d", currentDate, count),
			},
		},
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{
					{
						"email": recipientEmails,
					},
				},
			},
		},
	}

	// Marshal the message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Send the message using the SendGrid API
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = messageBytes
	response, err := sendgrid.API(request)
	if err != nil {
		return err
	}

	// Check the response status
	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to send email, status: %d, response: %s", response.StatusCode, response.Body)
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
