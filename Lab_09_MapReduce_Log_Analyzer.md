# Lab 9: MapReduce Log Analyzer

**CS 411: Computing Paradigms - Spring 2027**  
**Week 11: March 22-28, 2027**  
**Due:** March 28, 2027  
**Value:** 100 points

**Starter repository:** [Open the Spring 2027 starter](https://github.com/cwillpinto/cs411-spring2027-lab09-mapreduce-log-analyzer)

## Purpose

Build a MapReduce pipeline that counts request paths and status classes from web logs. Verify the distributed result against a sequential reference and prove that retrying a task does not double-count output.

## Learning objectives

- Trace Map Shuffle Reduce
- Choose stable intermediate keys
- Make retried tasks deterministic and idempotent

## Before you begin

1. Review the Week 11 slide deck and assigned VoiceThread or reading.
2. Create your own repository from the Spring 2027 template linked in Blackboard. It includes `data/access.log`, an independent sequential baseline, public tests, and TODOs in `map.go`, `reduce.go`, and `coordinator.go`.
3. Run `go test ./...` before editing and save the initial output in `Lab9_Terminal.txt`. The public trace and retry checks fail at the intended TODOs.

## Part A Record trace (25 points)

Follow one log line through every phase.

- In `Lab9_Record_Trace.md` or PDF, trace `GET /search 200` through parsing, map, partition, shuffle, and reduce. List both emitted pairs: `path:/search` and `status:2xx`, each with count 1.
- For three reducers, calculate each pair's partition with the supplied `Partition` function and show the grouped reducer input.
- State the final counts for those keys in the supplied data file. Explain why `bad record` is rejected rather than mapped.

## Part B Pipeline (45 points)

Implement deterministic map and reduce tasks.

- Complete `MapTask` to emit two stable keys per valid `METHOD PATH STATUS` record and count rejected records. Preserve input order in the emitted pairs.
- Complete `Reduce` to sum counts by key and `Coordinator.Commit` to accept a task ID only once. The supplied coordinator partitions intermediate keys; explain why every occurrence of a key reaches the same reducer.
- Run `go run ./cmd/analyzer data/access.log`. Its CLI compares output to the independent sequential baseline and prints final keys in sorted order.

## Part C Retry and skew (30 points)

Test operational behavior.

- The CLI submits the first map task twice with the same task ID. Verify that the committed output and rejected count appear once. Run `go test ./...` and `go test -race ./...` and save commands/output in `Lab9_Terminal.txt`.
- Save the CLI's sorted counts, rejected count, and `reducer_loads` line as `Lab9_Output.txt`.
- Calculate the largest reducer's share of emitted pair counts: `max(reducer_loads) / sum(reducer_loads)`. Explain any skew using the path and status keys, and name one way to reduce a hot-key bottleneck.

## Verification

- `go test ./...` and `go test -race ./...` pass.
- The CLI exits successfully, and sorted distributed output equals the sequential reference.
- A retried task does not change totals or rejected-record count.

## What to submit

Submit your own repository URL and:

1. Your own repository URL (created from the Spring 2027 template).
2. `map.go`, `reduce.go`, and `coordinator.go` in that repository.
3. `Lab9_Record_Trace.pdf` or `Lab9_Record_Trace.md`, including the skew calculation and explanation.
4. `Lab9_Output.txt` and `Lab9_Terminal.txt`.

## Rubric

| Criterion | Points |
| --- | ---: |
| Record trace | 20 |
| Map and shuffle design | 25 |
| Reduce correctness | 20 |
| Retry safety | 20 |
| Reference comparison and skew analysis | 15 |
| **Total** | **100** |

## AI usage and academic integrity

Add an **AI Usage** section to the required trace file. Name each tool and list the prompts or interactions that affected the work. If you did not use AI, write `AI Usage: None.` You must be able to explain and defend every submitted design and line of code.

## Late work

Late work may be accepted for up to three calendar days with a 25 percent deduction per day, following the course syllabus.

