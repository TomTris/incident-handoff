# Go Mastery — Real-World Challenge Handbook

> A hands-on curriculum for developers who want to go from Go beginner to production-ready engineer.
> Built for people who learn by doing — not by reading slides.

---

## Who is this for?

This handbook is for developers who:
- Already know basic of programming
- Want to learn Go properly — not just syntax, but how Go is actually used in production
- Want to understand Go backend and systems engineering deeply

## What you'll be able to do after this

By the end of these challenges you will be able to:
- Design systems using Go's interface model (the way the standard library does it)
- Write concurrent programs without race conditions
- Profile and optimize Go code for production workloads
- Build HTTP services that handle real traffic patterns
- Read and review Go code like a advanced engineer
- Understand diverse Go topics such as concurrency model, memory model, and type system

## Skill progression

| After Phase | What you can do |
|---|---|
| Phase 1 ✅ | Write basic Go — structs, slices, methods, file I/O |
| Phase 2 | Design with interfaces, read standard library code |
| Phase 3 | Build concurrent systems without race conditions |
| Phase 4 | Profile and optimize Go code for performance |
| Phase 5 | Ship a production-grade HTTP service end-to-end |
| Phase 6 | Read, review, and debug Go code under pressure |

---

# PHASE 1 — Foundations ✅
Check

---

# PHASE 2 — Interfaces & Type System

> **Why this phase matters**
> Interfaces are the backbone of Go's entire standard library. `http.Handler`, `io.Reader`, `io.Writer`, `error` — all interfaces. If you don't understand interfaces deeply, you'll struggle to read Go code written by others, and you'll write brittle code yourself. This phase teaches you to think the Go way: *program to behavior, not to concrete types.*

---

## Challenge 2.1 — Build a Multi-Format Logger
### `🟡 Beginner → Intermediate`
**🕐 Expected duration: 8–10 hours**

### 1. Context
Every production system logs events. But *where* those logs go changes depending on the environment: console during local development, structured files in staging, JSON for log aggregation tools like Datadog, Grafana Loki, or AWS CloudWatch in production.

A well-designed logging system should let you swap the destination without changing the code that *uses* the logger. This is exactly how Go's `io.Writer`, `log/slog`, and popular libraries like `uber-go/zap` work internally.

### 2. Goal
Build a logging system that can write to multiple destinations (console, file, JSON) using a common interface. The logger must be swappable — the rest of the code should not care where logs go.

### 3. Scope
- Define a `Logger` interface with at least one method: `Log(level, message string)`
- Implement 3 concrete loggers, all satisfying the `Logger` interface:
  - `ConsoleLogger` — prints to terminal with timestamp
  - `FileLogger` — writes to a `.log` file
  - `JSONLogger` — writes structured JSON lines to a file
- Write a function `RunApp(l Logger)` that takes any logger and logs 3 events (startup, a warning, a shutdown)
- In `main()`, call `RunApp` three times — once with each logger type
- No `if/else` based on logger type anywhere in `RunApp` — it must work purely through the interface

### 4. Expected Output
Console:
```
[2026-03-22 10:00:01] INFO  app started
[2026-03-22 10:00:01] WARN  high memory usage
[2026-03-22 10:00:01] INFO  app shutdown
```
File (`app.log`):
```
2026-03-22 10:00:01 | INFO  | app started
2026-03-22 10:00:01 | WARN  | high memory usage
2026-03-22 10:00:01 | INFO  | app shutdown
```
JSON (`app.json`):
```json
{"time":"2026-03-22T10:00:01","level":"INFO","message":"app started"}
{"time":"2026-03-22T10:00:01","level":"WARN","message":"high memory usage"}
{"time":"2026-03-22T10:00:01","level":"INFO","message":"app shutdown"}
```

