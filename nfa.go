package main

// NondeterministicFiniteAutomaton
type NFA[conditionType ConditionType, alphabetType AlphabetType] struct {
	body     map[*conditionType]map[*alphabetType][]conditionType
	alphabet []alphabetType
	current  *conditionType
	start    *conditionType
	end      []*conditionType
}

func (N NFA[conditionType, alphabetType]) GetCurrentCondition() []conditionType {
	//TODO implement me
	panic("implement me")
}

func (N NFA[conditionType, alphabetType]) GetAlphabet() []alphabetType {
	//TODO implement me
	panic("implement me")
}

func (N NFA[conditionType, alphabetType]) ToNext(alphabetType alphabetType) ([]conditionType, bool) {
	//TODO implement me
	panic("implement me")
}

func (N NFA[conditionType, alphabetType]) IsFinal() bool {
	//TODO implement me
	panic("implement me")
}
