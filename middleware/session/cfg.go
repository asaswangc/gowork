package session

type Cfg struct {
	Public  string `json:"public" toml:"public"`
	Private string `json:"private" toml:"private"`
}

func NewCfg(public string, private string) *Cfg {
	return &Cfg{Public: public, Private: private}
}
