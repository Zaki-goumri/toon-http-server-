# TOON HTTP Server

A Go HTTP server implementation comparing TOON (a custom text-based format) vs JSON serialization performance.

## Overview

This project implements two identical REST APIs:
- **TOON Server** (port 8080) - Custom TOON format serialization
- **JSON Server** (port 8081) - Standard JSON serialization

## Load Test Results

Full CRUD operations tested with k6: GET all, POST create, GET single, PUT update, DELETE.
Testing with up to 100 concurrent users over 3.5 minutes.

### TOON Server (port 8080)
```
HTTP Performance:
  ├─ Requests/sec:      281.99 req/s
  ├─ Avg duration:      27.84ms
  ├─ Median duration:   748.64µs
  ├─ P90 duration:      92.27ms
  ├─ P95 duration:      141.12ms
  └─ Max duration:      555.46ms

Operations:
  ├─ Total checks:      59,395 (100% passed)
  ├─ GET all users:     ✓ 11,879 requests
  ├─ POST create:       ✓ 11,879 requests
  ├─ GET single:        ✓ 11,879 requests
  ├─ PUT update:        ✓ 11,879 requests
  └─ DELETE:            ✓ 11,879 requests

Network:
  ├─ Data received:     9.2 GB (44 MB/s)
  └─ Data sent:         6.9 MB (33 kB/s)

Total: 59,395 requests, 11,879 iterations, 100% success rate
```

### JSON Server (port 8081)
```
HTTP Performance:
  ├─ Requests/sec:      290.41 req/s (3.0% faster)
  ├─ Avg duration:      21.27ms (23.6% faster)
  ├─ Median duration:   1.38ms (84.3% slower)
  ├─ P90 duration:      75.60ms (18.1% faster)
  ├─ P95 duration:      125.32ms (11.2% faster)
  └─ Max duration:      280.51ms (49.5% faster)

Operations:
  ├─ Total checks:      61,160 (100% passed)
  ├─ GET all users:     ✓ 12,232 requests
  ├─ POST create:       ✓ 12,232 requests
  ├─ GET single:        ✓ 12,232 requests
  ├─ PUT update:        ✓ 12,232 requests
  └─ DELETE:            ✓ 12,232 requests

Network:
  ├─ Data received:     18 GB (83 MB/s, 95.7% more data)
  └─ Data sent:         7.6 MB (36 kB/s)

Total: 61,160 requests, 12,232 iterations, 100% success rate
```

### Key Findings

1. **Throughput**: JSON handles 3% more requests/sec (290.41 vs 281.99)
   - JSON: 12,232 complete iterations
   - TOON: 11,879 complete iterations
   - JSON processed 353 more full CRUD cycles (2.97% more)

2. **Response Times**: JSON is faster at P90+ percentiles
   - Average: JSON 21.27ms vs TOON 27.84ms (**23.6% faster**)
   - P90: JSON 75.60ms vs TOON 92.27ms (**18.1% faster**)
   - P95: JSON 125.32ms vs TOON 141.12ms (**11.2% faster**)
   - Median: TOON 748µs vs JSON 1.38ms (TOON 45.7% faster for quick requests)

3. **Bandwidth Efficiency**: TOON uses **~49% less bandwidth**
   - TOON received: 9.2 GB
   - JSON received: 18 GB
   - TOON is significantly more efficient for data transfer

4. **Reliability**: Both formats achieved 100% success rate across all 5 operations

### Conclusion

**Performance Winner: JSON**
- 3% higher throughput (290 vs 282 req/s)
- 24% faster average response time
- 18% faster at P90, 11% faster at P95
- Better under high load

**Bandwidth Winner: TOON**
- ~49% less data transferred (9.2 GB vs 18 GB)
- Better for network-constrained environments
- Lower hosting costs for high-traffic APIs

**Use Cases:**
- **JSON**: High-performance APIs, low-latency requirements, modern web apps
- **TOON**: Bandwidth-limited scenarios, mobile apps, IoT devices, cost optimization

## Why is JSON Faster Despite Being Less Efficient?

TOON transfers 49% less data but JSON is still 24% faster. Here's why:

