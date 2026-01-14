package domain

type CardEvent struct {
	EventID          string
	UserRefId        string
	SourceRefId      string
	DestinationRefId string
	DcFlag           string // D or C
	PaymentType      string // DEEFERRED(NOT IMMEDIATE)
	PaymentMode      string // IMPS / UPI / NEFT / CARD
	NetworkTxnId     string
	GatewayTxnId     string
	EventType        string
	NetworkTxn       string
	Amount           int64
	Source           string // CARD_SWITCH, NETWORK, CORE
}
