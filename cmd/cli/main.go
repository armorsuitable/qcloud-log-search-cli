package main

import (
	"fmt"
	"log"
	"log-search/internal/qcloud"
	"os"

	"github.com/manifoldco/promptui"
)

func main() {
	validate := func(keywordQuery string) error {
		if len(keywordQuery) == 0 {
			return fmt.Errorf("keyword query cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Query keyword",
		Validate: validate,
	}

	keywordQuery, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	periodPrompt := promptui.Select{
		Label: "Select period format",
		Items: []string{"last15m", "last1h", "last6h", "last1d", "last7d"},
	}

	_, periodResult, err := periodPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("period: %s keyword: %s\n", periodResult, keywordQuery)

	topicId := os.Getenv("QCLOUD_TOPIC_ID")
	// fmt.Println(topicId)
	if len(topicId) == 0 {
		log.Fatal("QCLOUD_TOPIC_ID environment variable is not set")
	}

	client := qcloud.NewQCloudLogSearchClientContext()
	contents := client.SearchLogs(topicId, periodResult, keywordQuery)

	for _, content := range contents {
		fmt.Println(content.Content)
	}
}
