package analyzer

import (
	"strconv"
	"strings"
)

type Pair struct {
	Key   string
	Count int
}

// Parse accepts METHOD PATH STATUS (e.g. GET /search 200).
func Parse(line string) (path, class string, ok bool) {
	parts := strings.Fields(line)
	if len(parts) != 3 || (parts[0] != "GET" && parts[0] != "POST") || !strings.HasPrefix(parts[1], "/") {
		return "", "", false
	}
	status, err := strconv.Atoi(parts[2])
	if err != nil || status < 100 || status > 599 {
		return "", "", false
	}
	return parts[1], strconv.Itoa(status/100) + "xx", true
}

// TODO: emit path:<path> and status:<class> with count 1 per valid record;
// count invalid records separately. Pair ordering should follow input order.
func MapTask(lines []string) (pairs []Pair, rejected int) {
	_ = lines
	return nil, 0
}
