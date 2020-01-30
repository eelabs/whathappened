package logging

import (
	"context"
	"fmt"
	"github.com/containerd/containerd/log"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	InvokedFunction    = "invokedFunction"
	ExecutionID        = "executionId"
	Environment        = "environment"
	EnvironmentKey     = "ENVIRONMENT"
	AwsLambdaFnNameKey = "AWS_LAMBDA_FUNCTION_NAME"
)

var (
	LogLevelKey = "LOG_LEVEL"
)

func ConfigureJSONFormatting(logLevelKey *string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	if logLevelKey == nil {
		logLevelKey = &LogLevelKey
	}
	level, err := logrus.ParseLevel(os.Getenv(*logLevelKey))
	if err != nil {
		level = logrus.WarnLevel
	}
	logrus.SetLevel(level)
}
func AppDetails(commitID *string, env *string) {
	logrus.WithFields(logrus.Fields{
		"commitId":    commitID,
		"environment": defaultEnv(env),
	}).Info("Application details")
}

func NewLambdaInvocation(fnName *string, env *string) *logrus.Entry {
	logger := log.L.WithFields(logrus.Fields{
		Environment:     defaultEnv(env),
		InvokedFunction: defaultFnName(fnName),
	})
	logger.Info("New lambda invoked")
	return logger
}

func defaultEnv(env *string) string {
	if env == nil {
		return os.Getenv(EnvironmentKey)
	}
	return os.Getenv(*env)
}

func defaultFnName(fnName *string) string {
	if fnName == nil {
		fn := os.Getenv(AwsLambdaFnNameKey)
		if len(fn) == 0 {
			return fmt.Sprintf("FN_NAME_%s", uuid.New().String())
		}
		return fn
	}
	return os.Getenv(*fnName)
}

func AddFields(ctx context.Context, f logrus.Fields) context.Context {
	logger := log.G(ctx).WithFields(f)
	logger.Info()
	return log.WithLogger(ctx, logger)
}

// help track execution
type (
	executionKey struct{}
)

func WithExecutionID(ctx context.Context) context.Context {
	executionID := ctx.Value(executionKey{})
	if executionID == nil {
		executionID = uuid.New().String()
		ctx = context.WithValue(ctx, executionKey{}, executionID)
	}
	return AddFields(ctx, logrus.Fields{
		ExecutionID: executionID,
	})
}
func WithExecutionIDLogStart(ctx context.Context) context.Context {
	c := WithExecutionID(ctx)
	ExecutionStart(c)
	return c
}
func ExecutionStart(ctx context.Context) {
	now := time.Now().UTC()
	log.G(ctx).WithFields(logrus.Fields{
		"executionStart":         now.Unix(),
		"executionStartReadable": now,
	}).Info("START execution metric")
}
func ExecutionEnd(ctx context.Context) {
	now := time.Now().UTC()
	log.G(ctx).WithFields(logrus.Fields{
		"executionEnd":         now.Unix(),
		"executionEndReadable": now,
	}).Info("END execution metric")
}
