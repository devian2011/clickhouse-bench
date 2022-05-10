# Benchmark

Expressions: 
All tables contains - 12 000 000 rows.
Schema: click_id (string), user_id(int64), brand_id(string), balance_before(float), balance_after(float), amount(float), created_at(datetime)

All benches make 10 requests for different dates
Benches:
1. Get users count for brand id in 1 day
```sql
SELECT COUNT(user_id), brand_id FROM user_balance WHERE created_at >= '%s 00:00:00' AND created_at <= '%s 23:59:59' GROUP BY brand_id
```
2. Get all amount sum for users in 3 days
```sql
SELECT SUM(amount), user_id FROM user_balance WHERE created_at >= '%s 00:00:00' AND created_at <= '%s 23:59:59'  GROUP BY user_id
```
3. Get all data in 4 days range
```sql
SELECT click_id, brand_id, balance_before, balance_after, amount, user_id, created_at FROM user_balance WHERE created_at >= '%s 00:00:00' AND created_at <= '%s 23:59:59'
```
4. Get users count for brand id in 1 day for three brands - 5,9,4
```sql
SELECT COUNT(user_id), brand_id FROM user_balance WHERE created_at >= '%s 00:00:00' AND created_at <= '%s 23:59:59' WHERE brand_id IN (5,9,4) GROUP BY brand_id
```

## Clickhouse bench

Table schema:
```clickhouse
CREATE TABLE IF NOT EXISTS user_balance (
			  click_id String
			, user_id Int64
			, balance_before Float32
			, balance_after Float32
			, brand_id String
			, amount Float32
			, created_at DateTime
		) Engine = MergeTree() ORDER BY (created_at)
```

Result: 
```
==================== ClickHouse ========================
Operation: GroupByBrandIdLastDay
Times: 10
Detailed: 115.4645ms 68.193ms 67.649375ms 80.618833ms 75.355916ms 64.620375ms 65.519542ms 64.615709ms 65.216917ms 62.865042ms
Avg Time: 73.01192ms
Operation: GroupByUserIdSumAmountLastThreeDays
Times: 10
Detailed: 102.009209ms 82.1165ms 78.507666ms 110.578708ms 99.339625ms 76.882375ms 82.794291ms 101.239459ms 81.43475ms 81.602584ms
Avg Time: 89.650516ms
Operation: SelectFourDays
Times: 10
Detailed: 3.419763125s 2.923913834s 3.0297315s 3.15992975s 3.152563208s 2.844815125s 3.248616625s 2.9074735s 2.99714925s 2.914957833s
Avg Time: 3.059891375s
Operation: GroupByBrandIdTwoDayForThreeBrands
Times: 10
Detailed: 3.242806833s 2.837851167s 2.848761666s 2.821099s 2.787955041s 3.015570958s 3.024556042s 2.826949375s 2.686266959s 2.72121s
Avg Time: 2.881302704s
=============================================================
```

## Postgres No Partition

Table schema:
```sql
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
```

Result: 
```
==================== PostgresNoPartition ========================
Operation: GroupByBrandIdLastDay
Times: 10
Detailed: 508.347458ms 256.73825ms 312.089083ms 211.320375ms 191.327792ms 204.528375ms 193.6085ms 240.072334ms 178.427958ms 188.875917ms
Avg Time: 248.533604ms
Operation: GroupByUserIdSumAmountLastThreeDays
Times: 10
Detailed: 1.967857916s 1.57899775s 1.565463375s 1.647264959s 1.70596525s 1.615416709s 1.381085791s 1.586139875s 1.379186459s 1.908756541s
Avg Time: 1.633613462s
Operation: SelectTwoDays
Times: 10
Detailed: 8.307667s 7.7812515s 7.863944709s 7.490528917s 7.523603542s 7.645440041s 7.01303925s 7.279980917s 6.549775459s 7.781795417s
Avg Time: 7.523702675s
Operation: GroupByBrandIdTwoDayForThreeBrands
Times: 10
Detailed: 7.495559166s 7.918881s 7.78897325s 7.761524417s 7.461660875s 7.66936125s 8.117216875s 7.818721917s 7.642375959s 8.377242208s
Avg Time: 7.805151691s
=============================================================
```

## Postgres With Partitions

Table schema:
```sql
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
```

Partitions had been divided by created_at and brand_id:
```sql
CREATE TABLE IF NOT EXISTS user_balance_[[date_t]] PARTITION OF user_balance FOR VALUES FROM ('[[date_from]]') TO ('[[date_to]]') PARTITION BY LIST (brand_id);
CREATE TABLE IF NOT EXISTS user_balance_[[date_t]]_[[brand]] PARTITION OF user_balance_[[date_t]] FOR VALUES IN ([[brand]]);
```

Result:
```
==================== PostgresPartition ========================
Operation: GroupByBrandIdTwoDayForThreeBrands
Times: 10
Detailed: 5.007945458s 5.136841708s 4.91011875s 4.903406959s 5.077984959s 4.999220792s 5.019627167s 5.074499625s 5.110480583s 5.158685167s
Avg Time: 5.039881116s
Operation: GroupByBrandIdLastDay
Times: 10
Detailed: 341.644291ms 44.058916ms 167.9795ms 180.5355ms 104.438583ms 147.456458ms 116.836542ms 98.779583ms 44.123459ms 43.494792ms
Avg Time: 128.934762ms
Operation: GroupByUserIdSumAmountLastThreeDays
Times: 10
Detailed: 385.830375ms 101.823917ms 255.547666ms 262.034125ms 217.134292ms 178.094833ms 213.80125ms 175.584625ms 118.422541ms 116.952209ms
Avg Time: 202.522583ms
Operation: SelectFourDays
Times: 10
Detailed: 5.600529542s 5.180947916s 5.026262084s 5.484189917s 5.492043166s 5.056289959s 5.399556s 5.02647175s 5.4596215s 5.184391041s
Avg Time: 5.291030287s
=============================================================
```
