package config

import "time"

var Config *Cfg

var TimeLocation *time.Location

type Cfg struct {
	App        AppConfig
	Db         DB
	Credential Credential
	Jwt        Jwt
}

type AppConfig struct {
	Name     string
	Url      string
	Port     int
	Env      string
	Debug    bool
	Timezone string
}

type DB struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Name       string
	Connection DbConnConfig
}

type DbConnConfig struct {
	Open int
	TTL  int
	Idle int
}

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Jwt struct {
	Secret           string
	ExpiresIn        int32
	RefreshExpiresIn int32
}
