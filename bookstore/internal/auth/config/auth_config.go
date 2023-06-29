package config

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const AuthSecret string = "secrety secret"

const HashingCost int = bcrypt.DefaultCost

const TokenExpTime time.Duration = time.Minute

const RepoChangestreamSleepSeconds = 60
