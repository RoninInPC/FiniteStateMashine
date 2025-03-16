package main

type ConditionType any

type AlphabetType any

type FiniteStateMashine[conditionType ConditionType, alphabetType AlphabetType] interface {
	GetCurrentCondition() []conditionType
	GetAlphabet() []alphabetType
	ToNext(alphabetType) ([]conditionType, bool)
	IsFinal() bool
}
