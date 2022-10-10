package mysql_cli

type Cfg struct {
	*LogCfg
	*ConnCfg
}

type LogCfg struct {
	MaxAge     int    `json:"max_age"`
	Maxsize    int    `json:"max_size"`
	LogPath    string `json:"log_path"`
	Compress   bool   `json:"compress"`
	MaxBackups int    `json:"max_backups"`
}

type ConnCfg struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	DBName       string `json:"db_name"`
	UserName     string `json:"user_name"`
	Password     string `json:"password"`
	MaxLifeTime  int    `json:"max_life_time"`
	MaxIdleTime  int    `json:"max_idle_time"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxOpenConns int    `json:"max_open_conns"`
}
