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
	mockstore "github.com/august-kuhfuss/hotdamn/store/mock_store"
	sqlitestore "github.com/august-kuhfuss/hotdamn/store/sqlite_store"
	"github.com/august-kuhfuss/hotdamn/tasks"
	fetchtemperature "github.com/august-kuhfuss/hotdamn/tasks/fetch_temperature"
	"github.com/urfave/cli/v2"
)

var (
	defaultSqliteFilePath = path.Join(defaultDataDir, "hotdamn.db")

	s        store.Store
	ts       []tasks.Task
	demoMode = false
	cmd      = &cli.App{
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
						Usage:      "thermometer ip addresses",
						Name:       "thermometer-ips",
						EnvVars:    []string{"THERMOMETER_IPS"},
						Required:   true,
						HasBeenSet: true,
					},
					&cli.IntFlag{
						Usage:   "temperature fetch interval in seconds",
						Name:    "temperature-fetch-interval",
						Aliases: []string{"interval"},
						EnvVars: []string{"TEMPERATURE_FETCH_INTERVAL"},
						Value:   defaultFetchIntervalSeconds,
					},
				},
				Before: func(ctx *cli.Context) error {
					var err error
					if demoMode {
						s = mockstore.NewStore()
						slog.Info("demo mode enabled")
						return nil
					}

					s, err = sqlitestore.NewStore(ctx.String("sqlite-file-path"))
					if err != nil {
						slog.Error("unable to create database store", slog.String("msg", err.Error()))
						return err
					}
					defer sqlitestore.Close()

					if err := sqlitestore.Ping(); err != nil {
						slog.Error("unable to ping database", slog.String("msg", err.Error()))
						return err
					}

					tips := ctx.StringSlice("thermometer-ips")
					tfi := ctx.Int("temperature-fetch-interval")
					ts = []tasks.Task{
						fetchtemperature.NewTask(tips, time.Duration(tfi)*time.Second),
					}

					return nil
				},
				Action: func(ctx *cli.Context) error {
					port := ctx.Int("port")
					httpSrv := &http.Server{
						Addr:    fmt.Sprintf(":%d", port),
						Handler: handler.New(s),
					}

					context, stop := signal.NotifyContext(ctx.Context, os.Interrupt, syscall.SIGTERM)
					defer stop()

					if !demoMode {
						for _, task := range ts {
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
