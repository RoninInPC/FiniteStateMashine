package main

import "slices"

type DFA[conditionType ConditionType, alphabetType AlphabetType] struct {
	body map[uint64]map[uint64]uint64

	serializationCondition SerializationCondition[conditionType]
	serializationAlphabet  SerializationAlphabet[alphabetType]

	conditions map[uint64]*conditionType
	alphabet   map[uint64]*alphabetType

	current uint64
	start   uint64
	end     []uint64
}

func (D *DFA[conditionType, alphabetType]) Correct(alpha ...alphabetType) bool {
	isFinal := false

	for _, partAlpha := range alpha {
		_, isFinal = D.ToNext(partAlpha)
	}
	D.Reset()

	return isFinal
}

func (D *DFA[conditionType, alphabetType]) GetCurrentCondition() []conditionType {
	return []conditionType{*D.conditions[D.current]}
}

func (D *DFA[conditionType, alphabetType]) GetDFACurrentCondition() conditionType {
	return *D.conditions[D.current]
}

func (D *DFA[conditionType, alphabetType]) GetAlphabet() []alphabetType {
	alphabet := make([]alphabetType, len(D.alphabet))
	i := 0
	for _, alpha := range D.alphabet {
		alphabet[i] = *alpha
		i++
	}
	return alphabet
}

func (D *DFA[conditionType, alphabetType]) ToNext(alphabet alphabetType) ([]conditionType, bool) {
	alphabetTypeUint, _ := D.serializationAlphabet(alphabet)
	newCurrent, ok := D.body[D.current][alphabetTypeUint]
	if ok {
		D.current = newCurrent
	} else {
		return nil, false
	}
	return D.GetCurrentCondition(), D.IsFinal()
}

func (D *DFA[conditionType, alphabetType]) IsFinal() bool {
	return slices.Contains(D.end, D.current)
}

func (D *DFA[conditionType, alphabetType]) Reset() {
	D.current = D.start
}

type MakerDFA[conditionType ConditionType, alphabetType AlphabetType] struct {
	build DFA[conditionType, alphabetType]
}

func DefaultDFA[conditionType ConditionType, alphabetType AlphabetType]() MakerDFA[conditionType, alphabetType] {
	return MakerDFA[conditionType, alphabetType]{
		build: DFA[conditionType, alphabetType]{
			serializationCondition: SerializationConditionDefault[conditionType],
			serializationAlphabet:  SerializationAlphabetDefault[alphabetType],
			body:                   make(map[uint64]map[uint64]uint64),
			alphabet:               make(map[uint64]*alphabetType),
			conditions:             make(map[uint64]*conditionType),
			end:                    make([]uint64, 0),
		},
	}
}

func SerializationDFA[conditionType ConditionType, alphabetType AlphabetType](
	serializationCondition SerializationCondition[conditionType],
	serializationAlphabet SerializationAlphabet[alphabetType]) MakerDFA[conditionType, alphabetType] {
	return MakerDFA[conditionType, alphabetType]{
		build: DFA[conditionType, alphabetType]{
			serializationCondition: serializationCondition,
			serializationAlphabet:  serializationAlphabet,
			body:                   make(map[uint64]map[uint64]uint64),
			alphabet:               make(map[uint64]*alphabetType),
			conditions:             make(map[uint64]*conditionType),
			end:                    make([]uint64, 0),
		},
	}
}

func (m MakerDFA[conditionType, alphabetType]) Build() DFA[conditionType, alphabetType] {
	return m.build
}

func (m MakerDFA[conditionType, alphabetType]) AddTransition(start conditionType, transitionAlpha alphabetType, end conditionType) MakerDFA[conditionType, alphabetType] {
	startUint, _ := m.build.serializationCondition(start)
	transition, _ := m.build.serializationAlphabet(transitionAlpha)

	m.build.conditions[startUint] = &start

	m.build.alphabet[transition] = &transitionAlpha

	if _, ok := m.build.body[startUint]; !ok {
		m.build.body[startUint] = make(map[uint64]uint64)
	}

	m.build.body[startUint][transition], _ = m.build.serializationCondition(end)

	return m
}

func (m MakerDFA[conditionType, alphabetType]) SetStart(condition conditionType) MakerDFA[conditionType, alphabetType] {
	conditionUint, _ := m.build.serializationCondition(condition)

	m.build.conditions[conditionUint] = &condition
	m.build.start = conditionUint
	m.build.current = conditionUint
	return m
}

func (m MakerDFA[conditionType, alphabetType]) SetEnd(conditions []conditionType) MakerDFA[conditionType, alphabetType] {
	end := make([]uint64, len(conditions))
	for i, condition := range conditions {
		conditionUint, _ := m.build.serializationCondition(condition)
		end[i] = conditionUint
		m.build.conditions[conditionUint] = &condition
	}
	m.build.end = end
	return m
}
