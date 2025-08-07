//go:build e2e

package e2e_test

import (
	"context"
	"errors"
	"mimic/modules/app"
	"testing"
	"time"
)

const (
	goMimicPort = 3000
)

func TestEndToEnd(t *testing.T) {
	app := &app.App{}
	if err := setupTest(&app); err != nil {
		t.Fatal(err)
	}

	defer teardownTest(app) //nolint:errcheck

	go func() {
		if err := app.Run(t.Context()); err != nil {
			if !errors.Is(err, context.Canceled) {
				panic(err)
			}
		}
	}()

	time.Sleep(time.Second) // waiting for the server to bind to local port

	// Perform end-to-end tests here
	t.Run("http handler", testHttpHandler)
}
