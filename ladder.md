# Backend Engineering Ladder

Learn by building **one evolving project** (e.g. a link shortener or a tiny ticket-booking app).
Build each layer yourself before a framework hides it. Each stage should create a problem the next stage solves.

Progress is measured by **checkpoints, not time**. Move on once you can do everything under "Done when". Don't move on before that.

## Setup

- **Language:** Go. Its standard library covers HTTP, SQL and concurrency with nothing hidden, concurrency is built into the language (goroutines, channels), and it means getting comfortable in a language other than Python.
- **Layout:** one folder per stage while it's a throwaway (`raw-http/`). From stage 2 on, one project that keeps growing.
- **Tools:** terminal, git, `curl -v`.

---

## 1. Raw HTTP server
Build a TCP echo server, then an HTTP server on raw sockets. Parse the request line and headers yourself and return a valid response.

**Learn:** HTTP is text over TCP. Ports, request/response, status codes, headers.

**Done when:**
- [ ] A browser and `curl` both get a correct response from your socket server
- [ ] You can explain each line of a raw HTTP request and response
- [ ] Malformed requests get a `400` instead of crashing the server

## 2. JSON CRUD API (in-memory)
Standard-library HTTP server with create/read/update/delete over JSON.

**Learn:** routing, REST conventions, input validation, error handling.

**Done when:**
- [ ] All CRUD endpoints return the correct status codes (200/201/204/400/404/405)
- [ ] Invalid input returns structured error JSON

## 3. Persistence: file → SQLite → Postgres
**Learn:** why databases exist, SQL, schema design, migrations, indexes.

**Done when:**
- [ ] Data survives restarts
- [ ] You can explain what breaks with plain file storage (concurrent writes, partial writes, querying)
- [ ] Schema changes go through migrations
- [ ] You added an index and confirmed it with `EXPLAIN ANALYZE`

## 4. Users and auth
Sign-up/login with hashed passwords (bcrypt/argon2), sessions, then tokens. Users own their own resources.

**Learn:** authentication vs authorization, cookies, why you never store plain passwords.

**Done when:**
- [ ] User A can't read or modify User B's data (checked with a test)
- [ ] You can explain the trade-offs between sessions and tokens

## 5. Tests
Unit tests plus integration tests against a real database.

**Learn:** writing testable code, test fixtures, test isolation.

**Done when:**
- [ ] One command runs the whole suite against a fresh DB
- [ ] Every bug fixed from here on gets a regression test first

## 6. Concurrency bugs on purpose
Let two requests book the last ticket at the same moment, reproduce the double booking, then fix it.

**Learn:** race conditions, transactions, isolation levels, locks, idempotency keys.

**Done when:**
- [ ] A concurrent test reliably reproduces the bug before the fix and passes after it
- [ ] Retrying the same request with an idempotency key doesn't create duplicates

## 7. Rate limiter and cache
Token-bucket rate limiter, then a Redis cache on a hot read path.

**Learn:** caching trade-offs, cache invalidation, protecting a service from abuse.

**Done when:**
- [ ] Clients over the limit get `429` with a `Retry-After` header
- [ ] Cached reads are measurably faster, and updates never serve stale data longer than you intend

## 8. Background jobs
A worker that sends "emails" from a queue, with retries and backoff.

**Learn:** async processing, at-least-once delivery, failure handling, dead-letter queues.

**Done when:**
- [ ] Killing the worker mid-job loses no jobs
- [ ] A job that keeps failing ends up in a dead-letter queue instead of retrying forever

## 9. Deploy
Docker, a cheap VPS, nginx as reverse proxy, TLS, config through environment variables.

**Learn:** processes, reverse proxies, TLS, secrets management.

**Done when:**
- [ ] The app is reachable over HTTPS on a real domain
- [ ] There are no secrets in the repo
- [ ] Redeploying is a single command

## 10. Observability and load testing
Structured logs and metrics, then a load test (k6 or `hey`). Find the bottleneck, fix it, measure again.

**Learn:** observability, and performance work driven by numbers rather than guesses.

**Done when:**
- [ ] You have before/after numbers (p50/p99 latency, throughput) for one real optimization
- [ ] You can trace a single slow request through the logs

---

## Stretch: deeper fundamentals
- [ ] Your own key-value store (write-ahead log → in-memory index → compaction)
- [ ] A mini Redis that speaks the RESP protocol
- [ ] A simple message queue

## Habits
- **Before adding a tool, write down what problem it solves.** If you can't say, you aren't ready for it yet.
- **Break things on purpose:** kill the DB mid-request, send malformed JSON, fire 1,000 concurrent requests.
- **Read the raw traffic:** `curl -v`, browser dev tools, `EXPLAIN ANALYZE`.
- **After stage 6, read** *Designing Data-Intensive Applications*.
