package fragment

type config struct {
	IntervalRange     []uint32 `json:"interval_range"`      //ms
	PacketLengthRange []uint32 `json:"packet_length_range"` //[x,y], 为空则按最大的包长进行发送
	PacketNumberRange []uint32 `json:"packet_number_range"` //[x,y], 包序号在范围内的的包都会被处理, 为空则所有包均不处理
}