### 1. C-Level Optimized Parsers
- **Go's `encoding/json`**: Written in optimized Go with assembly-level optimizations for hot paths
- **Most frameworks use native JSON**: Python (C), Node.js (V8/C++), Rust (native), Java (JVM-optimized)
- **TOON**: Custom reflection-based implementation (~440 lines) that inspects types at runtime

### 2. Reflection is Expensive
TOON's encoder/decoder uses Go reflection heavily:
- `reflect.ValueOf()` - Creates reflection objects for every value
- `reflect.Type()` - Runtime type inspection
- Struct field iteration and tag parsing on every encode/decode
- Type switching and conversions for every field

JSON's native implementation avoids most reflection with pre-compiled code paths.

### 3. String Operations Overhead
TOON encoder performs more string operations:
- Custom escaping logic for special characters
- `strings.Repeat()` for indentation on every line
- `strings.Split()` and parsing in decoder
- Manual CSV-style formatting
- More memory allocations from string concatenation

### 4. Maturity and Optimization
- **JSON parsers**: 20+ years of optimization, battle-tested in production
- **Hardware optimization**: Modern CPUs have optimizations for JSON-like workloads
- **TOON**: First implementation, no optimization yet

### 5. Ecosystem Effects
- JSON benefits from OS-level caching, CDN optimization, HTTP compression
- Browsers and tools have native JSON support
- TOON requires custom tooling everywhere

### The Trade-off Table

| Factor | JSON | TOON |
|--------|------|------|
| Parsing Speed | Native C-level code | Reflection-based Go |
| CPU Usage | Low (optimized) | High (reflection overhead) |
| Memory Allocations | Optimized | Many string operations |
| Code Maturity | 20+ years | First implementation |
| Bandwidth | 18 GB transferred | 9.2 GB transferred (49% less) |
| Network Cost | Higher | Lower |
| Best Use Case | Performance-critical | Bandwidth-limited |

### Conclusion

**JSON wins on speed** because decades of optimization at the C/assembly level beats bandwidth savings. The CPU cost of reflection and custom parsing outweighs the network transfer time for smaller payloads.

**When TOON wins**: Slow networks (mobile 3G), expensive bandwidth (CDN costs), or when network I/O is the actual bottleneck, not CPU.

## Running the Servers

### TOON Server
```bash
go run main.go models.go toon.go
```

### JSON Server
```bash
go run main_json.go models.go
```

## Load Testing

```bash
# Run both load tests
./run_load_tests.sh

# Or run individually
k6 run load_test_toon.js
k6 run load_test_json.js
```

## API Endpoints

### TOON Server (port 8080)
- `GET /users` - Get all users
- `GET /users/:id` - Get user by ID
- `POST /users` - Create user
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

### JSON Server (port 8081)
- `GET /json/users` - Get all users
- `GET /json/users/:id` - Get user by ID
- `POST /json/users` - Create user
- `PUT /json/users/:id` - Update user
- `DELETE /json/users/:id` - Delete user

## Project Structure

```
toon-go-server/
├── cmd/
│   ├── toon-server/        # TOON server binary
│   │   └── main.go
│   └── json-server/        # JSON server binary
│       └── main.go
├── pkg/
│   ├── models/             # Data models
│   │   └── user.go
│   ├── toon/               # TOON encoder/decoder
│   │   └── encoder.go
│   └── handlers/           # HTTP handlers
│       └── handlers.go
├── tests/
│   └── load/               # k6 load tests
│       ├── load_test_toon.js
│       └── load_test_json.js
├── scripts/                # Utility scripts
│   ├── build.sh
│   └── run_load_tests.sh
├── benchmarks/             # Go benchmarks
└── bin/                    # Compiled binaries
```

## Quick Start (Updated Structure)

### Build Both Servers
```bash
./scripts/build.sh
```

### Run Servers

**TOON Server:**
```bash
go run cmd/toon-server/main.go
# OR
./bin/toon-server
```

**JSON Server:**
```bash
go run cmd/json-server/main.go
# OR
./bin/json-server
```

### Load Testing
```bash
# Ensure both servers are running, then:
./scripts/run_load_tests.sh
```

## Development

```bash
# Format code
go fmt ./...

# Run tests
go test ./...

# Build binaries
./scripts/build.sh
```
