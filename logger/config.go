package logger

type LogConfig struct {
	File      string `json:"file" yaml:"file"`
	Level     string `json:"level" yaml:"level"`
	FileSize  uint64 `json:"file_size" yaml:"file_size"`
	FileCount uint64 `json:"file_count" yaml:"file_count"`
	KeepDays  uint32 `json:"keep_days" yaml:"keep_days"`
	Console   bool   `json:"console" yaml:"console"`
}
