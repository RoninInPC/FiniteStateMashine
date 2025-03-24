package main

import "slices"

// NondeterministicFiniteAutomaton или Недетерменированный конечный автомат

type NFA[conditionType ConditionType, alphabetType AlphabetType] struct {
	body map[uint64]map[uint64][]uint64

	funcSerialization Serialization

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
	num, err := N.funcSerialization(alphabetType)
	if err != nil {
		return nil, false
	}

	conditionNext := make([]uint64, 0)
	for _, condition := range N.current {
		conditionNext = append(conditionNext, N.body[condition][num]...)
	}

	N.current = conditionNext
	answer := make([]conditionType, len(conditionNext))
	for i, condition := range conditionNext {
		answer[i] = *N.conditions[condition]
	}
	return answer, N.IsFinal()
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
	N.current = N.current[:0]
}

type MakerNFA[conditionType ConditionType, alphabetType AlphabetType] struct {
	build NFA[conditionType, alphabetType]
}

func DefaultNFA[conditionType ConditionType, alphabetType AlphabetType]() MakerNFA[conditionType, alphabetType] {
	return MakerNFA[conditionType, alphabetType]{
		build: NFA[conditionType, alphabetType]{
			funcSerialization: SerializationValueDefault},
	}
}

func SerializationNFA[conditionType ConditionType, alphabetType AlphabetType](
	serialization Serialization) MakerNFA[conditionType, alphabetType] {
	return MakerNFA[conditionType, alphabetType]{
		build: NFA[conditionType, alphabetType]{
			funcSerialization: serialization},
	}
}

func (m MakerNFA[conditionType, alphabetType]) Build() NFA[conditionType, alphabetType] {
	return m.build
}

func (m MakerNFA[conditionType, alphabetType]) SetAlphabet(alphabet []alphabetType) MakerNFA[conditionType, alphabetType] {

}
