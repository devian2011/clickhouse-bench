package bench

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"log"
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

func (chc *ClickHouseConnection) GroupByBrandIdLastDay(today string) {
	now, _ := time.Parse("2006-01-02", today)
	yesterday := now.AddDate(0, 0, -1)
	yesterdayFormat := yesterday.Format("2006-01-02")

	_, qErr := chc.conn.Query(chc.ctx,
		fmt.Sprintf(
			"SELECT COUNT(user_id), brand_id FROM user_balance WHERE created_at >= '%s 00:00:00' AND created_at <= '%s 23:59:59' GROUP BY brand_id",
			yesterdayFormat, yesterdayFormat))

	if qErr != nil {
		log.Println("Error for select by brand and last day Err: " + qErr.Error())
	}
}

func (chc *ClickHouseConnection) GroupByUserIdSumAmountLastThreeDays(today string) {
	now, _ := time.Parse("2006-01-02", today)
	toDate := now.AddDate(0, 0, -1)
	fromDate := now.AddDate(0, 0, -3)

	_, qErr := chc.conn.Query(chc.ctx,
		fmt.Sprintf(
			"SELECT SUM(amount), user_id FROM user_balance WHERE created_at >= '%s 00:00:00' AND created_at <= '%s 23:59:59'  GROUP BY user_id",
			fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")))

	if qErr != nil {
		log.Println("Error for select amount sum for user and last three days day Err: " + qErr.Error())
	}
}

func (chc *ClickHouseConnection) SelectLastTenDays(today string) {
	now, _ := time.Parse("2006-01-02", today)
	fromDate := now.AddDate(0, 0, -11)

	_, qErr := chc.conn.Query(chc.ctx,
		fmt.Sprintf(
			"SELECT click_id, brand_id, balance_before, balance_after, amount, user_id, created_at FROM user_balance WHERE created_at >= '%s 00:00:00'",
			fromDate.Format("2006-01-02")))

	if qErr != nil {
		log.Println("Error for select all data for last ten days Err: " + qErr.Error())
	}
}

func (chc *ClickHouseConnection) Close() {
	err := chc.conn.Close()
	if err != nil {
		log.Println("Cannot close clickhouse connection. Err: " + err.Error())
	}
}