### 5. Hints & Knowledge
- In Go, interfaces are implemented **implicitly** — no `implements` keyword. If your struct has the right methods, it satisfies the interface automatically.
- `io.Writer` is Go's most important interface: `Write(p []byte) (n int, err error)`. `os.Stdout` and `os.File` both implement it — that's why you can write to both the same way.
- `time.Now().Format("2006-01-02 15:04:05")` — Go uses a reference time for formatting (Jan 2, 2006 = Go's birthday).
- `encoding/json` — use `json.Marshal(struct)` to convert a struct to JSON bytes.
- `os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)` — open a file for appending.

### 6. Sources
- Go interfaces explained: https://go.dev/tour/methods/9
- `io.Writer`: https://pkg.go.dev/io#Writer
- `encoding/json`: https://pkg.go.dev/encoding/json
- `time.Format`: https://pkg.go.dev/time#Time.Format
- `os.OpenFile`: https://pkg.go.dev/os#OpenFile

### 7. Knowledge Gained
- ✅ How Go interfaces work (implicit implementation)
- ✅ Writing to `io.Writer` — the foundation of all Go I/O
- ✅ `encoding/json` for structured data
- ✅ Dependency injection via interfaces (pass behavior, not implementation)
- ✅ The design pattern used by `net/http`, `os`, `bufio`, and most Go packages

---

## Challenge 2.2 — Fix the Shape Calculator
### `🟢 Beginner`
**🕐 Expected duration: 3–4 hours**

### 1. Context
A junior developer tried to build a geometry calculator that computes the area of different shapes using Go interfaces. The code compiles in some places and panics in others. Your job is to fix it and make it robust.

### 2. Goal
Fix all bugs in the provided broken code, understand why each bug exists, and add one defensive improvement using a **type switch**.

### 3. Scope
Here is the broken code:

```go
package main

import (
    "fmt"
    "math"
)

type Shape interface {
    Area() float64
    Describe() string
}

type Circle struct {
    Radius float64
}

type Rectangle struct {
    Width, Height float64
}

type Triangle struct {
    Base, Height float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Describe() string {
    return fmt.Sprintf("Rectangle %.1f x %.1f", r.Width, r.Height)
}

func (c Circle) Describe() string {
    return fmt.Sprintf("Circle r=%.1f", c.Radius)
}

func printArea(s Shape) {
    fmt.Printf("%s → area: %.2f\n", s.Describe, s.Area)
}

func totalArea(shapes []Shape) float64 {
    total := 0
    for _, s := range shapes {
        total += s.Area()
    }
    return total
}

func main() {
    shapes := []Shape{
        Circle{Radius: 5},
        Rectangle{Width: 3, Height: 4},
        Triangle{Base: 6, Height: 8},
    }
    for _, s := range shapes {
        printArea(s)
    }
    fmt.Printf("Total area: %.2f\n", totalArea(shapes))
}
```

Find ALL bugs (there are 5), fix them, then add:
- A `type switch` inside `printArea` that prints `"(is a circle)"` if the shape is a `Circle`

### 4. Expected Output
```
Circle r=5.0 (is a circle) → area: 78.54
Rectangle 3.0 x 4.0 → area: 12.00
Triangle b=6.0 h=8.0 → area: 24.00
Total area: 114.54
```

### 5. Hints & Knowledge
- Missing method on a type = does NOT satisfy the interface → compile error
- `s.Describe` vs `s.Describe()` — calling a method needs `()`
- `total := 0` makes `total` an `int` — can't add `float64` to it
- Type switch syntax: `switch v := s.(type) { case Circle: ... }`
- A `Triangle` must implement ALL methods of `Shape` to be used as one

### 6. Sources
- Type assertions: https://go.dev/tour/methods/15
- Type switches: https://go.dev/tour/methods/16

### 7. Knowledge Gained
- ✅ Interface satisfaction rules — ALL methods must be implemented
- ✅ Type assertion and type switch
- ✅ Common interface bugs and how to spot them
- ✅ Zero values and type mismatches

---

# PHASE 3 — Goroutines & Channels

> **Why this phase matters**
> Concurrency is Go's killer feature. Goroutines are why Go is chosen popular. A goroutine costs ~2KB of memory vs ~8MB for a thread — you can run hundreds of thousands of them. Channels replace shared memory with message passing, eliminating entire classes of bugs. If you can write correct concurrent Go, you are got what is nice of Go.

---

## Challenge 3.1A — Log File Generator
### `🟢 Beginner`
**🕐 Expected duration: 2–3 hours**

### 1. Context
Before you can build a concurrent log scanner, you need log files to scan. This part generates realistic fake logs that Part B will consume.

### 2. Goal
Build `log_generator.go` that creates 7 fake `.log` files in `./logs/`, each with 300 lines of random log entries.

### 3. Scope
- Create `./logs/` directory safely (no crash if it already exists)
- Generate files named `gateway_1.log` through `auth-service_7.log`, by using service names (See [Given Variables](#4-given-variables)).
- Each line follows this format:
```
2026-05-01T08:00:03.000Z [INFO ] [gateway] [trace=a3f9c012] Request processed in 42ms
```
- Level distribution: 70% INFO, 20% WARN, 10% ERROR
- Each file uses a different service name
- Timestamps must strictly increase within each file, each message is spaced 1-15 seconds apart.
- Print progress as each file is created

### 4. Given Variables
```go
const (
	numFiles    = 7
	numLines    = 300
	outputDir   = "logs"
)
 
var (
	infoMessages = []string{
		"Service started successfully",
		"Configuration loaded from /etc/app/config.yaml",
		"Database connection pool initialized (size=10)",
		"Health check endpoint responding on :8080/health",
		"Cache warmed up with 1482 entries",
		"Scheduled job 'cleanup' registered (interval=5m)",
		"TLS certificate valid until 2027-03-15",
		"Request processed in 42ms",
		"User session created",
		"Metrics exported to Prometheus",
		"Batch import completed: 350 records",
		"Webhook delivered to https://hooks.example.com/notify",
		"Rate limiter reset for tenant acme-corp",
		"Graceful reload triggered by SIGHUP",
		"New worker spawned (pool=4/8)",
	}
 
	warnMessages = []string{
		"Response time exceeded 500ms threshold (actual=623ms)",
		"Disk usage at 82%% on /var/lib/data",
		"Retry attempt 2/5 for upstream service payment-api",
		"Deprecated header X-Request-ID used by client 10.0.3.44",
		"Connection pool nearing capacity (8/10)",
		"Clock skew detected: 1.3s drift from NTP server",
		"Certificate expires in 30 days",
		"Memory usage above 75%% (current=78%%)",
		"Slow query detected: SELECT * FROM orders (1.8s)",
		"Rate limit approaching for API key sk-...a3f9",
		"Fallback to secondary DNS resolver",
		"Stale cache entry served for key user:9281",
		"Config key 'legacy_mode' is deprecated, migrate to v2 schema",
		"Unrecognized query parameter 'debug' ignored",
		"Partial response returned: 3 of 5 shards responded",
	}
 
	errorMessages = []string{
		"Failed to connect to postgres://db:5432/app — connection refused",
		"Panic recovered in handler /api/v1/orders: index out of range [5]",
		"TLS handshake failed: certificate signed by unknown authority",
		"Out of memory: cannot allocate 256MB for image processing",
		"Deadlock detected between goroutine 47 and goroutine 52",
		"Kafka consumer lag exceeded 10000 messages on topic=events",
		"Disk write failed: /var/log/app.log — no space left on device",
		"Authentication failed for user admin@example.com (attempt 5/5)",
		"Circuit breaker OPEN for service inventory-api (failures=12)",
		"Unhandled exception in middleware chain: nil pointer dereference",
		"DNS resolution failed for api.partner.io",
		"S3 upload failed: AccessDenied on bucket prod-assets",
		"Request timeout after 30s: POST /api/v1/reports/generate",
		"Invalid JWT token: signature verification failed",
		"Migration 0042_add_index.sql failed: duplicate column name",
	}
 
	services = []string{
		"gateway", "auth-service", "order-engine",
		"payment-api", "notification-worker", "scheduler", "inventory-sync",
	}
)
```

And you can use "INFO", "WARN", "ERROR" as 3 levels.
### 5. Expected Output
```
Generating 7 log files (300 lines each)...
  ✔ logs/gateway-1.log  (300 lines)
  ✔ logs/auth-service-2.log  (300 lines)
  ✔ logs/order-engine-3.log  (300 lines)
  ✔ logs/payment-api-4.log  (300 lines)
  ✔ logs/notification-worker-5.log  (300 lines)
  ✔ logs/scheduler-6.log  (300 lines)
  ✔ logs/inventory-sync-7.log  (300 lines)
Done!
```

### 6. Hints
- `os.MkdirAll`, `os.Create`, `fmt.Sprintf("%08x", rand.Uint32())`
- Accumulate timestamps, don't calculate from base + index

### 7. Knowledge Gained
- ✅ File I/O — `os.Create`, `WriteString`
- ✅ `time.Time` formatting and arithmetic
- ✅ Weighted random generation

---
---

## Challenge 3.1B — Concurrent Log Scanner
### `🟡 Beginner → Intermediate`
**🕐 Expected duration: 8–12 hours**

### 1. Context
You have 5 log files from Part A. An incident just happened — scan all files simultaneously, count levels per file, capture every ERROR line, and print a report. Doing it sequentially is too slow.

### 2. Goal
Build `scanner.go` that processes all log files in parallel using goroutines, collects results through a channel AND a mutex-protected shared list, and prints a sorted report.

### 3. Scope
- Read all `.log` files from `./logs/`
- Launch **one goroutine per file** — fan-out
- Each worker scans its file line by line, counts INFO/WARN/ERROR
- Workers send per-file statistics through a **channel** — fan-in
- Workers append individual ERROR lines to a **shared slice protected by mutex**
- Use `sync.WaitGroup` to know when all workers are done
- Close the channel after all goroutines finish
- Print a sorted report: table of counts per file, totals, worst file, and first 5 ERROR lines
- Must pass `go run -race .` with zero race conditions

### 4. Given Structs
```go
type FileResult struct {
    Filename   string
    InfoCount  int
    WarnCount  int
    ErrorCount int
}

type ErrorLog struct {
    Filename string
    Line     string
}
```

```go
var (
    allErrors []ErrorLog
    mu        sync.Mutex
)
```

One struct goes through a channel. The other goes into a shared slice. Which is which — and why — is yours to figure out.

### 5. Must Use
```
goroutine (1 per file)       — fan-out
channel (for per-file stats) — fan-in
sync.WaitGroup               — coordinate completion
sync.Mutex                   — protect shared error list
bufio.Scanner                — read files line by line
```

### 6. Expected Output
```
Scanning 7 files concurrently...

[worker] gateway-1.log → INFO:138  WARN:41  ERROR:21
[worker] order-engine-3.log → INFO:140  WARN:39  ERROR:21
[worker] auth-service-2.log → INFO:142  WARN:38  ERROR:20
[worker] notification-worker-5.log → INFO:137  WARN:42  ERROR:21
[worker] payment-api-4.log → INFO:141  WARN:40  ERROR:19
[worker] scheduler-6.log → INFO:139  WARN:41  ERROR:20
[worker] inventory-sync-7.log → INFO:140  WARN:40  ERROR:20

══════════════════════════════════════════════════════
                    ERROR REPORT
══════════════════════════════════════════════════════
 File                          INFO    WARN    ERROR
──────────────────────────────────────────────────────
 gateway-1.log                  138      41       21
 auth-service-2.log             142      38       20
 order-engine-3.log             140      39       21
 payment-api-4.log              141      40       19
 notification-worker-5.log      137      42       21
 scheduler-6.log                139      41       20
 inventory-sync-7.log           140      40       20
──────────────────────────────────────────────────────
 TOTAL                          977     281      142
══════════════════════════════════════════════════════

Worst file: notification-worker-5.log (21 errors)

First 5 ERROR lines:
  [gateway-1.log] 2026-05-01T08:02:11.000Z [ERROR] Failed to connect...
  [order-engine-3.log] 2026-05-01T08:01:02.000Z [ERROR] Out of memory...
```

Note: 
- `[worker]` lines appear in random order (proves concurrency). Report table is sorted by filename.
- The 5 ERROR lines may differ (proves concurrency)

### 7. Design Constraint — Why Both Channel AND Mutex?
You must use **both** mechanisms in this challenge. Think about what kind of data each worker produces:
- One type is sent **exactly once** per worker
- The other is sent **zero to many times**, unpredictably

Which mechanism fits which? Figure this out — it's the core design decision.

### 8. Why This Matters in Production
Fan-out/fan-in is the core of Go's concurrency model. It's used in CI/CD systems (test N packages in parallel), data pipelines (process N files simultaneously), web scrapers (fetch N URLs concurrently), and Kubernetes controllers (reconcile N resources at the same time).

Knowing when to reach for a channel vs a mutex is what separates junior from mid-level Go engineers.

### 9. Common Mistakes to Avoid
- `wg.Add(1)` in the wrong place — silent data loss, no crash, no error, just missing results
- `wg.Wait(); close(ch)` in the wrong goroutine — deadlock
- Forgetting `close(ch)` entirely — deadlock
- Appending to a shared slice without protection — race condition
- Passing `wg` instead of `&wg` — `Done()` on a copy does nothing

### 10. Hints & Knowledge
- `os.ReadDir`, `bufio.NewScanner`, `strings.Contains`
- `sort.Slice` for sorting results
- `fmt.Fprintf` with `%-16s %6d` for aligned columns
- Buffered channel: `make(chan T, N)` — workers don't block waiting for reader
- `for result := range ch` — reads until `close(ch)`

### 11. Sources
- Goroutines: https://go.dev/tour/concurrency/1
- Channels: https://go.dev/tour/concurrency/2
- `sync.WaitGroup`: https://pkg.go.dev/sync#WaitGroup
- `sync.Mutex`: https://pkg.go.dev/sync#Mutex
- Race detector: https://go.dev/doc/articles/race_detector
- `bufio.Scanner`: https://pkg.go.dev/bufio#Scanner

### 12. Checklist
```
[ ] go run scanner.go          — runs without errors
[ ] go run -race scanner.go   — ZERO race conditions
[ ] go vet ./...               — zero warnings
[ ] Worker lines appear in random order
[ ] Report table sorted by filename
[ ] ERROR lines are captured and printed
[ ] Both channel AND mutex are used correctly
```

### 14. Knowledge Gained
```
✅ Goroutine         — launch and manage concurrent work
✅ Channel           — buffered, send/receive, close, range
✅ Directional chan  — chan<- (send-only) in function signatures
✅ sync.WaitGroup    — coordinate goroutine completion
✅ sync.Mutex        — protect shared state
✅ Fan-out / Fan-in  — most important Go concurrency pattern
✅ Channel vs Mutex  — when to use which
✅ Race detection    — go run -race
✅ bufio.Scanner     — efficient line-by-line file reading
```

---

## Challenge 3.2 — Build a Worker Pool URL Checker
### `🟠 Intermediate`
**🕐 Expected duration: 15–20 hours**

### 0. Introduction

In production, you never spawn unlimited goroutines. If 10,000 requests come in and you launch 10,000 goroutines to handle them simultaneously, your server runs out of memory. The solution: a **worker pool** — a fixed number of goroutines that pick jobs from a queue, process them, and stay alive for the next job.

Worker pools are how Go HTTP servers, job queues (like Faktory, Asynq), and data processors work under the hood. This is the second of the two most important Go concurrency patterns.

### 1. Context
Your team runs 20 internal services. Every minute, an automated health checker pings each service's /health endpoint and reports which ones are down. The previous engineer wrote it in Python — it checks services one by one, sequentially. On a bad day with 6 timed-out services (each taking 4 seconds to timeout), it takes 24 seconds per cycle.
Your job: rewrite it in Go using a **worker pool** with the **minimum number of workers possible**. All 20 services must be checked and complete within **7 seconds**, even on a bad day when 6 services time out. Under normal conditions (no timeouts or network delay), in avarage for each url, the checker should finish well under 500ms second.

### 2. Goal
Build a URL health checker using a fixed worker pool of minimum number of goroutines possible that processes 20 URLs with timeout control.

### 3. Scope
- Define **minimun worker goroutines**
- Feed 20 URLs as below as in this order into a jobs channel (mix of valid/invalid/timeout URLs)
- Each worker performs an HTTP GET with a **4-second timeout** using `context`
- Results sent to a results channel, printed as they arrive
- Graceful handling: timeouts, unreachable hosts, invalid URLs
- Print final summary: total success vs failed
- Workers must stop cleanly when there are no more jobs

### 4. Given Variables & Structs
```go

var urls = []string{
	"https://google.com",
	"https://github.com",
	"https://go.dev",
	"https://pkg.go.dev",
	"https://cloudflare.com",
	"https://fastly.com",
	"https://stackoverflow.com",
	"https://reddit.com",
	"https://news.ycombinator.com",
	"https://gitlab.com",
	"https://bitbucket.org",
	"https://hub.docker.com",
	"https://kubernetes.io",
	"https://prometheus.io",
	// 6 URLs that will timeout
	"https://httpstat.us/200?sleep=10000",
	"https://httpstat.us/200?sleep=15000",
	"https://httpstat.us/200?sleep=20000",
	"https://10.255.255.1",             // non-routable IP, hangs
	"https://192.0.2.1",               // TEST-NET, hangs
	"https://198.51.100.1",            // TEST-NET-2, hangs
}

type Job struct {
    ID  int
    URL string
}

type Result struct {
	Job        Job
	WorkerID   int
	StatusCode int
	Duration   time.Duration
	Err        error
}
```

### 5. Expected Output
```
[worker 4] 1  ✅ https://github.com           → 200 (130ms)
[worker 2] 5  ✅ https://fastly.com           → 200 (191ms)
[worker 6] 3  ✅ https://pkg.go.dev           → 200 (209ms)
[worker 4] 6 ❌ https://stackoverflow.com    → 403 (80ms)
[worker 3] 2  ✅ https://go.dev               → 200 (226ms)
[worker 1] 0  ✅ https://google.com           → 200 (249ms)
[worker 5] 4  ✅ https://cloudflare.com       → 200 (364ms)
[worker 3] 10 ✅ https://bitbucket.org        → 200 (144ms)
[worker 1] 11 ✅ https://hub.docker.com       → 200 (189ms)
[worker 5] 12 ✅ https://kubernetes.io        → 200 (101ms)
[worker 3] 13 ✅ https://prometheus.io        → 200 (99ms)
[worker 2] 7  ✅ https://reddit.com           → 200 (332ms)
[worker 4] 9  ✅ https://gitlab.com           → 200 (467ms)
[worker 6] 8  ✅ https://news.ycombinator.com → 200 (666ms)
[worker 1] 14 ❌ https://10.255.255.1         → Get "https://10.255.255.1": context deadline exceeded (4001ms)
[worker 5] 15 ❌ https://10.255.255.2         → Get "https://10.255.255.2": context deadline exceeded (4001ms)
[worker 3] 16 ❌ https://10.255.255.3         → Get "https://10.255.255.3": context deadline exceeded (4000ms)
[worker 2] 17 ❌ https://10.255.255.1         → Get "https://10.255.255.1": context deadline exceeded (4001ms)
[worker 4] 18 ❌ https://192.0.2.1            → Get "https://192.0.2.1": context deadline exceeded (4001ms)
[worker 6] 19 ❌ https://198.51.100.1         → Get "https://198.51.100.1": context deadline exceeded (4001ms)

══════════════════════════════════════════════════
                      SUMMARY
══════════════════════════════════════════════════
 ✅  Healthy  (2xx)   :  13
 ❌  Unreachable      :  7
──────────────────────────────────────────────────
 Total                :  20
 Fastest              :  https://stackoverflow.com (80ms)
 Slowest (healthy)    :  https://news.ycombinator.com (666ms)
 Total runtime        :  4877ms
```

### 6. Why This Matters in Production
Worker pools are used everywhere:
- **Payment processors** — exactly N workers process transactions to avoid overloading downstream APIs
- **Web crawlers** — N goroutines fetch pages, respecting rate limits
- **Background job systems** — N workers drain a Redis/SQS queue
- **Health checkers** — this exact challenge, in production (Prometheus blackbox exporter does this)

Knowing when to use a worker pool vs fan-out is a key advanced Go topic.

### 7. Common Mistakes to Avoid
- Not closing the jobs channel — workers loop forever waiting for jobs that never come (goroutine leak)
- Closing the results channel from a worker goroutine — if multiple workers do this, panic
- Not using `context` for timeouts — HTTP requests can hang forever without it
- Using `time.Sleep` for timeouts instead of `context.WithTimeout` — never do this
- Do not close response body after a successful request

### 8. Hints & Knowledge
- `context.WithTimeout(context.Background(), 3*time.Second)` — cancels after 3s
- `http.NewRequestWithContext(ctx, "GET", url, nil)` — attaches context to request
- `close(jobs)` — workers reading `for job := range jobs` stop automatically
- Channel direction: `jobs <-chan Job` (receive-only), `results chan<- Result` (send-only)
- `time.Since(start)` — measure elapsed time

### 9. Sources
- Worker pools: https://gobyexample.com/worker-pools
- `context` package: https://pkg.go.dev/context
- `net/http` client: https://pkg.go.dev/net/http
- `select` statement: https://go.dev/tour/concurrency/5

### 10. Knowledge Gained
- ✅ Worker pool — fixed concurrency pattern
- ✅ `context` — timeout and cancellation (essential for all network code)
- ✅ `net/http` client — making HTTP requests in Go
- ✅ Channel directionality — enforcing send/receive contracts
- ✅ Graceful goroutine shutdown

---

# PHASE 4 — Memory & Performance

> **Why this phase matters**
> Go is used for high-performance systems because it gives you control over memory — without the danger of C. Companies running Go at scale (Cloudflare processes 50M+ req/s with Go) care deeply about allocations per request, GC pauses, and heap pressure. Knowing how to measure, profile, and reduce allocations is what separates surface-level Go from deep Go.

---

## Challenge 4 — The Benchmark Battle
### `🟠 Intermediate → Advanced`
**🕐 Expected duration: 15–20 hours**

### 1. Context
A data pipeline at your company processes dozen millions of log lines per day. Each line is parsed into a key-value map. The current implementation is correct but slow — and it's causing GC pressure because it allocates a new map on every single call. Your job: measure it, understand why it's slow, and fix it.

This is a real scenario. Datadog, Cloudflare, and similar companies do this kind of optimization routinely on their log ingestion pipelines.

### 2. Goal
Benchmark two existing implementations, analyze their memory behavior using Go's built-in tooling, then write a faster third version that wins on both time and allocations.

### 3. Scope
- Write proper benchmark tests for Version A and Version B (provided). For each round, the parser needs to parse the whole file.
- Run `go test -bench=. -benchmem` and record: `ns/op`, `B/op`, `allocs/op`
- Run `go build -gcflags="-m" . 2>&1 | grep "escapes to heap"` to see escape analysis output — what goes to heap?
- Write Version C using `sync.Pool` to reduce heap allocations
- Version C must have fewer `ns/op`, `B/ops` and `allocs/op` than both A and B
- Write a short explanation (as comments) of what each optimization does and why

### 4. Given Code

```go
// parser.go
func ParseA(line string) map[string]string {
	result := map[string]string{}
	parts := strings.Split(line, "|")
	for _, p := range parts {
		kv := strings.Split(strings.TrimSpace(p), "=")
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	return result
}
 
func ParseB(line string) map[string]string {
	pairs := strings.Split(line, "|")
	result := make(map[string]string, len(pairs))
	for _, p := range pairs {
		p = strings.TrimSpace(p)
		if idx := strings.Index(p, "="); idx != -1 {
			result[p[:idx]] = p[idx+1:]
		}
	}
	return result
}

// parser_test.go
// loadLines reads all lines from a file into memory.
// Use this in benchmarks so we only measure parsing time,
// not file I/O.
func loadLines(path string) []string {
    file, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines
}
```

#### File structure
```
your_folder/
├── parser.go
├── parser_test.go
└── testdata/
    └── logs.txt       ← provided (10,000 sample log lines)
```
### 5. Expected Output
```
go test -bench=. -benchmem
goos: darwin
goarch: arm64
pkg: benchmark
cpu: Apple M4
BenchmarkParseA-10           268           4329510 ns/op         8315234 B/op     116476 allocs/op
BenchmarkParseB-10           504           2427512 ns/op         5220761 B/op      33408 allocs/op
BenchmarkParseC-10           586           2037913 ns/op         1302817 B/op      10001 allocs/op
PASS
ok      benchmark       3.822s
```

### 6. Why This Matters in Production
At 10M requests/day, the difference between 10000 allocs/op and 30000 alloc/op is 20 million fewer heap allocations. Each allocation the GC doesn't have to track = less GC pause = lower tail latency. This is why high-performance Go services obsess over allocations per request.

### 7. Common Mistakes to Avoid
- Optimizing without measuring first — "premature optimization is the root of all evil"
- Not calling `mapPool.Put(result)` (or Release) after use — defeats the purpose of `sync.Pool`
- Not clearing the map after Get() — stale data leaks
- Confusing `b.N` — Go determines the right N automatically, never hardcode it
- Not running benchmarks multiple times — use `-count=5` for stable results

### 8. Hints & Knowledge
- `func BenchmarkX(b *testing.B) { for i := 0; i < b.N; i++ { ... } }` — standard shape
- `go test -bench=. -benchmem` — run all benchmarks with memory stats
- `B/op` = bytes allocated per operation, `allocs/op` = number of heap allocations
- Escape analysis: if a local variable is used after the function returns, it "escapes" to heap
- `sync.Pool`: Get() → use → Put() — reuse without allocation
- func `clear` helps clear your hash-map

#### Tools You'll Need (or want)

| Command | What it does |
|---|---|
| `go test -v` | Run all tests |
| `go test -run TestX -v` | Run one specific test |
| `go test -bench=. -benchmem` | Run all benchmarks with memory stats |
| `go test -bench=. -benchmem -count=5` | Run 5x for stable numbers |
| `go build -gcflags="-m" .` | Show escape analysis |
| `go test -cpuprofile cpu.out -bench=.` - optional | Generate CPU profile |
| `go tool pprof cpu.out` - optional | Explore with pprof |

### 9. Sources
- Go benchmarks: https://pkg.go.dev/testing#hdr-Benchmarks
- `sync.Pool`: https://pkg.go.dev/sync#Pool
- Escape analysis deep dive: https://go.dev/doc/faq#stack_or_heap
- pprof tutorial: https://go.dev/blog/pprof

### 10. Knowledge Gained
- ✅ Writing and interpreting Go benchmark tests
- ✅ `benchmem` — reading allocation output
- ✅ Escape analysis — understanding stack vs heap
- ✅ `sync.Pool` — object reuse pattern
- ✅ How to approach performance work: measure → profile → optimize → re-measure

---
# PHASE 5 — Standard Library & Systems Integration

> **Why this phase matters**
> Building real Go services means combining everything: goroutines for concurrency, interfaces for flexibility, HTTP for APIs, JSON for data, and mutexes for safe state. This phase simulates the architecture of a real microservice — the kind you'd find at any company running Go in production.

---

## Challenge 5.1 — Build a Mini DevOps Dashboard
### `🔴 Intermediate → Advanced`
**🕐 Expected duration: 25–30 hours**

### 1. Context
You're joining a platform team at a mid-size company. They have dozens of services writing log files to a shared directory. The ops team needs a lightweight internal tool that automatically picks up new log files, analyzes them for error rates, and exposes the results via a simple HTTP API — without restarting the service.

This is a simplified version of what tools like Fluentd, Logstash, and Vector do. You're building the Go-native version from scratch.

### 2. Goal
Build a self-contained HTTP service that watches a directory for log files, processes them concurrently, and exposes results through a REST JSON API.

### 3. Scope
The service has 3 components working together:

**Component A — File Watcher** (goroutine)
- Every 5 seconds, scan `./logs/` for new `.log` files
- Send new filenames to a jobs channel
- Track already-seen files — don't reprocess

**Component B — Worker Pool** (3 goroutines)
- Read filenames from jobs channel
- Count: total lines, ERROR, WARN, INFO per file
- Store results in a thread-safe store (`sync.Mutex`)

**Component C — HTTP API**
- `GET /status` → JSON of all processed files + counts
- `GET /errors` → JSON of only files with at least 1 error
- `POST /scan` → trigger an immediate re-scan without waiting for ticker

### 4. Expected Output
```bash
$ curl http://localhost:8080/status
{
  "total_files": 3,
  "files": [
    {"name": "log_1.log", "total": 100, "errors": 7, "warns": 12, "infos": 81}
  ]
}

$ curl http://localhost:8080/errors
{
  "files_with_errors": [
    {"name": "log_1.log", "errors": 7}
  ]
}
```

### 5. Why This Matters in Production
This challenge's architecture is the architecture of most Go microservices:
- A **background goroutine** doing periodic work (cron jobs, health checks, cache refresh)
- A **worker pool** processing a queue (job processors, event consumers)
- A **mutex-protected store** as the source of truth
- An **HTTP API** as the external interface

Understanding how these three components communicate safely is the difference between Go code that works and Go code that works *under load*.

### 6. Common Mistakes to Avoid
- Protecting reads *and* writes with the mutex — a common mistake is only locking writes
- Launching a new goroutine on every HTTP request instead of using the existing worker pool
- Not generating test log files — the service has nothing to do without them
- Using `map[string]FileStats` without a mutex — concurrent map writes cause runtime panic
- Blocking the HTTP handler while waiting for a scan to complete — use `go` and return `202 Accepted`

### 7. What a Senior Would Do Differently
- Replace `sync.Mutex` + `map` with a proper repository interface — easier to swap for Redis/Postgres later
- Use `chi` or `net/http`'s `ServeMux` patterns for cleaner routing
- Add structured logging with `log/slog` — every request logged with duration and status
- Add `/healthz` and `/readyz` endpoints — standard in any Kubernetes-deployed service
- Use `context` propagation from HTTP request through to file processing

### 8. Hints & Knowledge
- `http.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {})` — register a handler
- `json.NewEncoder(w).Encode(data)` — write JSON directly to response writer
- `w.Header().Set("Content-Type", "application/json")` — always set before writing body
- `time.NewTicker(5 * time.Second)` — fires every 5s, like a cron job
- `os.ReadDir("./logs/")` — returns `[]os.DirEntry`
- `sync.Mutex`: `mu.Lock()` before read/write, `defer mu.Unlock()` immediately after

### 9. Sources
- `net/http`: https://pkg.go.dev/net/http
- `encoding/json`: https://pkg.go.dev/encoding/json
- `sync.Mutex`: https://pkg.go.dev/sync#Mutex
- `time.Ticker`: https://pkg.go.dev/time#Ticker
- `os.ReadDir`: https://pkg.go.dev/os#ReadDir
- Go HTTP patterns: https://gobyexample.com/http-servers

### 10. Knowledge Gained
- ✅ `net/http` — building production HTTP servers
- ✅ JSON encoding for REST APIs
- ✅ `sync.Mutex` — protecting shared state under concurrent access
- ✅ `time.Ticker` — background periodic tasks
- ✅ Wiring goroutines + channels + HTTP into a cohesive service architecture
- ✅ The standard Go microservice architecture pattern

---

# PHASE 6 — Stress test

> **Why this phase matters**
> Technical skill and performance under pressure are different skills. Phase 6 trains the second one: reading code under pressure, explaining your decisions out loud, and writing correct Go quickly. After Phase 5 you have the knowledge — Phase 6 makes sure you can demonstrate it.

---

## Challenge 6.1 — Code Review Gauntlet
### `🔴 Advanced`
**🕐 Expected duration: 10 hours**

### 1. Context
Every developer will face hard time while developing. You're shown real-looking code with subtle bugs — goroutine leaks, race conditions, nil panics, interface misuse — and asked to spot and explain them. No running the code.

### 2. Goal
Review 5 broken Go programs. For each: identify all issues, explain why each is a problem, and write the fix.

### 3. Scope
*(Programs provided when you reach this challenge)*
Covers: goroutine leak, race condition, nil interface panic, bad interface design, performance anti-pattern.

### 4. Why This Matters in Production
Every Go team does code review. The ability to read a PR and say *"this goroutine leaks if the context is cancelled"* or *"this map access needs a mutex"* is what separates a beginner from an advanced engineer.

### 5. Common Mistakes to Avoid
- Assuming code is correct because it compiles — Go's concurrency bugs are runtime bugs
- Missing nil interface subtlety: an interface holding a nil pointer is NOT nil itself
- Not checking if channels are ever closed — the most common goroutine leak

### 6. What a Senior Would Do Differently
- Use `go vet` and `staticcheck` as automated first passes before manual review
- Reference https://100go.co — the canonical Go mistakes resource

### 7. Knowledge Gained
- ✅ Critical code reading skills
- ✅ Goroutine leak patterns
- ✅ Race condition identification
- ✅ Nil interface gotchas

---

## Challenge 6.2 — Build Under Pressure
### `🔴 Advanced`
**🕐 Expected duration: 8–10 hours**

### 1. Context
The final challenge. 3 timed problems, 45 minutes each, no hints. Designed to simulate real pressure — the kind you face during incidents, tight deadlines, or live debugging.

### 2. Goal
Solve 3 problems under time pressure. After each attempt: review together — what you got right, what you missed.

### 3. Scope
*(Problems provided when you reach this challenge)*
One problem each on: interfaces, concurrency, and performance.

### 4. Why This Matters
Performance under pressure is a skill. Writing correct Go under a 45-minute clock while explaining your thinking out loud is different from writing it comfortably at home. This challenge trains that specific skill.

### 5. Tips for Performing Well
- Read the problem twice — misunderstanding costs more time than reading slowly
- Define your data structures before writing logic
- Write the happy path first, then add error handling
- Name things clearly — `workerCount` not `wc`
- Say what you're doing out loud as you type - so that others can also understand

### 6. Knowledge Gained
- ✅ Performing under time pressure
- ✅ Structuring solutions quickly
- ✅ Clear Go style and technical communication

---

# Full Roadmap Summary

| Phase | Challenge | Hours | Level | Topics |
|---|---|---|---|---|
| 2 | Multi-Format Logger | 3–5h | 🟡 Beginner→Intermediate | interfaces, io.Writer, json |
| 2 | Shape Calculator Fix | 1–2h | 🟢 Beginner | type switch, interface bugs |
| 3 | Concurrent Log Scanner | 7–12h | 🟡 Beginner→Intermediate | goroutines, channels, mutex, WaitGroup |
| 3 | Worker Pool URL Checker | 8–12h | 🟠 Intermediate | worker pool, context, http client |
| 4 | Benchmark Battle | 6–10h | 🟠 Intermediate→Advanced | benchmarks, sync.Pool, pprof |
| 5 | Mini DevOps Dashboard | 14–20h | 🔴 Intermediate→Advanced | http, json API, mutex, ticker |
| 6 | Code Review Gauntlet | 4-6h | 🔴 Advanced | spot bugs, leaks, races |
| 6 | Build Under Pressure | 4–6h | 🔴 Advanced | pressure simulation |
| | **Total** | **~50–70h** | | |

---

*Start with Challenge 2.1. Come back when you're done — or when you're stuck.*
*This handbook is designed to be worked through in order. Don't skip phases.*