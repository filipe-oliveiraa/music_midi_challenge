package config

type ConductorConf struct {
	// Root data path
	DataDir string `conf:"default:./conductor,flag:d" json:"dataDir"`

	InitAndExit bool `conf:"default:false,flag:x" json:"-"`

	Rest Rest `json:"rest"`

	Logger Logger `json:"logger"`
}
