package config

type Rest struct {
	// TLSCertFile is the certificate file
	TLSCertFile string `conf:"" json:"tlsCertFile"`

	// TLSKeyFile is the key file
	TLSKeyFile string `conf:"" json:"tlsKeyFile"`

	// EndpointAddress configures the address the node listens to for REST API calls.
	// Specify an IP and port or just port. For example,
	// 127.0.0.1:0 will listen on a random port on the localhost (preferring 8080).
	EndpointAddress string `conf:"default:127.0.0.1:0" json:"endpointAddress"`

	// IncomingConnectionsLimit specifies the max number of incoming connections.
	// 0 means no connections allowed. Must be non-negative.
	// Estimating 1.5MB per incoming connection, 1.5MB*2400 = 3.6GB
	IncomingConnectionsLimit int `conf:"default:100" json:"incomingConnectionsLimit"`

	// RestConnectionsSoftLimit is the maximum number of active requests the API server
	// When the number of http connections to the REST layer exceeds the soft limit,
	// we start returning http code 429 Too Many Requests.
	ConnectionsSoftLimit uint64 `conf:"default:1024" json:"connectionsSoftLimit"`

	// RestConnectionsHardLimit is the maximum number of active connections the API server
	// will accept before closing requests with no response.
	ConnectionsHardLimit uint64 `conf:"default:2048" json:"connectionsHardLimit"`

	// ReservedFDs is used to make sure the process does not run out of file descriptors (FDs).
	// Daemon ensures that RLIMIT_NOFILE >= IncomingConnectionsLimit + RestConnectionsHardLimit +
	// ReservedFDs. ReservedFDs are meant to leave room for short-lived FDs like
	// DNS queries, SQLite files, etc. This parameter shouldn't be changed.
	// If RLIMIT_NOFILE < IncomingConnectionsLimit + RestConnectionsHardLimit + ReservedFDs
	// then either RestConnectionsHardLimit or IncomingConnectionsLimit decreased.
	ReservedFDs uint64 `conf:"default:256" json:"reservedFDs"`

	// RestReadTimeoutSeconds is passed to the API servers rest http.Server implementation.
	ReadTimeoutSeconds int `conf:"default:15" json:"readTimeoutSeconds"`

	// RestWriteTimeoutSeconds is passed to the API servers rest http.Server implementation.
	WriteTimeoutSeconds int `conf:"default:120" json:"writeTimeoutSeconds"`
}

type Logger struct {
	// LogToStdout if set true will log to stdout
	LogToStdout bool `conf:"default:false,flag:o" json:"-"`

	// FileDir is an optional directory to store the log, node.log
	// If not specified, the node will use the HotDataDir.
	// The -o command line option can be used to override this output location.
	FileDir string `conf:"" json:"fileDir"`

	// ArchiveDir is an optional directory to store the log archive.
	// If not specified, the node will use the ColdDataDir.
	ArchiveDir string `conf:"" json:"archieveDir"`

	// BaseLoggerDebugLevel specifies the logging level (node.log).
	// The levels range from 0 (critical error / silent) to 5 (debug / verbose).
	// The default value is 4 (‘Info’ - fairly verbose).
	BaseLoggerDebugLevel int8 `conf:"default:4" json:"baseLoggerDebugLevel"`

	// ArchiveName text/template for creating log archive filename.
	// Available template vars:
	// Time at start of log: {{.Year}} {{.Month}} {{.Day}} {{.Hour}} {{.Minute}} {{.Second}}
	// Time at end of log: {{.EndYear}} {{.EndMonth}} {{.EndDay}} {{.EndHour}} {{.EndMinute}} {{.EndSecond}}
	//
	// If the filename ends with .gz or .bz2 it will be compressed.
	//
	// default: "node.archive.log" (no rotation, clobbers previous archive)
	LogArchiveName string `conf:"default:node.archive.log" json:"logArchiveName"`

	// ArchiveMaxAge will be parsed by time.ParseDuration().
	// Valid units are 's' seconds, 'm' minutes, 'h' hours
	LogArchiveMaxAge string `conf:"" json:"logArchiveMaxAge"`

	// SizeLimit is the log file size limit in bytes. When set to 0 logs will be written to stdout.
	LogSizeLimit uint64 `conf:"default:1073741824" json:"logSizeLimit"`
}
