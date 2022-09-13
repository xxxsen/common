package logger

type LogConfig struct {
	File      string `json:"file"`
	Level     string `json:"level"`
	FileSize  uint64 `json:"file_size"`
	FileCount uint64 `json:"file_count"`
	KeepDays  uint32 `json:"keep_days"`
	Console   bool   `json:"console"`
}
