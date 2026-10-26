package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"sort"

	a "cs411-lab09-mapreduce"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: analyzer data/access.log")
		os.Exit(2)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	c := a.NewCoordinator()
	for start := 0; start < len(lines); start += 3 {
		end := start + 3
		if end > len(lines) {
			end = len(lines)
		}
		pairs, rejected := a.MapTask(lines[start:end])
		id := fmt.Sprintf("map-%d", start/3)
		c.Commit(id, pairs, rejected)
		if start == 0 {
			c.Commit(id, pairs, rejected)
		} // injected retry
	}
	result, loads := c.Output(3)
	baseline, baselineRejected := a.Sequential(lines)
	if !reflect.DeepEqual(result, baseline) || c.Rejected != baselineRejected {
		fmt.Fprintln(os.Stderr, "distributed result differs from sequential reference")
		os.Exit(1)
	}
	keys := make([]string, 0, len(result))
	for key := range result {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("%s %d\n", key, result[key])
	}
	fmt.Printf("rejected %d\nreducer_loads %v\n", c.Rejected, loads)
}
