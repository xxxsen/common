package padding

type config struct {
	PaddingMin           uint `json:"padding_min"`              //default: 200
	PaddingMax           uint `json:"padding_max"`              //default: 1000
	PaddingIfLessThan    uint `json:"padding_if_less_than"`     //default: 4k
	MaxBusiDataPerPacket uint `json:"max_busi_data_per_packet"` //default: 16k
}
