package ratelimiter

import (
	"sync"
	"time"
)

// 固定窗口
type FixedWindowRateLimiter struct {
	sync.RWMutex
	clients map[string]int
	limit   int
	window  time.Duration
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}

// 检查是否允许某个 IP 地址在特定时间窗口内进行访问，并返回一个布尔值和一个时间间隔
func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.RLock()
	count, exists := rl.clients[ip]
	rl.RUnlock()

	// 没访问过，或者访问次数小于限制
	if !exists || count < rl.limit {
		
		rl.Lock()
		if !exists {
			go rl.resetCount(ip)
		}

		rl.clients[ip]++
		rl.Unlock()
		return true, 0
	}

	return false, rl.window
}


// 恢复出厂设置
func (rl *FixedWindowRateLimiter) resetCount(ip string) {
	// 预热
	time.Sleep(rl.window)

	rl.Lock()
	delete(rl.clients, ip)
	rl.Unlock()
}
