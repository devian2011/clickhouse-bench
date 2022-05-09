package seeder

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"log"
	"reportTest/pkg/faker"
	"sync"
	"time"
)

type ClickHouseConnection struct {
	ctx  context.Context
	conn driver.Conn
}

func NewClickHouseConnection(ctx context.Context) *ClickHouseConnection {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:19000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		//Debug:           true,
		DialTimeout:     time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		log.Fatalf("Cannot connect to clichouse")
	}

	return &ClickHouseConnection{conn: conn, ctx: ctx}
}

func (chc *ClickHouseConnection) CreateTable() {
	var err error
	if err = chc.conn.Exec(chc.ctx, `DROP TABLE IF EXISTS user_balance`); err != nil {
		log.Fatalln("Cannot drop clickhouse table user_balance. Err: " + err.Error())
	}
	err = chc.conn.Exec(chc.ctx, `
		CREATE TABLE IF NOT EXISTS user_balance (
			  click_id String
			, user_id Int64
			, balance_before Float32
			, balance_after Float32
			, brand_id String
			, amount Float32
			, created_at DateTime
		) Engine = MergeTree() ORDER BY (created_at)
	`)
	if err != nil {
		log.Fatalln("Cannot create clickhouse table user_balance. Err: " + err.Error())
	}
}

func (chc *ClickHouseConnection) InsertData(wg *sync.WaitGroup, insert <-chan faker.UserData) {
	var err error
	var batch driver.Batch
	batchCounter := 0
	for v := range insert {
		if batchCounter == 0 || batchCounter%1_000 == 0 {
			batch, err = chc.conn.PrepareBatch(chc.ctx, "INSERT INTO user_balance")
			if err != nil {
				log.Println("Fail to prepare insertion data")
			}
		}
		batchCounter++
		err := batch.Append(v.ClickId, v.UserId, v.BalanceBefore, v.BalanceAfter, v.BrandId, v.Amount, v.CreatedAt)
		if err != nil {
			log.Println("Fail append data to batch. Err: " + err.Error())
		}
		if batchCounter == 1_000 {
			batchCounter = 0
			err = batch.Send()
			if err != nil {
				log.Println("Fail to send batch. Err: " + err.Error())
			}

			log.Println("Write 1000 rows to clickhouse")
		}

	}
	err = batch.Send()
	if err != nil {
		log.Println("Fail to send batch. Err: " + err.Error())
	}
	wg.Done()
}

func (chc *ClickHouseConnection) Close() {
	err := chc.conn.Close()
	if err != nil {
		log.Println("Cannot close clickhouse connection. Err: " + err.Error())
	}
}
