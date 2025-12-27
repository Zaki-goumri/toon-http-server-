#!/bin/bash

echo "========================================"
echo "Starting Load Tests"
echo "========================================"
echo ""

# Check if servers are running
if ! curl -s http://localhost:8080/users > /dev/null; then
    echo "ERROR: TOON server not running on port 8080"
    echo "Start it with: go run cmd/toon-server/main.go"
    echo "           or: ./bin/toon-server"
    exit 1
fi

if ! curl -s http://localhost:8081/json/users > /dev/null; then
    echo "ERROR: JSON server not running on port 8081"
    echo "Start it with: go run cmd/json-server/main.go"
    echo "           or: ./bin/json-server"
    exit 1
fi

echo "✓ Both servers are running"
echo ""

# Run TOON load test
echo "========================================"
echo "Testing TOON Server (port 8080)"
echo "========================================"
k6 run tests/load/load_test_toon.js
echo ""

# Run JSON load test
echo "========================================"
echo "Testing JSON Server (port 8081)"
echo "========================================"
k6 run tests/load/load_test_json.js
echo ""

echo "========================================"
echo "Load Tests Complete!"
echo "========================================"
