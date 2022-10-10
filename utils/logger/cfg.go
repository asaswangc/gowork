package logger

type Cfg struct {
	MaxAge      int    `json:"max_age" toml:"max_age"`
	Maxsize     int    `json:"max_size" toml:"max_size"`
	ErrLog      string `json:"err_log" toml:"err_log"`
	WarnLog     string `json:"warn_log" toml:"warn_log"`
	InfoLog     string `json:"info_log" toml:"info_log"`
	MaxBackups  int    `json:"max_backups" toml:"max_backups"`
	Compress    bool   `json:"compress" toml:"compress"`
	JsonEncoder bool   `json:"json_encoder" toml:"json_encoder"`
}
