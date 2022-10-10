package redis_cli

type Cfg struct {
	Host        []string `json:"host"`
	Password    string   `json:"password"`
	PoolSize    int      `json:"pool_size"`
	PoolTimeout int      `json:"pool_timeout"`
}
