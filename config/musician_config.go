package config

type MusicianConf struct {
	// Root data path
	DataDir string `conf:"default:./conductor,flag:d" json:"dataDir"`

	InitAndExit bool `conf:"default:false,flag:x" json:"-"`

	Rest Rest `json:"rest"`

	Conductor Conductor `json:"conductor"`

	Logger Logger `json:"logger"`
}

type Conductor struct {
	AdvertiseAddr string `conf:"default:http://localhost:8090" json:"advertiseAddr"`
	ConductorAddr string `conf:"default:http://localhost:8080" json:"conductorAddr"`
}
