// This program performs administrative tasks for the garage sale service.
package main

import (
	// Core Packages
	"fmt"
	"log" // https://golang.org/pkg/log/
	"os"

	// Third-party packages
	"github.com/pkg/errors"

	// Internal applcation packages
    "github.com/deezone/HydroBytes-BaseStation/internal/platform/conf"
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/database"
	"github.com/deezone/HydroBytes-BaseStation/internal/schema"
)

// Main entry point for program.
func main() {

	// Only call Exit in main() to allow all defers to complete before shutdown in the case of an error
	if err := run(); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

// Main application logic.
func run() error {

	// =========================================================================
	// Configuration

	var cfg struct {
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:false"`
		}
		Args conf.Args
	}

	if err := conf.Parse(os.Args[1:], "STATIONS", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("STATIONS", &cfg)
			if err != nil {
				log.Fatalf("main : generating usage : %v", err)
			}
			fmt.Println(usage)
			return
		}
		log.Fatalf("error: parsing config: %s", err)
	}

	// =========================================================================
	// Database connection

	// Initialize dependencies.
	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		log.Fatalf("error: connecting to db: %s", err)
	}
	defer db.Close()

	// =========================================================================
	// Supported admin commands

	switch cfg.Args.Num(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Println("error applying migrations", err)
			os.Exit(1)
		}
		fmt.Println("Migrations complete")
		return

	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Println("error seeding database", err)
			os.Exit(1)
		}
		fmt.Println("Seed data complete")
		return
	}
}
