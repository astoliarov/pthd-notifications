package main

import (
	"context"
	"log"
	"pthd-notifications/pkg"
)

func main() {
	application, initErr := pkg.NewApplication()
	if initErr != nil {
		log.Fatalf("Failed to initialize application: %s", initErr)
	}

	ctx := context.Background()
	runErr := application.Run(ctx)
	if runErr != nil {
		log.Fatalf("Error of server: %s", runErr)
	}
}
