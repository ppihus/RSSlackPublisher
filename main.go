package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v2"
)

// SlackRequestBody structure for sending messages to Slack
type SlackRequestBody struct {
	Text string `json:"text"`
}

// Config structure for application configuration
type Config struct {
	RSSFeeds     []string `yaml:"rss_feeds"`
	SlackWebhook string   `yaml:"slack_webhook"`
	MaxNewsCount int      `yaml:"max_news_count"`
}

// SendSlackNotification sends a message to a Slack channel
func SendSlackNotification(webhookURL string, msg string) error {
	slackBody, err := json.Marshal(SlackRequestBody{Text: msg})
	if err != nil {
		return fmt.Errorf("error marshaling Slack message: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return fmt.Errorf("error creating Slack request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending Slack request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("non-OK response from Slack: %s, body: %s", resp.Status, string(body))
	}

	return nil
}

// loadConfig loads configuration from a YAML file
func loadConfig(filePath string) (Config, error) {
	var config Config
	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshaling config data: %w", err)
	}

	return config, nil
}

// initializeSentNewsFile ensures the existence of the sent news file
func initializeSentNewsFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("error creating sent news file: %w", err)
		}
	}
	return nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// main function of the program
func main() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if slackWebhook := os.Getenv("SLACK_WEBHOOK"); slackWebhook != "" {
		config.SlackWebhook = slackWebhook
	}

	if config.SlackWebhook == "" {
		log.Fatal("Slack webhook URL is not configured. Exiting.")
	}

	err = initializeSentNewsFile("sent_news.txt")
	if err != nil {
		log.Fatalf("Failed to initialize sent news file: %v", err)
	}

	data, err := os.ReadFile("sent_news.txt")
	if err != nil {
		log.Fatalf("Failed to read sent news file: %v", err)
	}
	sentNews := strings.Split(string(data), "\n")

	var totalSentCount int
	var unsentItemCount int
	fp := gofeed.NewParser()

	for _, feedURL := range config.RSSFeeds {
		feed, err := fp.ParseURL(feedURL)
		if err != nil {
			log.Printf("Error parsing RSS feed: %v", err)
			continue
		}

		for _, item := range feed.Items {
			if !contains(sentNews, item.Link) {
				unsentItemCount++ // Count only unsent news items
				if totalSentCount >= config.MaxNewsCount {
					// If limit is reached, don't process further items
					continue
				}

				msg := fmt.Sprintf("Title: %s\nLink: %s\nPublished: %s", item.Title, item.Link, item.Published)
				err := SendSlackNotification(config.SlackWebhook, msg)
				if err != nil {
					log.Printf("Error sending Slack notification: %v", err)
					continue
				}
				totalSentCount++

				f, err := os.OpenFile("sent_news.txt", os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					log.Printf("Error opening sent news file: %v", err)
					continue
				}
				_, err = f.WriteString(item.Link + "\n")
				f.Close()
				if err != nil {
					log.Printf("Error writing to sent news file: %v", err)
				}
			}
		}
	}

	if unsentItemCount > totalSentCount && totalSentCount > 0 {
		remainingNews := unsentItemCount - totalSentCount
		msg := fmt.Sprintf("Limit reached. Sent %d out of %d unsent news items. %d news items will be considered in the next run.", totalSentCount, unsentItemCount, remainingNews)
		SendSlackNotification(config.SlackWebhook, msg)
	}
}
