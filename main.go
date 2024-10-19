package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

// Initialize Redis client
func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis-server:6379", // Use the Docker Compose Redis service name
	})
}

// Simulate an HTTP request
func performHTTPRequest() {

	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		fmt.Println("HTTP request failed:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("HTTP request completed with status:", resp.StatusCode)
}

// Simulate a Redis operation
func performRedisOperation() {
	ctx := context.Background()
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		fmt.Println("Redis error:", err)
		return
	}
	fmt.Println("Redis operation completed")
}

// Monitor memory usage
func monitorResourceUsage() {
	var memStats runtime.MemStats
	var cpuUsage float64
	numCPU := runtime.NumCPU()
	prevCPUTime := getCPUTime()

	for {
		runtime.ReadMemStats(&memStats)
		currCPUTime := getCPUTime()
		cpuUsage = calculateCPUUsage(prevCPUTime, currCPUTime, numCPU)
		prevCPUTime = currCPUTime

		fmt.Printf("Allocated memory: %v MB, CPU Usage: %.2f%%\n", memStats.Alloc/1024/1024, cpuUsage)
		time.Sleep(2 * time.Second)
	}
}

func getCPUTime() int64 {
	return time.Now().UnixNano()

}

func calculateCPUUsage(prevTime, currTime int64, numCPU int) float64 {
	cpuTimeDiff := float64(currTime - prevTime)
	totalTime := float64(2 * time.Second.Nanoseconds() * int64(numCPU))
	return (cpuTimeDiff / totalTime) * 100
}

// Benchmark handler to trigger tasks
func benchmarkHandler(c *gin.Context) {
	// Number of tasks to run
	numTasks := 10
	if value := c.Query("tasks"); value != "" {
		fmt.Sscanf(value, "%d", &numTasks)
	}

	for i := 0; i < numTasks; i++ {

		performHTTPRequest()    // Simulate HTTP request
		performRedisOperation() // Simulate Redis operation
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Ran %d HTTP and Redis tasks", numTasks)})
}

func main() {

	// Initialize Redis
	initRedis()

	// Start monitoring memory usage in a Go routine
	go monitorResourceUsage()

	// Create Gin router
	r := gin.Default()

	// Define the benchmark route
	r.GET("/benchmark", benchmarkHandler)

	// Start the Gin server
	r.Run(":8080")
}
