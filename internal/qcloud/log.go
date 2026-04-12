package qcloud

import (
	"fmt"
	"log"
	"os"
	"time"

	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"

	_ "github.com/joho/godotenv/autoload"
)

type QCloudLogSearchClientContext struct {
	ApiClient *cls.Client
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
		endpoint = "internal.tencentcloudapi.com"
	}

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = fmt.Sprintf("cls.%s", endpoint)
	client, initErr := cls.NewClient(cred, region, cpf)
	if initErr != nil {
		log.Fatal(initErr)
	}

	return &QCloudLogSearchClientContext{
		ApiClient: client,
	}
}

func (c *QCloudLogSearchClientContext) SearchLogs(topicId, periodFormat, query string) []string {
	if len(periodFormat) == 0 {
		periodFormat = "last15m"
	}

	utc8, _ := time.LoadLocation("Asia/Shanghai")
	startTime := time.Now()
	endStr := time.Now().Format(time.DateTime)
	endTimeStamp, _ := time.ParseInLocation(time.DateTime, endStr, utc8)

	switch periodFormat {
	case "last15m":
		startTime = startTime.Add(-15 * time.Minute)
	case "last1h":
		startTime = startTime.Add(-1 * time.Hour)
	case "last6h":
		startTime = startTime.Add(-6 * time.Hour)
	case "last1day":
		startTime = startTime.Add(-24 * time.Hour)
	case "last7d":
		startTime = startTime.Add(-7 * 24 * time.Hour)
	default:
		log.Printf("Unsupported period format: %s. Defaulting to last15m.", periodFormat)
		startTime = startTime.Add(-15 * time.Minute)
	}

	startStr := startTime.Format(time.DateTime)
	startTimeStamp, _ := time.ParseInLocation(time.DateTime, startStr, utc8)

	f := startTimeStamp.UnixMilli()
	t := endTimeStamp.UnixMilli()
	sort := "asc"
	limit := int64(500)

	logReq := cls.NewSearchLogRequest()
	logReq.TopicId = common.StringPtr(topicId)
	logReq.From = &f
	logReq.To = &t
	logReq.Query = common.StringPtr(query)
	logReq.Sort = common.StringPtr(sort)
	logReq.Limit = &limit

	resp, searErr := c.ApiClient.SearchLog(logReq)
	if searErr != nil {
		log.Fatal(searErr)
	}

	log.Printf("logSearch response:%s", resp.ToJsonString())

	if len(resp.Response.Results) == 0 {
		log.Printf("not result log found.")
		return []string{}
	}

	logContent := make([]string, 0)

	for _, result := range resp.Response.Results {
		logContent = append(logContent, *result.LogJson)
	}

	return logContent
}
