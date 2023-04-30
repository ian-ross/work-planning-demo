package server

type Config struct {
	// DevMode is a development mode flag: if true, logging is in a
	// human-readable format; if false, logging is in a JSON format.
	DevMode bool `env:"DEV_MODE,default=false"`

	// DBURL is the store connection URL. This is either a Postgres
	// connection URL, or "memory" to use the in-memory store.
	StoreURL string `env:"STORE_URL,required"`

	// Port is the port to run the HTTP server on.
	Port int `env:"PORT,default=8080"`

	// AccessTokenLease is the time (in seconds) for which a JWT access
	// token is valid.
	AccessTokenLease int `env:"ACCESS_TOKEN_LEASE,default=600"`

	// RefreshTokenLease is the time (in seconds) for which a JWT
	// refresh token is valid.
	RefreshTokenLease int `env:"REFRESH_TOKEN_LEASE,default=86400"`

	// AuthKey is a secret string used for generating JWT tokens.
	AuthKey string `env:"AUTH_KEY,required"`
}
