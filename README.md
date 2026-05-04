# Architecture & Scalability Upfluence technical test

This project is my submitted solution to the Upfluence technical test, written in Go, reading from an SSE stream of data to compute some metrics. 

This was done over the course of the last two weeks, _weekends excluded_.

## Architecture

This project follows the **hexagonal architecture** pattern to maintain a strict separation between business logic and infrastructure.

### Layers

- **Services**: atomic business logic, divided in two components; the _stream_ service manages data flow, and the _compute_ service handles mathematics.
- **UseCases**: contains the core business rules that will use the services. This is where I specify what percentiles to compute, fetch the timestamps, and store the analyzed metric. It is independent of frameworks.
- **Router**: uses Echo for HTTP routing and validation. I chose Echo because it's a modern library, still maintained to this day, and provides simple and efficient ways to ensure the data we send is correct (through the use of *validators*).
- **Platforms**: external repository adapters. Implements the interface that will fetch data from the Upfluence SSE API.

### Production Trade-offs

1. <u>Complexity of the code</u>

The main downside of the hexagonal architecture is that it complexifies the code a lot. This simple application required a lot of groundwork, even if it will make debugging and errors handling much faster afterwards.

2. <u>Memory footprint</u>

Currently, we store all valid metrics in a slice to compute exact percentiles. For very long durations or huge throughput, this could lead to big memory usage. This is where a repo like _Redis_ could be useful.

3. <u>Code efficiency</u>

I used a sorting method on the metrics array, inside my service, so when I call the method ComputePercentile three times in a row in my usecase, it will check for the array being sorted three times.


## Testing Strategy

The project aims for high reliability with a pragmatic testing approach, with both unit and integration tests. It lacks end to end tests, because it would have taken a lot of time to make efficient e2e tests. I aimed to at least reach 80% coverage on every tested packages.

## Installation & Running

### Prerequisites

- Go: 1.24.3+
- Make (optional)

### Download dependencies (with make):

make deps


### Run the server (with make):

make run

### Run the server (without make):

go run ./cmd/server/main.go

---

The server will listen by default on port 8080 of the localhost, but you can configure whatever port you want in the provided .env file, under the "PORT" field.


## Functional Notes

- **Accuracy**: percentiles are calculated using linear interpolation.
- **Validation**: invalid durations or dimensions return a 400 Bad Request with detailed JSON errors (RFC 7807).
- **Graceful Shutdown**: the server listens for SIGINT/SIGTERM to close active stream connections properly.

## Technical Stack

- **Framework**: Echo v4 (chosen for performance and middleware support).
- **Logging**: slog (structured JSON logging).
- **Standard Library**: used for everything else (math, strings, time, buffering).

---
Thank you for your time.