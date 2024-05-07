package config

import (
	"github.com/joho/godotenv"
)
func init() {
    initEnv()
}

// Loads required variables from configuration file.
func initEnv() {
    if err := godotenv.Load("./config/wbl0_vars.env"); err != nil {
        panic(err)
    }
}
