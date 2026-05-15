package dag

import "testing"

func TestExecutionOrderGenerator(t *testing.T) {
	graph := map[string][]string{
		"A": {"B", "C"},
		"B": {"D", "E"},
		"C": {"F"},
		"D": {"G"},
		"E": {"G"},
		"F": {"H"},
		"G": {},
		"H": {},
	}

	got,err := ExecutionOrderGenerator(graph)
	if err != nil {
		t.Error(got)
	}
}

func BenchmarkExecutionOrderGenerator(b *testing.B) {
	graph := map[string][]string{
		"A": {"B", "C"},
		"B": {"D", "E"},
		"C": {"F"},
		"D": {"G"},
		"E": {"G"},
		"F": {"H"},
		"G": {},
		"H": {},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ExecutionOrderGenerator(graph)
	}
}

func BenchmarkExecutionOrderGeneratorWithCycle(b *testing.B) {
	graph := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": {"D"},
		"D": {"A"}, // cycle
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ExecutionOrderGenerator(graph)
	}
}