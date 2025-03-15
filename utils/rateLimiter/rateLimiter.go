package rateLimiter

import (
	"sync"
	"time"
	"golang.org/x/time/rate"
	"DBHS/config"
)

// Client holds the rate limiter and last activity time for a client
type Client struct {
	Limiter       *rate.Limiter
	lastSeen      time.Time
}

// IPRateLimiter stores clients with their rate limiters and last activity
type IPRateLimiter struct {
	clients map[string]*Client
	mu      *sync.RWMutex
	r       rate.Limit
	b       int
}

// NewIPRateLimiter creates a new rate limiter that allows r requests per second with a burst of b
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		clients: make(map[string]*Client),
		mu:      &sync.RWMutex{},
		r:       r,
		b:       b,
	}
}


// CleanupStaleClients removes clients that haven't been seen for the specified duration
func (i *IPRateLimiter) CleanupStaleClients(staleTimeout time.Duration) int {
	i.mu.Lock()
	defer i.mu.Unlock()

	now := time.Now()
	count := 0

	for ip, client := range i.clients {
		// remove the client from the clients map when it has not been seen for longer than the stale timeout
		if now.Sub(client.lastSeen) > staleTimeout {
			delete(i.clients, ip)
			count++
		}
	}

	return count
}

func StartCleanupStaleClients(i *IPRateLimiter)  {
    if len(i.clients) > 1000 {
		go func() {
		    for {
				cleaned := Newlimiter.CleanupStaleClients(time.Minute * 30)
				config.App.InfoLog.Printf("Cleaned up %d stale clients\n", cleaned)
			}
		}()
	}
}

// AddClient creates a new client with a rate limiter and adds it to the clients map, using the IP address as the key
func (i *IPRateLimiter) AddClient(ip string) *Client {
    StartCleanupStaleClients(i); // Start the cleanup goroutine

	i.mu.Lock()
	defer i.mu.Unlock()

	client := &Client{
		Limiter:  rate.NewLimiter(i.r, i.b),
		lastSeen: time.Now(),
	}
	i.clients[ip] = client

	return client
}

// GetClient returns the client for the provided IP address if it exists
// Otherwise calls AddClient to add a new client for the IP
func (i *IPRateLimiter) GetClient(ip string) *Client {
    StartCleanupStaleClients(i); // Start the cleanup goroutine

	i.mu.RLock()
	client, exists := i.clients[ip]
	i.mu.RUnlock()

	if !exists {
		return i.AddClient(ip)
	}

	// Update the last seen time
	i.mu.Lock()
	client.lastSeen = time.Now()
	i.mu.Unlock()

	return client
}

// Create a global rate limiter that allows 15 requests per second with a burst of 25
var Newlimiter = NewIPRateLimiter(15, 25)
