package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	Mp   map[string]*Node
	Lock *sync.Mutex
}

type Node struct {
	LastRefillTime int64
	Tokens         float64
}

const (
	MaxTokens  = 5
	RefillRate = 5.0 / 60.0 // 5 tokens 60 seconds -->  this is rate
)

type ResponseStats struct {
	Tokens     float64
	WaitForSec float64
}

func (m *Middleware) UserStatus(ctx *gin.Context) {
	ip := ctx.ClientIP()

	m.Lock.Lock()
	defer m.Lock.Unlock()
	if _, ok := m.Mp[ip]; !ok {
		ctx.JSON(200, "user dont exist")
		return
	}
	node := m.Mp[ip]
	now := time.Now().UnixNano()
	elapsed := float64(now-node.LastRefillTime) / 1e9 // seconds
	refilled := elapsed * RefillRate

	node.Tokens = minFloat(MaxTokens, node.Tokens+refilled)
	node.LastRefillTime = now

	m.Mp[ip] = node
	ctx.JSON(200, node)
}

func (m *Middleware) RateLimitCheck(ctx *gin.Context) {
	ip := ctx.ClientIP()
	fmt.Println(ip)
	now := time.Now().UnixNano()

	m.Lock.Lock()
	defer m.Lock.Unlock()

	node, exists := m.Mp[ip]

	if !exists {
		// First request  --> already consumed one token
		m.Mp[ip] = &Node{
			LastRefillTime: now,
			Tokens:         MaxTokens - 1,
		}
		return
	}

	// Step 1: Refill tokens
	elapsed := float64(now-node.LastRefillTime) / 1e9 // seconds
	refilled := elapsed * RefillRate

	node.Tokens = minFloat(MaxTokens, node.Tokens+refilled)

	// Step 2: Check tokens
	if node.Tokens < 1 {
		tokensNeeded := 1 - node.Tokens
		waitSeconds := tokensNeeded / RefillRate

		res := &ResponseStats{
			Tokens:     node.Tokens,
			WaitForSec: waitSeconds,
		}
		ctx.AbortWithStatusJSON(429, res)
		return
	}

	// Step 3: Consume token
	node.Tokens -= 1
	node.LastRefillTime = now

	m.Mp[ip] = node
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
