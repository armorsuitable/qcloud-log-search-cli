package main

import (
	"fmt"
	"log-search/internal/qcloud"
	"regexp"
)

func main() {
	client := qcloud.NewQCloudLogSearchClientContext()
	argument := client.CreateCliParameter()

	fmt.Printf("%v\n", argument)

	contents := client.SearchLogs(argument)

	for _, content := range contents {
		highlighted := highlightKeyword(content.Content, argument.Keyword)
		fmt.Println(highlighted)
	}
}

func highlightKeyword(text, keyword string) string {
	// ANSI color codes for red background (grep style)
	const (
		redBG = "\033[41m"
		reset = "\033[0m"
	)

	// Case-insensitive regex pattern
	pattern := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(keyword))

	return pattern.ReplaceAllString(
		text,
		redBG+"$0"+reset,
	)
}
