package main

import (
	"context"
	"log"
	"reportTest/pkg/bench"
	"reportTest/pkg/froze"
)

func main() {
	ctx := context.Background()
	clickhouse := bench.NewClickHouseConnection(ctx)
	postgres := bench.NewPostgresConnection()

	defer func() {
		clickhouse.Close()
		postgres.Shutdown()
	}()
	log.Println("Bench has been started")
	benchClickHouse(clickhouse)
	benchPostgresNoPartition(postgres)
	benchPostgresPartition(postgres)
	log.Println("Testing stop")
}

func benchClickHouse(clickhouse *bench.ClickHouseConnection) {
	log.Println("\nClickhouse bench start\n")
	callBench(func() {
		clickhouse.GroupByBrandIdLastDay("2022-05-08")
	}, "Clickhouse-GroupByBrandIdLastDay")
	callBench(func() {
		clickhouse.GroupByUserIdSumAmountLastThreeDays("2022-05-08")
	}, "Clickhouse-GroupByUserIdSumAmountLastThreeDays")
	callBench(func() {
		clickhouse.SelectLastTenDays("2022-05-08")
	}, "Clickhouse-SelectLastTenDays")
	log.Println("\nClickhouse bench end\n")
}

func benchPostgresNoPartition(postgres *bench.PostgresConnection) {
	log.Println("\nPostgres no partition bench has been started\n")
	callBench(func() {
		postgres.GroupByBrandIdLastDay("2022-05-08", "user_balance_l")
	}, "PostgresNoPartition-GroupByBrandIdLastDay")
	callBench(func() {
		postgres.GroupByUserIdSumAmountLastThreeDays("2022-05-08", "user_balance_l")
	}, "PostgresNoPartition-GroupByUserIdSumAmountLastThreeDays")
	callBench(func() {
		postgres.SelectLastTenDays("2022-05-08", "user_balance_l")
	}, "PostgresNoPartition-SelectLastTenDays")
	log.Println("\nPostgres no partition bench end\n")
}

func benchPostgresPartition(postgres *bench.PostgresConnection) {
	log.Println("\nPostgres partition bench has been started\n")
	callBench(func() {
		postgres.GroupByBrandIdLastDay("2022-05-08", "user_balance")
	}, "PostgresPartition-GroupByBrandIdLastDay")
	callBench(func() {
		postgres.GroupByUserIdSumAmountLastThreeDays("2022-05-08", "user_balance")
	}, "PostgresPartition-GroupByUserIdSumAmountLastThreeDays")
	callBench(func() {
		postgres.SelectLastTenDays("2022-05-08", "user_balance")
	}, "PostgresPartition-SelectLastTenDays")
	log.Println("\nPostgres partition bench end\n")
}

func callBench(execution func(), operation string) {
	time := &froze.Froze{}
	log.Println("Operation: " + operation + " operation has been started")
	time.Start()
	execution()
	time.Stop()
	log.Println("Operation: " + operation + " took: " + time.GetDiff().String())
}
