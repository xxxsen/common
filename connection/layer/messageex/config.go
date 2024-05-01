package messageex

type action struct {
	Type      string `json:"type"`
	MinLength uint   `json:"min_length"` //valid when type == send
	MaxLength uint   `json:"max_length"` //valid when type == send
}

type config struct {
	Actions []string `json:"actions"`
}
