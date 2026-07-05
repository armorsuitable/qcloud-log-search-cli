package qcloud

import (
	"flag"
	"fmt"
	"strings"
)

var QueryPeriod = []string{
	"last15m",
	"last1h",
	"last6h",
	"last1d",
	"last7d",
	"last30d",
}

func NonInteractiveCommandLineQuery(topicId string, logLimit int64) QCloudLogQuery {
	query := flag.String("query", "", "Search keyword for query")
	period := flag.String("period", "last15m", "Time period for log search (e.g., last15m, last1h, last6h, last1d, last7d, last30d)")
	paramTopicId := flag.String("topicId", topicId, "QCloud CLS Topic ID")
	sortType := flag.String("sort", "asc", "Sort type for log search (asc or desc)")

	queryContext := flag.Bool("queryContext", false, "Query context for log search (optional)")
	pkgId := flag.String("pkgId", "", "Package ID for log search (optional)")
	pkgLogId := flag.Int64("pkgLogId", 0, "Package Log ID for log search (optional)")
	pkgLogTime := flag.String("pkgLogTime", "", "Package Log Time for log search (optional)")

	flag.Parse()

	if *queryContext {
		return QCloudLogQuery{
			Period:          *period,
			TopicId:         *paramTopicId,
			SortType:        *sortType,
			LogLimit:        logLimit,
			QueryLogContext: true,
			PackageId:       *pkgId,
			PackageLogId:    *pkgLogId,
			PackageLogTime:  *pkgLogTime,
		}
	}

	return QCloudLogQuery{
		Keyword: convertAndQuery(*query),
		Period:  *period,
		TopicId: *paramTopicId,

		SortType: *sortType,
		LogLimit: logLimit,
	}
}

func convertAndQuery(query string) string {
	if !strings.Contains(query, "AND") {
		return query
	}

	parts := strings.Split(query, "AND")
	convertedParts := make([]string, 0)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len(part) > 0 {
			convertedParts = append(convertedParts, fmt.Sprintf(`\"%s\"`, part))
		}
	}

	return strings.Join(convertedParts, " AND ")
}
