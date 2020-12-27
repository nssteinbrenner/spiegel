package utils

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func InitLogger() *logrus.Entry {
	logger := logrus.WithFields(logrus.Fields{
		"correlation_id": uuid.New().String(),
	})

	return logger
}
