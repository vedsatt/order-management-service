package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/config"
)

func main() {
	var migrationsPath string
	var command string

	flag.StringVar(&migrationsPath, "path", "./migrations", "Path to migration files")
	flag.StringVar(&command, "command", "up", "Migration command: up or down")
	flag.Parse()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	addr := net.JoinHostPort(cfg.PostgresCfg.Host, cfg.PostgresCfg.Port)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.PostgresCfg.User,
		cfg.PostgresCfg.Password,
		addr,
		cfg.PostgresCfg.DBName,
	)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dsn,
	)

	if err != nil {
		log.Fatalf("failed to create migrations: %v", err)
	}
	defer m.Close()

	switch command {
	case "up":
		if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Printf("migration up failed: %v", err)
			return
		}
		log.Println("migration up completed successfully")

	case "down":
		if err = m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Printf("migration down failed: %v", err)
			return
		}
		log.Println("migration down completed successfully")

	default:
		log.Printf("unknown command: %s", command)
		return
	}
}
