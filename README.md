# Upfluence Technical Test

This project is my submitted solution to the Upfluence technical test. This was done over the course of the last two weeks, _weekends excluded_.

## Architecture

This project is written in **Go** and follows the **hexagonal architecture** pattern to maintain a strict separation between business logic and infrastructure. I choose Go over other languages because I wanted to use the same tools Upfluence uses, to better align with their technical stack. The hexagonal architecture pattern is one of the most utilized pattern in the industry, proving great isolation of services and repository adaptability though ports and adapters components. The only downside of this architecture is the time needed to setup even a simple project as this one, making the code very complex although making debugging and error handling smoothers and simpler in the future.

### Layers

- **Services**: atomic business logic, divided in two components; the _stream_ service manages data flow, and the _compute_ service handles mathematics.
- **UseCases**: contains the core business rules that will use the services. This is where I specify what percentiles to compute, fetch the timestamps, and store the analyzed metric. It is independent of frameworks.
- **Router**: uses Echo for HTTP routing and validation. I chose Echo because it's a modern library, still maintained to this day, and provides simple and efficient ways to ensure the data we send is correct (through the use of *validators*).
- **Platforms**: external repository adapters. Implements the interface that will fetch data from the Upfluence SSE API.

### Routes

* /analysis

The route takes 2 **mandatory** query parameters:

- duration, which is the total duration the app will listen to the stream data flux. It's formatted as an int followed by either the character 's', 'm', or 'h'.
- dimension, which is a string that represents the metric we'll analyze in the retrieve posts. It can be equal to either "likes", "comments', "favorites" or "retweets".

Every other query parameter will return an error 400 (Bad Request).

### Production Trade-offs

1. <u>Memory footprint</u>

Currently, we store all valid metrics in a slice to compute exact percentiles. For very long durations or huge throughput, this could lead to big memory usage. This is where a repo like _Redis_ could be useful.
I also didn't add HTTP connection pools, as an output of my repository. It would've been great to handle big amount of users using my API.

2. <u>Code efficiency</u>

I used a sorting method on the metrics array, inside my service, so when I call the method ComputePercentile three times in a row in my usecase, it will check for the array being sorted three times.

3. <u>Lack of End-to-End tests</u>

In my testing startegy, I tried to aim at 80% minimum coverage, through the use of unit and integration tests. However, I chose to not implement E2E tests, as it would have taken much more time, trying to pinpoint every example usage and making efficient and interesting tests.

4. <u>Lack of tools</u>

In a real project, shipped for production, I would've add more tools, such as documentation with Swagger, monitoring through apps like Prometeus to gather metrics about the code efficiency and memory usage, Docker containerization to simplify application launch, and a real validating system for the user inputs (the one currently implemented is basic and not very robust). However, this would've taken much more time, so I chose not to implement them.


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

### Launching the app

curl "http://localhost:8080/analysis?duration=5m&dimension=likes"

-> will output, after 5 minutes, a JSON payload containing the total number of posts analyzed, the timestamp range of the analyzed posts, and the 50th, 90th and 99th percentile of the dimension queried.

---

The server will listen by default on port 8080 of the localhost, but you can configure whatever port you want in the provided .env file, versioned for simplicity, under the "PORT" field.


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