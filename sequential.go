package analyzer

// Sequential is the supplied correctness baseline. It does not call MapTask,
// Reduce, or Coordinator, so it can detect their mistakes.
func Sequential(lines []string) (map[string]int, int) {
	counts := make(map[string]int)
	rejected := 0
	for _, line := range lines {
		path, class, ok := Parse(line)
		if !ok {
			rejected++
			continue
		}
		counts["path:"+path]++
		counts["status:"+class]++
	}
	return counts, rejected
}
