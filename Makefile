
.PHONY: seed
seed: build_seeder run_seeder

.PHONY: bench
bench: build_bench run_bench

.PHONY: build_seeder
build_seeder:
	go build -o seeder ./cmd/seeder/main.go

.PHONY: run_seeder
run_seeder:
	./seeder

.PHONY: build_bench
build_bench:
	go build -o benchmark ./cmd/fetch/benchmark.go

.PHONY: run_bench
run_bench:
	./benchmark -test=all

.PHONY: run_bench_clickhouse
run_bench_clickhouse: build_bench
	./benchmark -test=clickhouse

.PHONY: run_bench_postgres
run_bench_postgres: build_bench
	./benchmark -test=postgres

.PHONY: run_bench_postgrespartial
run_bench_postgrespartial: build_bench
	./benchmark -test=postgrespartial
