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

func (N *EpsilonNFA[conditionType, alphabetType]) ToNext(alphabet alphabetType) ([]conditionType, bool) {
	num, err := N.serializationAlphabet(alphabet)
	if err != nil {
		return nil, false
	}
	N.nextEpsilon()

	c := len(N.current)
	conditionNext := make([]uint64, 0)
	for _, condition := range N.current {
		conditions, ok := N.body[condition][num]
		if !ok {
			c--
			continue
		}
		conditionNext = append(conditionNext, conditions...)
	}
	if c == 0 {
		return nil, false
	}
	N.nextEpsilon()

	N.current = conditionNext

	return N.GetCurrentCondition(), N.IsFinal()
}

func (N *EpsilonNFA[conditionType, alphabetType]) GetEpsilon() alphabetType {
	return *N.alphabet[N.epsilon]
}

func DefaultEpsilonNFA[conditionType ConditionType, alphabetType AlphabetType](epsilon alphabetType) MakerEpsilonNFA[conditionType, alphabetType] {
	return MakerEpsilonNFA[conditionType, alphabetType]{
		DefaultNFA[conditionType, alphabetType](),
		epsilon,
	}
}

func SerializationEpsilonNFA[conditionType ConditionType, alphabetType AlphabetType](
	serializationCondition SerializationCondition[conditionType],
	serializationAlphabet SerializationAlphabet[alphabetType],
	epsilon alphabetType) MakerEpsilonNFA[conditionType, alphabetType] {
	return MakerEpsilonNFA[conditionType, alphabetType]{
		SerializationNFA[conditionType, alphabetType](serializationCondition, serializationAlphabet),
		epsilon,
	}
}

type MakerEpsilonNFA[conditionType ConditionType, alphabetType AlphabetType] struct {
	MakerNFA[conditionType, alphabetType]
	epsilon alphabetType
}

func (m MakerEpsilonNFA[conditionType, alphabetType]) AddTransitions(start conditionType, transitionAlpha alphabetType, end []conditionType) MakerEpsilonNFA[conditionType, alphabetType] {
	m.MakerNFA.AddTransitions(start, transitionAlpha, end)
	return m
}

func (m MakerEpsilonNFA[conditionType, alphabetType]) AddTransition(start conditionType, transitionAlpha alphabetType, end conditionType) MakerEpsilonNFA[conditionType, alphabetType] {
	m.MakerNFA.AddTransition(start, transitionAlpha, end)
	return m
}

func (m MakerEpsilonNFA[conditionType, alphabetType]) SetStart(condition conditionType) MakerEpsilonNFA[conditionType, alphabetType] {
	m.MakerNFA.SetStart(condition)
	return m
}

func (m MakerEpsilonNFA[conditionType, alphabetType]) SetEnd(conditions []conditionType) MakerEpsilonNFA[conditionType, alphabetType] {
	m.MakerNFA.SetEnd(conditions)
	return m
}

func (m MakerEpsilonNFA[conditionType, alphabetType]) Build() EpsilonNFA[conditionType, alphabetType] {
	epsilonUint, _ := m.build.serializationAlphabet(m.epsilon)
	m.build.alphabet[epsilonUint] = &m.epsilon
	return EpsilonNFA[conditionType, alphabetType]{
		m.build,
		epsilonUint,
	}
}
