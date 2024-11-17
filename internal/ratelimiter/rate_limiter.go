package ratelimiter

import "time"


type Limiter interface {
	Allow(ip string) (bool, time.Duration)
}

type Config struct {
	RequestsPerTimeFrame int              // 每个时间段的请求量  
	TimeFrame            time.Duration
	Enabled              bool
}
