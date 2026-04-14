package qcloud

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"

	_ "github.com/joho/godotenv/autoload"
)

type QCloudLogQuery struct {
	Keyword string
	Period  string
	TopicId string

	SortType string
	LogLimit int64
}

type QCloudLogJsonFormat struct {
	Content string `json:"__CONTENT__"`
	Tag     any    `json:"__TAG__"`
}

type QCloudLogSearchClientContext struct {
	ApiClient           *cls.Client
	InteractiveArgModel bool
}

func NewQCloudLogSearchClientContext() *QCloudLogSearchClientContext {
	cred := common.NewCredential(
		os.Getenv("QCLOUD_SECRET_ID"),
		os.Getenv("QCLOUD_SECRET_KEY"),
	)

	region := os.Getenv("QCLOUD_REGION")
	if len(region) == 0 {
		region = regions.Beijing
	}

	endpoint := os.Getenv("QCLOUD_ENDPOINT")
	if len(endpoint) == 0 {
		endpoint = "tencentcloudapi.com"
	}

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = fmt.Sprintf("cls.%s", endpoint)
	client, initErr := cls.NewClient(cred, region, cpf)
	if initErr != nil {
		log.Fatal(initErr)
	}

	interactiveArgModel := os.Getenv("INTERACTIVE_PARAM_MODE") == "true"

	return &QCloudLogSearchClientContext{
		ApiClient:           client,
		InteractiveArgModel: interactiveArgModel,
	}
}

func (c *QCloudLogSearchClientContext) SearchLogs(query QCloudLogQuery) []QCloudLogJsonFormat {
	if len(query.Period) == 0 {
		query.Period = "last15m"
	}

	utc8, _ := time.LoadLocation("Asia/Shanghai")
	startTime := time.Now()
	endStr := time.Now().Format(time.DateTime)
	endTimeStamp, _ := time.ParseInLocation(time.DateTime, endStr, utc8)

	switch query.Period {
	case "last15m":
		startTime = startTime.Add(-15 * time.Minute)
	case "last1h":
		startTime = startTime.Add(-1 * time.Hour)
	case "last6h":
		startTime = startTime.Add(-6 * time.Hour)
	case "last1d":
		startTime = startTime.Add(-24 * time.Hour)
	case "last7d":
		startTime = startTime.Add(-7 * 24 * time.Hour)
	default:
		log.Printf("Unsupported period format: %s. Defaulting to last15m.", query.Period)
		startTime = startTime.Add(-15 * time.Minute)
	}

	startStr := startTime.Format(time.DateTime)
	startTimeStamp, _ := time.ParseInLocation(time.DateTime, startStr, utc8)

	f := startTimeStamp.UnixMilli()
	t := endTimeStamp.UnixMilli()
	// sort := "asc"
	// limit := int64(500)

	logReq := cls.NewSearchLogRequest()
	logReq.TopicId = common.StringPtr(query.TopicId)
	logReq.From = &f
	logReq.To = &t
	logReq.Query = common.StringPtr(query.Keyword)
	logReq.Sort = common.StringPtr(query.SortType)
	logReq.Limit = &query.LogLimit

	resp, searErr := c.ApiClient.SearchLog(logReq)
	if searErr != nil {
		log.Fatal(searErr)
	}

	// log.Printf("logSearch response:%s", resp.ToJsonString())

	if len(resp.Response.Results) == 0 {
		log.Printf("not result log found.")
		return []QCloudLogJsonFormat{}
	}

	logContent := make([]QCloudLogJsonFormat, 0)

	for _, result := range resp.Response.Results {
		// logContent = append(logContent, *result.LogJson)
		var logFormat QCloudLogJsonFormat
		unmarshalErr := json.Unmarshal([]byte(*result.LogJson), &logFormat)
		if unmarshalErr != nil {
			log.Printf("Failed to unmarshal log JSON: %v", unmarshalErr)
			continue
		}
		logContent = append(logContent, logFormat)
	}

	return logContent
}

func (c *QCloudLogSearchClientContext) CreateCliParameter() QCloudLogQuery {
	topicId := os.Getenv("QCLOUD_TOPIC_ID")
	if len(topicId) == 0 {
		log.Fatal("QCLOUD_TOPIC_ID environment variable is not set")
	}

	logLimitStr := os.Getenv("QUERY_LOG_LIMIT")
	logLimit := int64(500) // default log limit
	if len(logLimitStr) != 0 {
		logLimit, _ = strconv.ParseInt(logLimitStr, 10, 32)
	}

	if c.InteractiveArgModel {
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
			log.Fatalf("Prompt failed %v\n", err)
		}

		periodPrompt := promptui.Select{
			Label: "Select period format",
			Items: []string{
				"last15m",
				"last1h",
				"last6h",
				"last1d",
				"last7d",
			},
		}

		_, periodResult, err := periodPrompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}

		sortPrompt := promptui.Prompt{
			Label: "Sort type (asc/desc), default asc",
			Validate: func(s string) error {
				if s != "asc" && s != "desc" && len(s) != 0 {
					return fmt.Errorf("invalid sort type, must be 'asc', 'desc', or empty for default")
				}
				return nil
			},
		}

		sortResult, err := sortPrompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}

		return QCloudLogQuery{
			Keyword:  keywordQuery,
			Period:   periodResult,
			TopicId:  topicId,
			SortType: sortResult,
			LogLimit: logLimit,
		}
	}

	query := flag.String("query", "", "Search keyword for query")
	period := flag.String("period", "last15m", "Time period for log search (e.g., last15m, last1h, last6h, last1d, last7d)")
	paramTopicId := flag.String("topicId", topicId, "QCloud CLS Topic ID")
	sortType := flag.String("sort", "asc", "Sort type for log search (asc or desc)")
	flag.Parse()

	return QCloudLogQuery{
		Keyword: *query,
		Period:  *period,
		TopicId: *paramTopicId,

		SortType: *sortType,
		LogLimit: logLimit,
	}
}
