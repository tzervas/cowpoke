package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tzervas/cowpoke/internal/features"
	"github.com/tzervas/cowpoke/internal/similarity"
)

var (
	version = "dev"
)

// Config holds the agent configuration
type Config struct {
	WindowSize          int
	MetricsInterval     int
	SimilarityThreshold float64
	Verbose             bool
}

// Agent represents the Cowpoke agent
type Agent struct {
	config           Config
	featureExtractor *features.Extractor
	similarityEngine *similarity.Engine
	ctx              context.Context
	cancel           context.CancelFunc
}

func main() {
	// Parse flags
	config := Config{}
	flag.IntVar(&config.WindowSize, "window-size", 300, "Time window for feature extraction in seconds")
	flag.IntVar(&config.MetricsInterval, "metrics-interval", 15, "Metrics collection interval in seconds")
	flag.Float64Var(&config.SimilarityThreshold, "similarity-threshold", 0.85, "Cosine similarity threshold")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose logging")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("cowpoke version %s\n", version)
		os.Exit(0)
	}

	// Initialize agent
	agent, err := NewAgent(config)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start agent
	log.Printf("Starting Cowpoke agent (version %s)", version)
	log.Printf("Configuration: window=%ds, interval=%ds, threshold=%.2f",
		config.WindowSize, config.MetricsInterval, config.SimilarityThreshold)

	go agent.Run()

	// Wait for termination signal
	sig := <-sigCh
	log.Printf("Received signal %s, shutting down gracefully...", sig)
	agent.Shutdown()
}

// NewAgent creates a new Cowpoke agent
func NewAgent(config Config) (*Agent, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize components
	featureExtractor := features.NewExtractor(config.WindowSize)
	similarityEngine := similarity.NewEngine(config.SimilarityThreshold)

	agent := &Agent{
		config:           config,
		featureExtractor: featureExtractor,
		similarityEngine: similarityEngine,
		ctx:              ctx,
		cancel:           cancel,
	}

	return agent, nil
}

// Run starts the main agent loop
func (a *Agent) Run() {
	ticker := time.NewTicker(time.Duration(a.config.MetricsInterval) * time.Second)
	defer ticker.Stop()

	log.Println("Agent loop started")

	for {
		select {
		case <-a.ctx.Done():
			log.Println("Agent loop stopped")
			return
		case <-ticker.C:
			a.collectAndProcess()
		}
	}
}

// collectAndProcess collects metrics and processes them
func (a *Agent) collectAndProcess() {
	if a.config.Verbose {
		log.Println("Collecting metrics...")
	}

	// TODO: Collect actual CPU metrics from Kubernetes API
	// For now, this is a placeholder that demonstrates the flow

	// Placeholder: simulate CPU metrics
	cpuMetrics := []float64{50.0, 52.0, 48.0, 55.0, 51.0}

	// Extract features
	featureVector, err := a.featureExtractor.Extract(cpuMetrics)
	if err != nil {
		log.Printf("Failed to extract features: %v", err)
		return
	}

	if a.config.Verbose {
		log.Printf("Extracted features: %v", featureVector)
	}

	// Compute similarity (placeholder - would compare against historical patterns)
	// TODO: Implement pattern storage and retrieval
	similarity := a.similarityEngine.ComputeSimilarity(featureVector, featureVector)

	if a.config.Verbose {
		log.Printf("Similarity score: %.4f", similarity)
	}

	// Make scaling decision based on similarity
	if similarity >= a.config.SimilarityThreshold {
		if a.config.Verbose {
			log.Println("Pattern matches historical workload - maintaining current scale")
		}
	} else {
		if a.config.Verbose {
			log.Println("New workload pattern detected - analyzing scaling recommendations")
		}
	}
}

// Shutdown gracefully stops the agent
func (a *Agent) Shutdown() {
	log.Println("Shutting down agent...")
	a.cancel()
	// Give time for cleanup
	time.Sleep(1 * time.Second)
	log.Println("Agent shutdown complete")
}
