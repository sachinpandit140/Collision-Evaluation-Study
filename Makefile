.PHONY: test bench run clean

test:
	go test ./... -v -count=1

bench:
	go test -bench=. ./... -benchmem

run:
	go run ./cmd/bench/ -n 1000 -frames 100 -seed 42 -out results.csv

clean:
	rm -f *.csv
