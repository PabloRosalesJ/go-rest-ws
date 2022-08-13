package server

/* Config necessary to create a server */
type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

/* A server should be implements this interface */
type Server interface {
	Config() *Config
}
