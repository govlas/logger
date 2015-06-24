package logger_test

import (
	"testing"

	"github.com/govlas/logger"
)

func TestLogger(t *testing.T) {

	logger.Info("Info")
	logger.Warning("Warning")
	logger.Error("Error")
	logger.Debug("Debug")

	logger.EnableColored()

	logger.Info("Info colored")
	logger.Warning("Warning colored")
	logger.Error("Error colored")
	logger.Debug("Debug colored")

	logger.DisableBTrace()

	logger.Info("Info colored")
	logger.Warning("Warning colored")
	logger.Error("Error colored")
	logger.Debug("Debug colored")

	logger.SetFileName(logger.FileNameShort)

	logger.Info("Info short filename")
	logger.Warning("Warning short filename")
	logger.Error("Error short filename")
	logger.Debug("Debug short filename")

	logger.SetFileName(logger.FileNameLong)

	logger.Info("Info long filename")
	logger.Warning("Warning long filename")
	logger.Error("Error long filename")
	logger.Debug("Debug long filename")
}
