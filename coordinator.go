package analyzer

import "hash/fnv"

type Coordinator struct {
	committed map[string][]Pair
	Rejected  int
}

func NewCoordinator() *Coordinator { return &Coordinator{committed: make(map[string][]Pair)} }

// TODO: commit a task ID at most once. Retried output with the same ID must
// not add pairs or rejected records a second time. Copy the accepted pairs.
func (c *Coordinator) Commit(taskID string, pairs []Pair, rejected int) {
	c.committed[taskID] = append(c.committed[taskID], pairs...)
	c.Rejected += rejected
}

func Partition(key string, reducers int) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int(h.Sum32() % uint32(reducers))
}

func (c *Coordinator) Output(reducers int) (map[string]int, []int) {
	if reducers < 1 {
		reducers = 1
	}
	buckets := make([][]Pair, reducers)
	loads := make([]int, reducers)
	for _, pairs := range c.committed {
		for _, pair := range pairs {
			index := Partition(pair.Key, reducers)
			buckets[index] = append(buckets[index], pair)
			loads[index] += pair.Count
		}
	}
	result := make(map[string]int)
	for _, bucket := range buckets {
		for key, count := range Reduce(bucket) {
			result[key] += count
		}
	}
	return result, loads
}
