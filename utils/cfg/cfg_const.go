package cfg

type MysqlLogger struct {
	MaxAge     int    `json:"max_age" toml:"max_age"`
	Maxsize    int    `json:"max_size" toml:"max_size"`
	LogPath    string `json:"log_path" toml:"log_path"`
	Compress   bool   `json:"compress" toml:"compress"`
	MaxBackups int    `json:"max_backups" toml:"max_backups"`
}

type Mysql struct {
	Host         string      `json:"host" toml:"host"`
	Port         string      `json:"port" toml:"port"`
	DbName       string      `json:"db_name" toml:"db_name"`
	UserName     string      `json:"user_name" toml:"user_name"`
	Password     string      `json:"password" toml:"password"`
	MaxLifeTime  int         `json:"max_life_time" toml:"max_life_time"`
	MaxIdleTime  int         `json:"max_idle_time" toml:"max_idle_time"`
	MaxIdleConns int         `json:"max_idle_conns" toml:"max_idle_conns"`
	MaxOpenConns int         `json:"max_open_conns" toml:"max_open_conns"`
	Logger       MysqlLogger `json:"logger" toml:"logger"`
}

type Redis struct {
	Host        []string `json:"host"`
	Password    string   `json:"password" toml:"password"`
	PoolSize    int      `json:"pool_size" toml:"pool_size"`
	PoolTimeout int      `json:"pool_timeout" toml:"pool_timeout"`
}

type ZapLog struct {
	MaxAge      int    `json:"max_age" toml:"max_age"`
	Maxsize     int    `json:"max_size" toml:"max_size"`
	ErrLog      string `json:"err_log" toml:"err_log"`
	WarnLog     string `json:"warn_log" toml:"warn_log"`
	InfoLog     string `json:"info_log" toml:"info_log"`
	MaxBackups  int    `json:"max_backups" toml:"max_backups"`
	Compress    bool   `json:"compress" toml:"compress"`
	JsonEncoder bool   `json:"json_encoder" toml:"json_encoder"`
}

type KeyPath struct {
	Public  string `json:"public" toml:"public"`
	Private string `json:"private" toml:"private"`
}

type Toml struct {
	Mysql   Mysql   `json:"mysql" toml:"mysql"`
	Redis   Redis   `json:"redis" toml:"redis"`
	ZapLog  ZapLog  `json:"zap_log" toml:"zap_log"`
	KeyPath KeyPath `json:"key_path" toml:"key_path"`
}
