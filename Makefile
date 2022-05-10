
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
	./benchmark


