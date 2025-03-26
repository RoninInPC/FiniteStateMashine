package main

type EpsilonNFA[conditionType ConditionType, alphabetType AlphabetType] struct {
	NFA[conditionType, alphabetType]

	epsilon uint64
}

func nextEpsilon(body map[uint64]map[uint64][]uint64, current []uint64, epsilon uint64) []uint64 {
	conditionNext := make([]uint64, 0)

	for _, condition := range current {
		conditions, ok := body[condition][epsilon]

		if !ok {
			continue
		}

		conditionNext = append(conditionNext, nextEpsilon(body, conditions, epsilon)...)
		conditionNext = fixCurrent(conditionNext)
	}

	return conditionNext
}

func (N *EpsilonNFA[conditionType, alphabetType]) nextEpsilon() {
	N.current = nextEpsilon(N.body, N.current, N.epsilon)

}

func (N *EpsilonNFA[conditionType, alphabetType]) ToNext(alphabetType alphabetType) ([]conditionType, bool) {
	num, err := N.serializationAlphabet(alphabetType)
	if err != nil {
		return nil, false
	}
	N.nextEpsilon()

	conditionNext := make([]uint64, 0)
	for _, condition := range N.current {
		conditions, ok := N.body[condition][num]
		if !ok {
			continue
		}
		conditionNext = append(conditionNext, conditions...)
	}

	N.current = conditionNext

	return N.GetCurrentCondition(), N.IsFinal()
}
