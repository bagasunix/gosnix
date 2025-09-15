package errors

import (
	"fmt"
	"os"

	"github.com/phuslu/log"
)

func HandlerWithOSExit(logger *log.Logger, err error, args ...interface{}) {
	if err == nil {
		return
	}
	logger.Error().Err(err).Any("args", args).Msg("Fatal error, exiting...")
	os.Exit(1)
}

func HandlerWithLoggerReturnedError(logger *log.Logger, err error, args ...interface{}) error {
	if err == nil {
		return nil
	}
	logger.Error().Err(err).Any("args", args).Msg("Error occurred")
	return err
}

func HandlerWithLoggerReturnedVoid(logger *log.Logger, err error, args ...interface{}) {
	if err == nil {
		return
	}
	logger.Error().Err(err).Any("args", args).Msg("Error occurred")
}

func HandlerReturnedVoid(err error, args ...interface{}) {
	if err == nil {
		return
	}
	fmt.Println("err:", err, "args:", args)
}
