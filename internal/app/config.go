package app

import (
	"os"
	"strconv"
)

type Config struct {
	JWTPublicKeyPath  string
	JWTIssuer         string
	MessageChannelCap int
	MessageWorkersCap int
}

func NewConfig() *Config {

	mChanCap, err := strconv.Atoi(os.Getenv("MCHANCAP"))
	if err != nil {
		return nil
	}

	mWorkerCap, err := strconv.Atoi(os.Getenv("MWORKERSCAP"))
	if err != nil {
		return nil
	}

	return &Config{
		MessageChannelCap: mChanCap,
		MessageWorkersCap: mWorkerCap,
		JWTPublicKeyPath:  os.Getenv("JWT_PUBLIC_KEY_PATH"),
		JWTIssuer:         os.Getenv("JWT_ISSUER"),
	}
}
