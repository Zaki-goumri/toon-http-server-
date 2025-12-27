TOON vs JSON Server Benchmark
==============================

A professional Go HTTP server implementation comparing TOON (custom text-based format) 
vs JSON serialization performance.

## Quick Commands

Build:
  ./scripts/build.sh

Run TOON Server:
  go run cmd/toon-server/main.go
  OR
  ./bin/toon-server

Run JSON Server:
  go run cmd/json-server/main.go
  OR
  ./bin/json-server

Load Test:
  ./scripts/run_load_tests.sh

Full documentation: README.md
