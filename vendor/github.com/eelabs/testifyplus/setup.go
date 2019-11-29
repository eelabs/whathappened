package testifyplus

import (
	"context"
	"github.com/containerd/containerd/log"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func InitialiseLogrus() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func UnitTestContextLogger(t *testing.T) (context.Context, *logrus.Entry) {
	ctx := context.TODO()
	logger := log.G(ctx).WithFields(logrus.Fields{
		"environment": "unit-test",
		"testing":     t.Name(),
	})
	logger.Info()
	ctx = log.WithLogger(ctx, logger)
	return ctx, log.G(ctx)
}
