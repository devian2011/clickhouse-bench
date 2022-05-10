package seeder

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"reportTest/pkg/faker"
	"strings"
	"sync"
	"time"
)

var schema = `
DROP TABLE IF EXISTS user_balance_l;
CREATE TABLE IF NOT EXISTS user_balance_l (
	id bigserial not null,
	click_id varchar(100),
	brand_id varchar(100),
	balance_before decimal(10,2),
	balance_after decimal(10,2),
	amount decimal(10,2),
	user_id bigint,
	created_at timestamp WITHOUT TIME ZONE
);
CREATE INDEX user_balance_l_created_at_brand ON user_balance_l USING btree (created_at, brand_id);
CREATE INDEX user_balance_l_brand_id ON user_balance_l using btree (brand_id);
CREATE INDEX user_balance_l_user_id ON user_balance_l USING btree (user_id);

DROP TABLE IF EXISTS user_balance;
CREATE TABLE IF NOT EXISTS user_balance (
	id bigserial not null,
	click_id varchar(100),
	brand_id varchar(100),
	balance_before decimal(10,2),
	balance_after decimal(10,2),
	amount decimal(10,2),
	user_id bigint,
	created_at timestamp WITHOUT TIME ZONE
) PARTITION BY RANGE (created_at);
CREATE INDEX user_balance_created_at_brand ON user_balance USING btree (created_at, brand_id);
CREATE INDEX user_balance_brand_id ON user_balance using btree (brand_id);
CREATE INDEX user_balance_user_id ON user_balance USING btree (user_id);
`

var createPartitionSql = `
CREATE TABLE IF NOT EXISTS user_balance_[[date_t]] PARTITION OF user_balance FOR VALUES FROM ('[[date_from]]') TO ('[[date_to]]') PARTITION BY LIST (brand_id);
CREATE TABLE IF NOT EXISTS user_balance_[[date_t]]_[[brand]] PARTITION OF user_balance_[[date_t]] FOR VALUES IN ([[brand]]);
`

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
	cn.SetConnMaxLifetime(time.Hour)

	err = cn.Ping()
	if err != nil {
		log.Fatalln("Postgres ping error: " + err.Error())
	}

	return &PostgresConnection{
		conn: cn,
	}
}

func (p *PostgresConnection) CreateTable() {
	p.conn.MustExec(schema)
}

func (p *PostgresConnection) InsertData(wg *sync.WaitGroup, insert <-chan faker.UserData) {
	batchCounter := 0
	userDataRows := make([]faker.UserData, 0)
	for v := range insert {
		batchCounter++
		userDataRows = append(userDataRows, v)

		crPartitionSql := strings.ReplaceAll(createPartitionSql, "[[date_from]]", v.CreatedAt.Format("2006-01-02")+" 00:00:00")
		crPartitionSql = strings.ReplaceAll(crPartitionSql, "[[date_to]]", v.CreatedAt.Format("2006-01-02")+" 23:59:59")
		crPartitionSql = strings.ReplaceAll(crPartitionSql, "[[date_t]]", v.CreatedAt.Format("2006_01_02"))
		crPartitionSql = strings.ReplaceAll(crPartitionSql, "[[brand]]", strings.ReplaceAll(v.BrandId, "-", ""))

		p.conn.MustExec(crPartitionSql)

		if batchCounter == 100 {
			_, err := p.conn.NamedExec("INSERT INTO user_balance (click_id, brand_id, balance_before, balance_after, amount, user_id, created_at) VALUES (:click_id, :brand_id, :balance_before, :balance_after, :amount, :user_id, :created_at)", userDataRows)
			if err != nil {
				log.Println("Fail to insert batch to postgres user_balance. Err: " + err.Error())
			}
			_, err = p.conn.NamedExec("INSERT INTO user_balance_l (click_id, brand_id, balance_before, balance_after, amount, user_id, created_at) VALUES (:click_id, :brand_id, :balance_before, :balance_after, :amount, :user_id, :created_at)", userDataRows)
			if err != nil {
				log.Println("Fail to insert batch to postgres user_balance_l. Err: " + err.Error())
			}
			batchCounter = 0
			userDataRows = nil
			userDataRows = make([]faker.UserData, 0)
			log.Println("Write 100 rows to postgres")
		}
	}

	wg.Done()
}

func (p *PostgresConnection) Shutdown() {
	err := p.conn.Close()
	if err != nil {
		log.Fatalln("Cannot close postgres connection. Err: " + err.Error())
	}
}
