package logging

import (
	"context"
	"github.com/containerd/containerd/log"
	"github.com/eelabs/testifyplus"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func init() {
	testifyplus.InitialiseLogrus()
}

func TestWithExecutionID(t *testing.T) {
	ctx, _ := testifyplus.UnitTestContextLogger(t)

	ctx = WithExecutionID(ctx)
	eID := log.GetLogger(ctx).Data[ExecutionID]
	executionID, err := uuid.Parse(eID.(string))
	assert.NoError(t, err)
	assert.NotNil(t, executionID)

	ctx = WithExecutionID(ctx)
	eID2 := log.GetLogger(ctx).Data[ExecutionID]
	executionID2, err := uuid.Parse(eID2.(string))
	assert.NoError(t, err)
	assert.Equal(t, executionID, executionID2)

	//new context
	ctx = context.Background()
	ctx = WithExecutionIDLogStart(ctx)
	eID3 := log.GetLogger(ctx).Data[ExecutionID]
	executionID3, err := uuid.Parse(eID3.(string))
	assert.NoError(t, err)
	assert.NotEqual(t, executionID, executionID3)
	assert.NotEqual(t, executionID2, executionID3)
}

func TestNewLambdaInvocation(t *testing.T) {
	env := "PRODUCTION"
	envKey := "SYS_ENVIRONMENT_KEY"
	assert.NoError(t, os.Setenv(envKey, env))
	funcName := "funcBeingExecuted"
	funcNameKey := "PROVIDED_FUNC_NAME_KEY"
	assert.NoError(t, os.Setenv(funcNameKey, funcName))
	logger := NewLambdaInvocation(&funcNameKey, &envKey)

	actualFuncName := logger.Data[InvokedFunction]
	actualEnv := logger.Data[Environment]
	assert.Equal(t, funcName, actualFuncName)
	assert.Equal(t, env, actualEnv)

	env = "DEFAULTED_ENV"
	assert.NoError(t, os.Setenv(EnvironmentKey, env))
	funcName = "lambdaFn"
	assert.NoError(t, os.Setenv(AwsLambdaFnNameKey, funcName))
	logger = NewLambdaInvocation(nil, nil)

	actualEnv = logger.Data[Environment]
	assert.Equal(t, env, actualEnv)
	actualFuncName = logger.Data[InvokedFunction]
	assert.Equal(t, funcName, actualFuncName)
}

// not a test just a convenience to execute func
func Test_ExecutionStart(t *testing.T) {
	ctx, _ := testifyplus.UnitTestContextLogger(t)

	ExecutionStart(ctx)
}

func Test_ExecutionEnd(t *testing.T) {
	ctx, _ := testifyplus.UnitTestContextLogger(t)

	ExecutionEnd(ctx)
}
