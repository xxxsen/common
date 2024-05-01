package message

type action struct {
	Type   string `json:"type"`
	Length uint   `json:"length"`
}

type config struct {
	Actions []*action `json:"actions"`
}
