package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/garrettladley/peli/internal/service/restaurant"
	"github.com/garrettladley/peli/internal/sqlite"
	"github.com/garrettladley/peli/internal/tui"
)

var seedData = []restaurant.SeedRestaurant{
	{Name: "Sushi Nakazawa", Cuisine: "Japanese", Position: 0},
	{Name: "Taco Bell Cantina", Cuisine: "Mexican", Position: 1},
	{Name: "Joe's Pizza", Cuisine: "Italian", Position: 2},
	{Name: "Ichiran Ramen", Cuisine: "Japanese", Position: 3},
	{Name: "Shake Shack", Cuisine: "American", Position: 4},
	{Name: "Pad Thai Palace", Cuisine: "Thai", Position: 5},
	{Name: "Curry House", Cuisine: "Indian", Position: 6},
	{Name: "Waffle House", Cuisine: "American", Position: 7},
}

func run() error {
	seed := flag.Bool("seed", false, "reset database with sample data (deletes all existing restaurants)")
	flag.Parse()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home directory: %w", err)
	}

	dataDir := filepath.Join(homeDir, ".peli")
	if err := os.MkdirAll(dataDir, 0o750); err != nil {
		return fmt.Errorf("create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "peli.db")

	ctx := context.Background()
	database, err := sqlite.New(ctx, dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "close database: %v\n", err)
		}
	}()

	svc := restaurant.New(database)

	if *seed {
		if err := svc.Seed(ctx, seedData); err != nil {
			return fmt.Errorf("seed database: %w", err)
		}
		fmt.Println("Database seeded successfully")
		return nil
	}

	if err := tui.Run(svc); err != nil {
		return fmt.Errorf("run tui: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
