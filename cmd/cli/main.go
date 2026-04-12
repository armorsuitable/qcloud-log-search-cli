package main

import (
	"fmt"
	"log-search/internal/qcloud"
	"os"
)

func main() {
	topicId := os.Getenv("QCLOUD_TOPIC_ID")
	fmt.Println(topicId)

	client := qcloud.NewQCloudLogSearchClientContext()
	contents := client.SearchLogs(topicId, "last1h", "ERROR")

	for _, content := range contents {
		fmt.Println(content)
	}
}
