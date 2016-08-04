package strategy

// ExtendedStrategy takes an error as the second argument to analyze it.
type ExtendedStrategy func(attempt uint, err error) bool

// ExtendStrategies takes list of Strategy and returns list of ExtendedStrategy.
func ExtendStrategies(strategies ...Strategy) []ExtendedStrategy {
	if len(strategies) == 0 {
		return nil
	}

	extendedStrategies := make([]ExtendedStrategy, 0, len(strategies))
	for i := range strategies {
		extendedStrategies = append(extendedStrategies, (func(strategy Strategy) func(uint, error) bool {
			return func(attempt uint, err error) bool {
				return strategy(attempt)
			}
		})(strategies[i]))
	}

	return extendedStrategies
}
