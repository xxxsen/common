package database

type DBConfig struct {
	Host string `json:"host"`
	Port uint32 `json:"port"`
	User string `json:"user"`
	Pwd  string `json:"pwd"`
	DB   string `json:"db"`
}
