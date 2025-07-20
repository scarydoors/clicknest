#!/usr/bin/env bash

set -e

# Create temp files
RESULTS=$(mktemp)

# Cleanup on exit
trap "rm -f $RESULTS" EXIT

# Run benchmarks
go test -bench='Benchmark(Insert|Async|Batch)' -count=10 -benchmem ./internal/clickhouse | tee $RESULTS

# Compare with benchstat
echo "Comparing results..."
benchstat -row=.fullname $RESULTS
