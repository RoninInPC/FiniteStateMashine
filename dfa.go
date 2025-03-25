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

func (D *DFA[conditionType, alphabetType]) ToNext(alphabetType alphabetType) ([]conditionType, bool) {
	alphabetTypeUint, _ := D.serializationAlphabet(alphabetType)
	D.current = D.body[D.current][alphabetTypeUint]

	return D.GetCurrentCondition(), D.IsFinal()
}

func (D *DFA[conditionType, alphabetType]) IsFinal() bool {
	return slices.Contains(D.end, D.current)
}

func (D *DFA[conditionType, alphabetType]) Reset() {
	D.current = D.start
}
