package main

import (
	"context"
	"log"
	"reportTest/pkg/faker"
	"reportTest/pkg/seeder"
	"sync"
)

func main() {
	ctx := context.Background()
	clickHouseCh := make(chan faker.UserData, 1000)
	postgresCh := make(chan faker.UserData, 1000)
	wg := &sync.WaitGroup{}

	clickHouse := seeder.NewClickHouseConnection(ctx)
	clickHouse.CreateTable()

	postgres := seeder.NewPostgresConnection()
	postgres.CreateTable()

	defer func() {
		clickHouse.Close()
		postgres.Shutdown()
	}()

	wg.Add(1)
	go clickHouse.InsertData(wg, clickHouseCh)
	wg.Add(1)
	go postgres.InsertData(wg, postgresCh)

	out := []chan<- faker.UserData{clickHouseCh, postgresCh}

	faker.GenerateFakeData(1, 12_000_000, 10, 1000, out)

	wg.Wait()
	log.Println("Data has been loaded")

}
