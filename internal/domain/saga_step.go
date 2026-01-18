package domain

type SagaStep struct {
	ID string

	TransactionID string
	TxState       string
	StepName      string
	Status        string
}
