package domain

type SagaStep struct {
	ID string

	TransactionID string
	StepName      string
	Status        string
}
