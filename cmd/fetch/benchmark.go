package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"reportTest/pkg/bench"
	"reportTest/pkg/froze"
	"strings"
	"sync"
	"time"
)

type benchResult struct {
	name       string
	operations map[string][]*froze.Froze
}

func (br *benchResult) generateBenchResult() {
	fmt.Println("==================== " + br.name + " ========================")
	for name, o := range br.operations {
		fmt.Printf("Operation: %s\n", name)
		var all []string
		var durations time.Duration
		for _, v := range o {
			dur := v.GetDiff()
			all = append(all, dur.String())
			durations += dur
		}
		fmt.Printf("Times: %d\n", len(o))
		fmt.Printf("Detailed: %s\n", strings.Join(all, " "))
		fmt.Printf("Avg Time: %s\n", (durations / time.Duration(len(o))).String())
	}
	fmt.Println("=============================================================\n\n")
}

func main() {
	testType := flag.String("test", "all", "Test name - all, clickhouse, postgres, postgrespartial")
	flag.Parse()

	ctx := context.Background()
	clickhouse := bench.NewClickHouseConnection(ctx)
	postgres := bench.NewPostgresConnection()
	noPartionPostgres := bench.NewPostgresConnection()

	dates := generateRandDates()

	defer func() {
		clickhouse.Close()
		postgres.Shutdown()
		noPartionPostgres.Shutdown()
	}()
	log.Println("Bench has been started")
	wg := &sync.WaitGroup{}
	switch *testType {
	case "clickhouse":
		wg.Add(1)
		go benchClickHouse(clickhouse, wg, dates)
	case "postgres":
		wg.Add(1)
		go benchPostgresPartition(postgres, wg, dates)
	case "postgrespartial":
		wg.Add(1)
		go benchPostgresPartition(postgres, wg, dates)
	case "all":
		wg.Add(3)
		go benchClickHouse(clickhouse, wg, dates)
		go benchPostgresPartition(postgres, wg, dates)
		go benchPostgresNoPartition(noPartionPostgres, wg, dates)
	default:
		log.Fatalln("Unknown bench start type")
	}

	wg.Wait()
	log.Println("Testing end")
}

func benchClickHouse(clickhouse *bench.ClickHouseConnection, wg *sync.WaitGroup, dates []string) {
	benchResult := benchResult{
		name:       "ClickHouse",
		operations: make(map[string][]*froze.Froze, 0),
	}

	for _, date := range dates {
		benchResult.operations["GroupByBrandIdLastDay"] = append(
			benchResult.operations["GroupByBrandIdLastDay"],
			callBench(func() { clickhouse.GroupByBrandIdLastDay(date) }))

		benchResult.operations["GroupByUserIdSumAmountLastThreeDays"] = append(
			benchResult.operations["GroupByUserIdSumAmountLastThreeDays"],
			callBench(func() { clickhouse.GroupByUserIdSumAmountLastThreeDays(date) }))

		benchResult.operations["SelectFourDays"] = append(
			benchResult.operations["SelectFourDays"],
			callBench(func() { clickhouse.SelectFourDays(date) }))

		benchResult.operations["GroupByBrandIdTwoDayForThreeBrands"] = append(
			benchResult.operations["GroupByBrandIdTwoDayForThreeBrands"],
			callBench(func() { clickhouse.SelectFourDays(date) }))
	}
	benchResult.generateBenchResult()
	wg.Done()
}

func benchPostgresNoPartition(postgres *bench.PostgresConnection, wg *sync.WaitGroup, dates []string) {
	benchResult := benchResult{
		name:       "PostgresNoPartition",
		operations: make(map[string][]*froze.Froze, 0),
	}

	for _, date := range dates {
		benchResult.operations["GroupByBrandIdLastDay"] = append(
			benchResult.operations["GroupByBrandIdLastDay"],
			callBench(func() { postgres.GroupByBrandIdLastDay(date, "user_balance_l") }))

		benchResult.operations["GroupByUserIdSumAmountLastThreeDays"] = append(
			benchResult.operations["GroupByUserIdSumAmountLastThreeDays"],
			callBench(func() { postgres.GroupByUserIdSumAmountLastThreeDays(date, "user_balance_l") }))

		benchResult.operations["SelectTwoDays"] = append(
			benchResult.operations["SelectTwoDays"],
			callBench(func() { postgres.SelectFourDays(date, "user_balance_l") }))

		benchResult.operations["GroupByBrandIdTwoDayForThreeBrands"] = append(
			benchResult.operations["GroupByBrandIdTwoDayForThreeBrands"],
			callBench(func() { postgres.SelectFourDays(date, "user_balance_l") }))
	}
	benchResult.generateBenchResult()
	wg.Done()
}

func benchPostgresPartition(postgres *bench.PostgresConnection, wg *sync.WaitGroup, dates []string) {
	benchResult := benchResult{
		name:       "PostgresPartition",
		operations: make(map[string][]*froze.Froze, 0),
	}

	for _, date := range dates {
		benchResult.operations["GroupByBrandIdLastDay"] = append(
			benchResult.operations["GroupByBrandIdLastDay"],
			callBench(func() { postgres.GroupByBrandIdLastDay(date, "user_balance") }))

		benchResult.operations["GroupByUserIdSumAmountLastThreeDays"] = append(
			benchResult.operations["GroupByUserIdSumAmountLastThreeDays"],
			callBench(func() { postgres.GroupByUserIdSumAmountLastThreeDays(date, "user_balance") }))

		benchResult.operations["SelectFourDays"] = append(
			benchResult.operations["SelectFourDays"],
			callBench(func() { postgres.SelectFourDays(date, "user_balance") }))

		benchResult.operations["GroupByBrandIdTwoDayForThreeBrands"] = append(
			benchResult.operations["GroupByBrandIdTwoDayForThreeBrands"],
			callBench(func() { postgres.SelectFourDays(date, "user_balance") }))
	}
	benchResult.generateBenchResult()
	wg.Done()
}

func generateRandDates() []string {
	dates := make([]string, 10)
	now := time.Now()
	monthAgo := now.AddDate(0, -1, 0)
	min := monthAgo.Unix()
	delta := now.Unix() - min
	for c := 0; c < 10; c++ {
		dates[c] = randDate(min, delta)
	}

	return dates
}

func randDate(min int64, delta int64) string {
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0).Format("2006-01-02")
}

func callBench(execution func()) *froze.Froze {
	t := &froze.Froze{}
	t.Start()
	execution()
	t.Stop()
	return t
}
