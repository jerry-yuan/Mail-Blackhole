package config

import (
	"flag"

	"github.com/ian-kent/envconf"
)

type LogConfig struct {
	Level string
}

type SMTPConfig struct {
	BindAddr string
	Hostname string
}

type StorageConfig struct {
	StorageType string

	MongoURI  string
	MongoDb   string
	MongoColl string

	MaildirPath string
}

type APIConfig struct {
	BindAddr   string
	WebPath    string
	CORSOrigin string
}

type WebUIConfig struct {
	BindAddr string
	WebPath  string
	APIHost  string
	AuthFile string
}

type Config struct {
	PrintVersion bool

	Log LogConfig

	SMTP    SMTPConfig
	Storage StorageConfig
	API     APIConfig
	WebUI   WebUIConfig
}

var globalConfig = DefaultConfig()

func init() {
	// Global configuration
	flag.BoolVar(&globalConfig.PrintVersion, "version", false, "Print version")

	// Logrus configuration
	flag.StringVar(&globalConfig.Log.Level, "log-level", envconf.FromEnvP("MBH_LOG_LEVEL", "INFO").(string), "Log level(available:panic;fatal;error;warn;info;debug;trace)")

	// SMTP configuration
	flag.StringVar(&globalConfig.SMTP.BindAddr, "smtp-bind-addr", envconf.FromEnvP("MBH_SMTP_BIND_ADDR", "0.0.0.0:1025").(string), "SMTP bind interface and port, e.g. 0.0.0.0:1025 or just :1025")
	flag.StringVar(&globalConfig.SMTP.Hostname, "smtp-hostname", envconf.FromEnvP("MBH_SMTP_HOSTNAME", "mail.blackhole.example").(string), "Hostname for SMTP EHLO/HELO response, e.g. mail.blackhole.example")

	// Storage configuration
	flag.StringVar(&globalConfig.Storage.StorageType, "storage-type", envconf.FromEnvP("MBH_STORAGE", "memory").(string), "Message storage: 'memory' (default), 'mongodb' or 'maildir'")

	flag.StringVar(&globalConfig.Storage.MongoURI, "mongo-uri", envconf.FromEnvP("MBH_MONGO_URI", "127.0.0.1:27017").(string), "MongoDB URI, e.g. 127.0.0.1:27017")
	flag.StringVar(&globalConfig.Storage.MongoDb, "mongo-db", envconf.FromEnvP("MBH_MONGO_DB", "mailhog").(string), "MongoDB database, e.g. mailhog")
	flag.StringVar(&globalConfig.Storage.MongoColl, "mongo-coll", envconf.FromEnvP("MBH_MONGO_COLLECTION", "messages").(string), "MongoDB collection, e.g. messages")

	flag.StringVar(&globalConfig.Storage.MaildirPath, "maildir-path", envconf.FromEnvP("MBH_MAILDIR_PATH", "").(string), "Maildir path (if storage type is 'maildir')")

	// API configuration
	flag.StringVar(&globalConfig.API.BindAddr, "api-bind-addr", envconf.FromEnvP("MBH_API_BIND_ADDR", "0.0.0.0:8025").(string), "HTTP bind interface and port for API, e.g. 0.0.0.0:8025 or just :8025")
	flag.StringVar(&globalConfig.API.WebPath, "api-web-path", envconf.FromEnvP("MBH_API_WEB_PATH", "").(string), "WebPath under whitch the API is served (without leading or trailing slashes), e.g. 'mailblackhole'. Value defaults to ''")
	flag.StringVar(&globalConfig.API.CORSOrigin, "api-cors-origin", envconf.FromEnvP("MBH_API_CORS_ORIGIN", "").(string), "CORS Access-Control-Allow-Origin header for API endpoints")

	// Web UI configuration
	flag.StringVar(&globalConfig.WebUI.BindAddr, "ui-bind-addr", envconf.FromEnvP("MBH_WEBUI_BIND_ADDR", "0.0.0.0:8025").(string), "HTTP bind interface and port for UI, e.g. 0.0.0.0:8025 or just :8025")
	flag.StringVar(&globalConfig.WebUI.AuthFile, "ui-auth-file", envconf.FromEnvP("MBH_WEBUI_AUTH_FILE", "").(string), "A username:bcryptpw mapping file")
	flag.StringVar(&globalConfig.WebUI.WebPath, "ui-web-path", envconf.FromEnvP("MBH_UI_WEB_PATH", "").(string), "WebPath under which the UI is served (without leading or trailing slashes), e.g. 'mailhog'. Value defaults to ''")
	flag.StringVar(&globalConfig.WebUI.APIHost, "ui-api-host", envconf.FromEnvP("MBH_API_HOST", "").(string), "API URL for MailHog UI to connect to, e.g. http://some.host:1234")

	flag.Parse()
}

func DefaultConfig() *Config {
	return &Config{}
}

func Configure() *Config {

	/* WebUI Initialize */
	//sanitize webpath
	//add a leading slash
	if globalConfig.WebUI.WebPath != "" && !(globalConfig.WebUI.WebPath[0] == '/') {
		globalConfig.WebUI.WebPath = "/" + globalConfig.WebUI.WebPath
	}

	return globalConfig
}
