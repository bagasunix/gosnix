package errors

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/phuslu/log"
	"gorm.io/gorm"
)

const (
	ERR_SOMETHING_WRONG = "oops something wrong. dont worry we will fix it"
	ERR_NOT_FOUND       = "not found"
	ERR_DUPLICATE_KEY   = "duplicate key"
	ERR_ALREADY_EXISTS  = "is already exists"
	ERR_INVALID_KEY     = "invalid"
	ERR_NOT_AUTHORIZED  = "unauthorized"
)

func CustomError(err string) error {
	return errors.New(err)
}

func ErrRecordNotFound(logger *log.Logger, entity string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(fmt.Sprintf("%s %s", entity, ERR_NOT_FOUND))
	}
	return ErrSomethingWrong(logger, err)
}

func ErrDuplicateValue(logger *log.Logger, entity string, err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(strings.ToLower(err.Error()), ERR_DUPLICATE_KEY) {
		return errors.New(fmt.Sprintf("%s %s", entity, ERR_ALREADY_EXISTS))
	}
	return ErrSomethingWrong(logger, err)
}

func ErrSomethingWrong(logger *log.Logger, err error) error {
	if err == nil {
		return nil
	}
	logger.Error().Err(err).Msg("Something went wrong")
	return errors.New(ERR_SOMETHING_WRONG)
}

func ErrInvalidAttributes(attributes string) error {
	return errors.New(fmt.Sprintf("%s %s", ERR_INVALID_KEY, attributes))
}

func ErrUnAuthorized() error {
	return errors.New(ERR_NOT_AUTHORIZED)
}

func FatalError(logger *log.Logger, err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}
	logger.Fatal().Err(err).Any("args", args).Msg(msg)
	os.Exit(1)
}

func LogAndReturnError(logger *log.Logger, err error, msg string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	logger.Error().Err(err).Any("args", args).Msg(msg)
	return err
}

func LogError(logger *log.Logger, err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}
	logger.Error().Err(err).Any("args", args).Msg(msg)
}

func LogWarning(logger *log.Logger, err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}
	logger.Warn().Err(err).Any("args", args).Msg(msg)
}

func LogDebug(logger *log.Logger, err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}
	logger.Debug().Err(err).Any("args", args).Msg(msg)
}
func ErrDataAlready(entity string) error {
	return errors.New(fmt.Sprint(entity, " ", ERR_ALREADY_EXISTS))
}
