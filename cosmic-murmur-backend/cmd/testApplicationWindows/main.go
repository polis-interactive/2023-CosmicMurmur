package main

import (
	"github.com/polis-interactive/2023-CosmicMurmur/data"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/application"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conf := &application.Config{
		LightingConfig: &application.LightingConfig{
			SegmentDefinition: data.DefaultLightingSegmentDefinition,
			SegmentCount:      1,
		},
		GraphicsConfig: &application.GraphicsConfig{
			DefaultShader:  "basic",
			PixelSize:      7,
			Frequency:      33 * time.Millisecond,
			ReloadOnUpdate: true,
		},
		ControllerConfig: &application.ControllerConfig{
			LocalAddress:    "2.0.0.1",
			NodeDefinitions: data.DefaultNodeDefinitions,
		},
		ServiceBusConfig: &application.ServiceBusConfig{
			EventQueueSize: 50,
			BusyTimeout:    1 * time.Second,
		},
		ProgramName: "cosmic-murmur-backend",
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().
		Str("method", "main").Msg("starting")

	app, err := application.NewApplication(conf)
	if err != nil {
		log.Panic().
			Str("method", "main").Err(err).Msg("couldn't create application instance; closing")
		panic(err)
	}

	err = app.Startup()
	if err != nil {
		log.Panic().
			Str("method", "main").Err(err).Msg("couldn't startup; shutting down")

		err2 := app.Shutdown()
		if err2 != nil {
			log.Panic().
				Str("method", "main").Err(err2).Msg("couldn't force shut down")
		}
		panic(err)
	}

	log.Info().Msg("Main: running")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info().
		Str("method", "main").Msg("closing")

	err = app.Shutdown()
	if err != nil {
		log.Panic().
			Str("method", "main").Err(err).Msg("issue shutting down")
	}

	log.Info().
		Str("method", "main").Msg("closed")
}
