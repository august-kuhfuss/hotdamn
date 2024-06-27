package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	_ "embed"

	"github.com/august-kuhfuss/hotdamn"
	"github.com/august-kuhfuss/hotdamn/handler"
	"github.com/august-kuhfuss/hotdamn/store"
	"github.com/august-kuhfuss/hotdamn/store/mock"
	"github.com/august-kuhfuss/hotdamn/store/sqlite"
	"github.com/august-kuhfuss/hotdamn/tasks"
	"github.com/urfave/cli/v2"
)

var (
	demoMode              = false
	defaultSqliteFilePath = path.Join(defaultDataDir, "hotdamn.db")

	cmd = &cli.App{
		Name:    "hotdamn",
		Usage:   "check weather in room (yes, it's hot in here)",
		Version: hotdamn.Version(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Usage:      "data directory",
				Name:       "data-dir",
				Aliases:    []string{"d"},
				EnvVars:    []string{"DATA_DIR"},
				Value:      defaultDataDir,
				Required:   true,
				HasBeenSet: true,
				Action: func(ctx *cli.Context, s string) error {
					if err := os.MkdirAll(s, 0755); err != nil {
						return err
					}
					return nil
				},
			},
			&cli.StringFlag{
				Usage:      "sqlite file path",
				Name:       "sqlite-file-path",
				Aliases:    []string{"s"},
				EnvVars:    []string{"SQLITE_FILE_PATH"},
				Value:      defaultSqliteFilePath,
				Required:   true,
				HasBeenSet: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "database",
				Aliases: []string{"db"},
				Usage:   "database management",
				Subcommands: []*cli.Command{
					{
						Name:    "migrate-up",
						Aliases: []string{"up"},
						Usage:   "migrate database up",
						Action: func(ctx *cli.Context) error {
							return sqlite.MigrateUp()
						},
					},
					{
						Name:    "migrate-down",
						Aliases: []string{"down"},
						Usage:   "migrate database down",
						Action: func(ctx *cli.Context) error {
							return sqlite.MigrateDown()
						},
					},
				},
			},
			{
				Name:    "run",
				Aliases: []string{"start", "serve"},
				Usage:   "start the http server",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Usage:      "port to listen on",
						Name:       "port",
						Aliases:    []string{"p"},
						EnvVars:    []string{"PORT"},
						Value:      defaultHTTPPort,
						Required:   true,
						HasBeenSet: true,
					},
					&cli.StringSliceFlag{
						Usage:      "sensor ip addresses",
						Name:       "sensor-ips",
						EnvVars:    []string{"SENSOR_IPS"},
						Required:   true,
						HasBeenSet: true,
					},
					&cli.IntFlag{
						Usage:   "temperature fetch interval in seconds",
						Name:    "sensor-fetch-interval",
						Aliases: []string{"interval"},
						EnvVars: []string{"SENSOR_FETCH_INTERVAL"},
						Value:   defaultFetchIntervalSeconds,
					},
				},
				Action: func(ctx *cli.Context) error {
					var st store.Store
					if demoMode {
						st = mock.NewStore()
						slog.Info("demo mode enabled")
						return nil
					}

					var err error
					st, err = sqlite.NewStore(ctx.String("sqlite-file-path"))
					if err != nil {
						slog.Error("unable to create database store", slog.String("msg", err.Error()))
						return err
					}
					defer sqlite.Close()

					if err := sqlite.Ping(); err != nil {
						slog.Error("unable to ping database", slog.String("msg", err.Error()))
						return err
					}

					if err := sqlite.MigrateUp(); err != nil {
						slog.Error("unable to migrate database up", slog.String("msg", err.Error()))
						return err
					}

					tsks := []tasks.Task{
						tasks.NewFetchTemperatureTask(
							ctx.StringSlice("sensor-ips"),
							time.Duration(ctx.Int("sensor-fetch-interval"))*time.Second,
							st,
						),
					}

					port := ctx.Int("port")
					httpSrv := &http.Server{
						Addr:    fmt.Sprintf(":%d", port),
						Handler: handler.New(st),
					}

					context, stop := signal.NotifyContext(ctx.Context, os.Interrupt, syscall.SIGTERM)
					defer stop()

					if !demoMode {
						for _, task := range tsks {
							go task.Start(ctx.Context)
						}
					}

					go func() {
						if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
							slog.Error("unable to start http server", slog.String("msg", err.Error()))
						}
					}()
					slog.Info("http server started", slog.String("addr", httpSrv.Addr))

					<-context.Done()
					slog.Info("http server stopping")
					if err := httpSrv.Shutdown(ctx.Context); err != nil {
						slog.Error("unable to stop http server", slog.String("msg", err.Error()))
						return err
					}
					slog.Info("http server stopped")

					return nil
				},
			},
		},
	}
)

func main() {
	ctx := context.Background()
	if err := cmd.RunContext(ctx, os.Args); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	os.Exit(0)
}
