package test

import (
	"os"
	"time"

	"github.com/amahdian/ai-assistant-be/pkg/logger"
	"github.com/amahdian/ai-assistant-be/pkg/logger/logging"
)

func SetupTestingEnv() {
	time.Local = time.UTC
	logger.Configure(logging.DebugLevel, logging.TextFormat)
	_ = os.Setenv("PROFILE", "testing")
}
