package strategy

import "testing"

func TestExtendStrategies(t *testing.T) {
	max := uint(2)
	strategies := []Strategy{
		Limit(max),
	}
	extendedStrategy := ExtendStrategies(strategies...)

	for i := uint(0); i < max+2; i++ {
		result1 := strategies[0](i)
		result2 := extendedStrategy[0](i, nil)
		if result1 != result2 {
			t.Errorf("expected the same result: expected %t, got %t", result1, result2)
		}
	}
}

func TestExtendStrategiesWithEmptyList(t *testing.T) {
	var strategies []ExtendedStrategy

	strategies = ExtendStrategies(nil...)
	if len(strategies) != 0 {
		t.Errorf("expected empty list of strategies got count %d", len(strategies))
	}

	strategies = ExtendStrategies([]Strategy{}...)
	if len(strategies) != 0 {
		t.Errorf("expected empty list of strategies got count %d", len(strategies))
	}
}
