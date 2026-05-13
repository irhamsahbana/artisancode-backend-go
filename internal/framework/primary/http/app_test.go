package http

import (
	"context"
	"errors"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

func TestAppRunReturnsBuildError(t *testing.T) {
	wantErr := errors.New("build failed")
	app := NewApp(AppConfig{})
	app.buildFn = func(appName string, appVersion string, appEnvironment string) (*fiber.App, error) {
		return nil, wantErr
	}
	app.serveFn = func(app *fiber.App, port string) error {
		t.Fatal("serve should not be called when build fails")
		return nil
	}
	app.waitForShutdownFn = func(ctx context.Context) error {
		t.Fatal("waitForShutdown should not be called when build fails")
		return nil
	}

	err := app.Run(context.Background(), "app", "v1", "test", "3000")

	require.ErrorIs(t, err, wantErr)
}

func TestAppRunReturnsServeErrorBeforeShutdown(t *testing.T) {
	wantErr := errors.New("listen failed")
	app := NewApp(AppConfig{})
	app.buildFn = func(appName string, appVersion string, appEnvironment string) (*fiber.App, error) {
		return fiber.New(), nil
	}
	app.serveFn = func(app *fiber.App, port string) error {
		return wantErr
	}
	app.waitForShutdownFn = func(ctx context.Context) error {
		t.Fatal("waitForShutdown should not be called when serve fails first")
		return nil
	}

	err := app.Run(context.Background(), "app", "v1", "test", "3000")

	require.ErrorIs(t, err, wantErr)
}

func TestAppRunReturnsShutdownError(t *testing.T) {
	wantErr := errors.New("shutdown failed")
	releaseServe := make(chan struct{})
	app := NewApp(AppConfig{})
	app.buildFn = func(appName string, appVersion string, appEnvironment string) (*fiber.App, error) {
		return fiber.New(), nil
	}
	app.serveFn = func(app *fiber.App, port string) error {
		<-releaseServe
		return nil
	}
	app.waitForShutdownFn = func(ctx context.Context) error {
		close(releaseServe)
		return wantErr
	}

	ctx, cancel := context.WithCancel(context.Background())
	go cancel()

	err := app.Run(ctx, "app", "v1", "test", "3000")

	require.ErrorIs(t, err, wantErr)
}

func TestAppRunReturnsNilOnGracefulShutdown(t *testing.T) {
	releaseServe := make(chan struct{})
	app := NewApp(AppConfig{})
	app.buildFn = func(appName string, appVersion string, appEnvironment string) (*fiber.App, error) {
		return fiber.New(), nil
	}
	app.serveFn = func(app *fiber.App, port string) error {
		<-releaseServe
		return nil
	}
	app.waitForShutdownFn = func(ctx context.Context) error {
		close(releaseServe)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	go cancel()

	err := app.Run(ctx, "app", "v1", "test", "3000")

	require.NoError(t, err)
}

func TestAppServeUsesConfiguredPort(t *testing.T) {
	fiberApp := fiber.New()
	var gotAddr string
	app := NewApp(AppConfig{})
	app.localIPv4sFn = func() []string {
		return []string{"192.168.0.10"}
	}
	app.listenFn = func(app *fiber.App, addr string) error {
		gotAddr = addr
		return nil
	}

	err := app.serve(fiberApp, "3000")

	require.NoError(t, err)
	require.Equal(t, ":3000", gotAddr)
}
