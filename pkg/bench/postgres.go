package bench

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type PostgresConnection struct {
	conn *sqlx.DB
}

func NewPostgresConnection() *PostgresConnection {
	cn, err := sqlx.Open("postgres", "dbname=test host=localhost port=15432 user=test password=test sslmode=disable")
	if err != nil {
		log.Fatalln("Cannot connect to postgres. Err: " + err.Error())
	}
	cn.SetMaxIdleConns(5)
	cn.SetMaxOpenConns(5)
	cn.SetConnMaxIdleTime(time.Hour)
	cn.SetConnMaxLifetime(time.Hour)

	err = cn.Ping()
	if err != nil {
		log.Fatalln("Postgres ping error: " + err.Error())
	}

	return &PostgresConnection{
		conn: cn,
	}
}

func (chc *PostgresConnection) GroupByBrandIdLastDay(today string, table string) {
	now, _ := time.Parse("2006-01-02", today)
	yesterday := now.AddDate(0, 0, -1)
	yesterdayFormat := yesterday.Format("2006-01-02")

	r, qErr := chc.conn.Query(
		fmt.Sprintf(
			"SELECT COUNT(user_id), brand_id FROM %s WHERE created_at >= '%s 00:00:00' AND created_at <= '%s 23:59:59' GROUP BY brand_id",
			table, yesterdayFormat, yesterdayFormat))

	if qErr != nil {
		log.Println("Error for select by brand and last day Err: " + qErr.Error())
	}
	r.Close()
}

func (chc *PostgresConnection) GroupByUserIdSumAmountLastThreeDays(today string, table string) {
	now, _ := time.Parse("2006-01-02", today)
	toDate := now.AddDate(0, 0, -1)
	fromDate := now.AddDate(0, 0, -3)

	r, qErr := chc.conn.Query(
		fmt.Sprintf(
			"SELECT SUM(amount), user_id FROM %s WHERE created_at >= '%s 00:00:00' AND created_at <= '%s 23:59:59'  GROUP BY user_id",
			table, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")))
	if qErr != nil {
		log.Println("Error for select amount sum for user and last three days day Err: " + qErr.Error())
	}
	r.Close()
}

func (chc *PostgresConnection) SelectFourDays(today string, table string) {
	now, _ := time.Parse("2006-01-02", today)
	fromDate := now.AddDate(0, 0, -4)

	r, qErr := chc.conn.Query(
		fmt.Sprintf(
			"SELECT click_id, brand_id, balance_before, balance_after, amount, user_id, created_at FROM %s WHERE created_at >= '%s 00:00:00' AND created_at <= '%s 23:59:59'",
			table, fromDate.Format("2006-01-02"), now.Format("2006-01-02")))

	if qErr != nil {
		log.Println("Error for select all data for last ten days Err: " + qErr.Error())
	}
	r.Close()
}

func (p *PostgresConnection) Shutdown() {
	err := p.conn.Close()
	if err != nil {
		log.Fatalln("Cannot close postgres connection. Err: " + err.Error())
	}
}
