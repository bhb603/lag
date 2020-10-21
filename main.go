package main

import (
	"github.com/bhb603/lag/cmd"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	cmd.Execute()
}
