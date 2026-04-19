package cmd

import (
	secondarypostgres "codebase-app/internal/framework/secondary/db/postgres"
	"codebase-app/internal/infrastructure/config"
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type dbBenchResult struct {
	Name           string
	Host           string
	ConnectTime    time.Duration
	QueryDurations []time.Duration
}

func RunDBBench(cmd *flag.FlagSet, args []string) {
	query := cmd.String("query", "SELECT 1", "SQL query to benchmark")
	iterations := cmd.Int("iterations", 5, "number of warm queries to run")
	includeDirect := cmd.Bool("include-direct", true, "also benchmark the derived direct endpoint when current host uses Neon pooler")

	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing db-bench flags")
	}

	if *iterations <= 0 {
		log.Fatal().Int("iterations", *iterations).Msg("db-bench iterations must be greater than zero")
	}

	targets := []dbBenchTarget{
		{
			Name: "configured",
			Host: config.Envs.Postgres.Host,
		},
	}

	if *includeDirect {
		directHost := deriveDirectHost(config.Envs.Postgres.Host)
		if directHost != "" && directHost != config.Envs.Postgres.Host {
			targets = append(targets, dbBenchTarget{
				Name: "direct",
				Host: directHost,
			})
		}
	}

	results := make([]dbBenchResult, 0, len(targets))
	for _, target := range targets {
		result, err := benchmarkDatabaseTarget(target, *query, *iterations)
		if err != nil {
			log.Fatal().Err(err).Str("target", target.Name).Str("host", target.Host).Msg("Error while benchmarking database target")
		}

		results = append(results, result)
	}

	for _, result := range results {
		log.Info().
			Str("target", result.Name).
			Str("host", result.Host).
			Dur("connect_time", result.ConnectTime).
			Strs("warm_queries", formatDurations(result.QueryDurations)).
			Str("avg_warm_query", averageDuration(result.QueryDurations).String()).
			Msg("DB benchmark completed")
	}
}

type dbBenchTarget struct {
	Name string
	Host string
}

func benchmarkDatabaseTarget(target dbBenchTarget, query string, iterations int) (dbBenchResult, error) {
	cfg := secondarypostgres.Config{
		Username:        config.Envs.Postgres.Username,
		Password:        config.Envs.Postgres.Password,
		Database:        config.Envs.Postgres.Database,
		Host:            target.Host,
		Port:            config.Envs.Postgres.Port,
		SSLMode:         config.Envs.Postgres.SslMode,
		ChannelBinding:  config.Envs.Postgres.ChannelBinding,
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnMaxLifetime: config.Envs.DB.ConnMaxLifetime,
	}

	connectStartedAt := time.Now()
	db, err := secondarypostgres.New(cfg)
	if err != nil {
		return dbBenchResult{}, fmt.Errorf("connect %s: %w", target.Name, err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Error().Err(closeErr).Str("target", target.Name).Msg("Error while closing benchmark database connection")
		}
	}()

	result := dbBenchResult{
		Name:        target.Name,
		Host:        target.Host,
		ConnectTime: time.Since(connectStartedAt),
	}

	for i := 0; i < iterations; i++ {
		duration, err := runBenchmarkQuery(db, query)
		if err != nil {
			return dbBenchResult{}, fmt.Errorf("run warm query %d for %s: %w", i+1, target.Name, err)
		}

		result.QueryDurations = append(result.QueryDurations, duration)
	}

	return result, nil
}

func runBenchmarkQuery(db *sqlx.DB, query string) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Envs.DB.ConnectionTimeout)*time.Second)
	defer cancel()

	startedAt := time.Now()
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return time.Since(startedAt), nil
}

func deriveDirectHost(host string) string {
	if !strings.Contains(host, "-pooler") {
		return ""
	}

	return strings.Replace(host, "-pooler", "", 1)
}

func formatDurations(durations []time.Duration) []string {
	formatted := make([]string, 0, len(durations))
	for _, duration := range durations {
		formatted = append(formatted, duration.String())
	}

	return formatted
}

func averageDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	var total time.Duration
	for _, duration := range durations {
		total += duration
	}

	return total / time.Duration(len(durations))
}
