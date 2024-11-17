package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`  // 偏移量, 从第几个开始
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`   // 标签列表, 最多5个
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`  // 起始时间
	Until  string   `json:"until"`  // 截至时间
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}

		fq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		l, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}

		fq.Offset = l
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort  // DESC 降序, ASC 升序
	}

	tags := qs.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	} else {
		fq.Tags = []string{}
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		fq.Since = parseTime(since)
	}

	until := qs.Get("until")
	if until != "" {
		fq.Until = parseTime(until)
	}

	return fq, nil
}

func parseTime(s string) string {
	// layout, time.DateTime == 2006-01-02 15:04:05
	
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}
