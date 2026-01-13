package events

type CardEvent struct {
	EventID     string `json:"eventId"`
	EventType   string `json:"eventType"`
	NetworkTxn  string `json:"networkTxnId"`
	Amount      int64  `json:"amount"`
}
