package config

import (
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type (
	Opts struct {
		Config    any
		Paths     []string
		Filenames []string
	}
)

func Load(opts Opts) error {
	configLoaded := false
	loadedFiles := make([]string, 0)

	for _, f := range opts.Filenames {
		for _, p := range opts.Paths {
			fp := filepath.Join(p, f)
			if _, fileErr := os.Stat(fp); fileErr != nil {
				log.Debug().Str("file", fp).Msg("Config file not found, skipping")
				continue
			}
			if err := cleanenv.ReadConfig(fp, opts.Config); err != nil {
				log.Warn().Str("file", fp).Err(err).Msg("Failed to parse config file, will use env vars")
				continue
			}
			configLoaded = true
			loadedFiles = append(loadedFiles, fp)
			log.Info().Str("file", fp).Msg("Config file loaded successfully")
		}
	}

	if err := cleanenv.ReadEnv(opts.Config); err != nil {
		return err
	}

	if configLoaded {
		log.Info().Strs("files", loadedFiles).Msg("Configuration loaded from files")
		log.Info().Msg("Environment variables will override file values if set")
	} else {
		log.Info().Msg("No config files found, using environment variables only")
	}

	return nil
}
