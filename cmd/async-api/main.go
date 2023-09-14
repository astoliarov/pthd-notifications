package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"pthd-notifications/pkg"
)

func main() {
	application, initErr := pkg.NewApplication()
	if initErr != nil {
		log.Fatal().
			Err(initErr).
			Msgf("failed to initialize API")
	}

	ctx := context.Background()
	runErr := application.RunRedisConsumer(ctx)
	if runErr != nil {
		log.Fatal().
			Err(runErr).
			Msgf("AsyncAPI server error")
	}
}
