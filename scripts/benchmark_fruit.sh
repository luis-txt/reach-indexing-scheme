#!/bin/bash

echo "================="
echo "Benching FRUIT..."
echo "================="
echo ""
echo "==== Creating test data..."
python3 scripts/create_test_data.py
echo "==== Running benchmarks on test data..."
python3 scripts/run_benches.py
echo "==== Creating plots..."
python3 scripts/create_plots.py
echo "==== Creating tables..."
python3 scripts/create_tables.py
