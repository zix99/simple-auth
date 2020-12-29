package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
)

func watchForCleanShutdown(e *echo.Echo) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Received SIGINT, shutting down...")
		e.Shutdown(context.Background())
	}()
}
