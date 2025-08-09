package infra

import "os"

type OIDCConfig struct {
	Issuer   string
	Audience string
	SkipVerify bool
}

type Config struct {
	DatabaseURL string
	DBMaxConns  int
	OIDC        OIDCConfig
}

func LoadConfigFromEnv() Config {
	db := os.Getenv("DATABASE_URL")
	if db == "" {
		db = "postgres://scheduler:secret@localhost:5432/scheduler?sslmode=disable"
	}
	skip := false
	if os.Getenv("OIDC_SKIP_VERIFY") == "1" { skip = true }
	return Config{
		DatabaseURL: db,
		DBMaxConns:  5,
		OIDC: OIDCConfig{
			Issuer: os.Getenv("OIDC_ISSUER"),
			Audience: os.Getenv("OIDC_AUD") ,
			SkipVerify: skip,
		},
	}
}
