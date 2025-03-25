package main

import "slices"

// NondeterministicFiniteAutomaton или Недетерменированный конечный автомат

type NFA[conditionType ConditionType, alphabetType AlphabetType] struct {
	body map[uint64]map[uint64][]uint64

	serializationCondition SerializationCondition[conditionType]
	serializationAlphabet  SerializationAlphabet[alphabetType]

	conditions map[uint64]*conditionType
	alphabet   map[uint64]*alphabetType

	current []uint64
	start   uint64
	end     []uint64
}

func (N *NFA[conditionType, alphabetType]) Correct(alpha ...alphabetType) bool {
	isFinal := false

	for _, partAlpha := range alpha {
		_, isFinal = N.ToNext(partAlpha)
	}
	N.Reset()

	return isFinal
}

func (N *NFA[conditionType, alphabetType]) GetCurrentCondition() []conditionType {
	currentCondition := make([]conditionType, len(N.current))

	for i, condition := range N.current {
		currentCondition[i] = *N.conditions[condition]
	}

	return currentCondition
}

func (N *NFA[conditionType, alphabetType]) GetAlphabet() []alphabetType {
	alphabet := make([]alphabetType, len(N.alphabet))

	i := 0
	for _, alpha := range N.alphabet {
		alphabet[i] = *alpha
		i++
	}

	return alphabet
}

func (N *NFA[conditionType, alphabetType]) ToNext(alphabetType alphabetType) ([]conditionType, bool) {
	num, err := N.serializationAlphabet(alphabetType)
	if err != nil {
		return nil, false
	}

	conditionNext := make([]uint64, 0)
	for _, condition := range N.current {
		conditionNext = append(conditionNext, N.body[condition][num]...)
	}

	N.current = conditionNext

	return N.GetCurrentCondition(), N.IsFinal()
}

func (N *NFA[conditionType, alphabetType]) IsFinal() bool {
	if len(N.end) == 0 {
		return false
	}
	for _, condition := range N.current {
		if slices.Contains(N.end, condition) {
			return true
		}
	}
	return false
}

func (N *NFA[conditionType, alphabetType]) Reset() {
	N.current = []uint64{N.start}
}

type MakerNFA[conditionType ConditionType, alphabetType AlphabetType] struct {
	build NFA[conditionType, alphabetType]
}

func DefaultNFA[conditionType ConditionType, alphabetType AlphabetType]() MakerNFA[conditionType, alphabetType] {
	return MakerNFA[conditionType, alphabetType]{
		build: NFA[conditionType, alphabetType]{
			serializationCondition: SerializationConditionDefault[conditionType],
			serializationAlphabet:  SerializationAlphabetDefault[alphabetType],
			body:                   make(map[uint64]map[uint64][]uint64),
			alphabet:               make(map[uint64]*alphabetType),
			conditions:             make(map[uint64]*conditionType),
			current:                make([]uint64, 0),
			end:                    make([]uint64, 0),
		},
	}
}

func SerializationNFA[conditionType ConditionType, alphabetType AlphabetType](
	serializationCondition SerializationCondition[conditionType],
	serializationAlphabet SerializationAlphabet[alphabetType]) MakerNFA[conditionType, alphabetType] {
	return MakerNFA[conditionType, alphabetType]{
		build: NFA[conditionType, alphabetType]{
			serializationCondition: serializationCondition,
			serializationAlphabet:  serializationAlphabet,
			body:                   make(map[uint64]map[uint64][]uint64),
			alphabet:               make(map[uint64]*alphabetType),
			conditions:             make(map[uint64]*conditionType),
			current:                make([]uint64, 0),
			end:                    make([]uint64, 0),
		},
	}
}

func (m MakerNFA[conditionType, alphabetType]) Build() NFA[conditionType, alphabetType] {
	return m.build
}

func (m MakerNFA[conditionType, alphabetType]) AddTransition(start conditionType, transitionAlpha alphabetType, end []conditionType) MakerNFA[conditionType, alphabetType] {
	startUint, _ := m.build.serializationCondition(start)
	transition, _ := m.build.serializationAlphabet(transitionAlpha)
	endUint := make([]uint64, len(end))

	m.build.conditions[startUint] = &start

	m.build.alphabet[transition] = &transitionAlpha

	for i, condition := range end {
		endUint[i], _ = m.build.serializationCondition(condition)
		m.build.conditions[endUint[i]] = &condition
	}
	if _, ok := m.build.body[startUint]; !ok {
		m.build.body[startUint] = make(map[uint64][]uint64)
	}

	m.build.body[startUint][transition] = endUint

	return m
}

func (m MakerNFA[conditionType, alphabetType]) SetStart(condition conditionType) MakerNFA[conditionType, alphabetType] {
	conditionUint, _ := m.build.serializationCondition(condition)

	m.build.conditions[conditionUint] = &condition
	m.build.start = conditionUint
	m.build.current = []uint64{conditionUint}
	return m
}

func (m MakerNFA[conditionType, alphabetType]) SetEnd(conditions []conditionType) MakerNFA[conditionType, alphabetType] {
	end := make([]uint64, len(conditions))
	for i, condition := range conditions {
		conditionUint, _ := m.build.serializationCondition(condition)
		end[i] = conditionUint
		m.build.conditions[conditionUint] = &condition
	}
	m.build.end = end
	return m
}
