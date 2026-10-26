package analyzer

import (
	"reflect"
	"testing"
)

var sample = []string{"GET /search 200", "POST /cart 201", "GET /search 503", "bad record", "GET /home 404"}

func TestMapRecordTrace(t *testing.T) {
	pairs, rejected := MapTask(sample[:2])
	want := []Pair{{"path:/search", 1}, {"status:2xx", 1}, {"path:/cart", 1}, {"status:2xx", 1}}
	if rejected != 0 || !reflect.DeepEqual(pairs, want) {
		t.Fatalf("pairs=%v rejected=%d", pairs, rejected)
	}
}

func TestReduceAndRetry(t *testing.T) {
	c := NewCoordinator()
	pairs1, rejected1 := MapTask(sample[:3])
	pairs2, rejected2 := MapTask(sample[3:])
	c.Commit("map-1", pairs1, rejected1)
	c.Commit("map-1", pairs1, rejected1) // retry must be ignored
	c.Commit("map-2", pairs2, rejected2)
	got, loads := c.Output(3)
	want, wantRejected := Sequential(sample)
	if !reflect.DeepEqual(got, want) || c.Rejected != wantRejected {
		t.Fatalf("output=%v rejected=%d; want=%v rejected=%d", got, c.Rejected, want, wantRejected)
	}
	total := 0
	for _, load := range loads {
		total += load
	}
	if total != 2*(len(sample)-wantRejected) {
		t.Fatalf("reducer loads sum to %d", total)
	}
}
