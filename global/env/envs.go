package env

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sethvargo/go-envconfig"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

// Envs contains all application configurations.
// The configuration are parsed from environment variables.
// For more info see: https://github.com/sethvargo/go-envconfig
type Envs struct {
	Server struct {
		GinMode         string `env:"GIN_MODE, default=debug"`
		LogLevel        string `env:"LOG_LEVEL, default=debug"`
		LogFormat       string `env:"LOG_FORMAT, default=text"`
		HttpPort        string `env:"HTTP_PORT, default=8090"`
		SwaggerHostAddr string `env:"SWAGGER_HOST_ADDR"`
		AssetsDir       string `env:"ASSETS_DIR, required"`
		JwtSecret       string `env:"JWT_SECRET, required"`
	}

	Db struct {
		Dsn      string `env:"DB_DSN, required"`
		LogLevel string `env:"DB_LOG_LEVEL, default=error"`
	}
}

// Load loads the environment variables from the .env files
// we use a base env file called ".env" and override it for different environment (e.g. dev, test, prod, or testing)
func Load(basePath string) (*Envs, error) {
	envFiles, err := filepath.Glob(filepath.Join(basePath, ".env*"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to find env files")
	}

	profileToEnvFile := map[string]string{}
	for _, envFile := range envFiles {
		// get the profile from the env file name e.g. .env.dev -> dev
		profile := filepath.Base(envFile)
		profile, _ = strings.CutPrefix(profile, ".env")
		profile, _ = strings.CutPrefix(profile, ".")
		profileToEnvFile[profile] = envFile
	}

	// first let load the env file for the current profile
	profile := os.Getenv("PROFILE")
	activeEnvFiles := make([]string, 0)
	if profile != "" {
		envFile, ok := profileToEnvFile[profile]
		if !ok {
			return nil, fmt.Errorf("no env file found for %q profile", profile)
		}
		activeEnvFiles = append(activeEnvFiles, envFile)
	}

	// then let's add the default profile as well if it exists
	if envFile, ok := profileToEnvFile[""]; ok {
		activeEnvFiles = append(activeEnvFiles, envFile)
	}

	if len(activeEnvFiles) == 0 {
		log.Print("no .env file found")
	}

	log.Printf(fmt.Sprintf("trying to load %v env files", activeEnvFiles))
	err = godotenv.Load(activeEnvFiles...)
	if err != nil {
		return nil, err
	}

	var configs Envs
	if err = envconfig.Process(context.Background(), &configs); err != nil {
		return nil, err
	}

	return &configs, nil
}
