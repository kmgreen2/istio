// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package log provides the canonical logging functionality used by Go-based
// Istio components.
//
// Istio's logging subsystem is built on top of the [Zap](https://godoc.org/go.uber.org/zap) package.
// High performance scenarios should use the Error, Warn, Info, and Debug methods. Lower perf
// scenarios can use the more expensive convenience methods such as Debugf and Warnw.
//
// The package provides direct integration with the Cobra command-line processor which makes it
// easy to build programs that use a consistent interface for logging. Here's an example
// of a simple Cobra-based program using this log package:
//
//		func main() {
//			// get the default logging options
//			options := log.NewOptions()
//
//			rootCmd := &cobra.Command{
//				Run: func(cmd *cobra.Command, args []string) {
//
//					// configure the logging system
//					if err := log.Configure(options); err != nil {
//                      // print an error and quit
//                  }
//
//					// output some logs
//					log.Info("Hello")
//					log.Sync()
//				},
//			}
//
//			// add logging-specific flags to the cobra command
//			options.AttachCobraFlags(rootCmd)
//			rootCmd.SetArgs(os.Args[1:])
//			rootCmd.Execute()
//		}
//
// Once configured, this package intercepts the output of the standard golang "log" package as well as anything
// sent to the global zap logger (zap.L()).
package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
)

// Global variables against which all our logging occurs.
var logger *zap.Logger = zap.NewNop()
var sugar *zap.SugaredLogger = logger.Sugar()

// Configure initializes Istio's logging subsystem.
//
// You typically call this once at process startup.
// Once this call returns, the logging system is ready to accept data.
func Configure(options *Options) error {
	return configure(options, func(c *zap.Config) (*zap.Logger, error) { return c.Build() })
}

type builder func(c *zap.Config) (*zap.Logger, error)

func configure(options *Options, b builder) error {
	outputLevel, err := options.GetOutputLevel()
	if err != nil {
		return err
	}

	stackTraceLevel, err := options.GetStackTraceLevel()
	if err != nil {
		return err
	}

	if outputLevel == None {
		// stick with the Nop default
		logger = zap.NewNop()
		sugar = logger.Sugar()
		logger = zap.NewNop()
		sugar = logger.Sugar()
		return nil
	}

	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(outputLevel),
		Development: false,

		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},

		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stack",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},

		OutputPaths:       options.OutputPaths,
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     !options.IncludeCallerSourceLocation,
		DisableStacktrace: stackTraceLevel == None,
	}

	if options.JSONEncoding {
		zapConfig.Encoding = "json"
	}

	l, err := b(&zapConfig)
	if err != nil {
		return err
	}

	logger = l.WithOptions(zap.AddCallerSkip(1), zap.AddStacktrace(stackTraceLevel))
	sugar = logger.Sugar()

	// capture global zap logging and force it through our logger
	_ = zap.ReplaceGlobals(l)

	// capture standard golang "log" package output and force it through our logger
	_ = zap.RedirectStdLog(logger)

	// capture gRPC logging
	grpclog.SetLogger(zapgrpc.NewLogger(logger.WithOptions(zap.AddCallerSkip(2))))

	return nil
}

// Debug outputs a message at debug level.
// This call is a wrapper around [Logger.Debug](https://godoc.org/go.uber.org/zap#Logger.Debug)
func Debug(msg string, fields ...zapcore.Field) {
	logger.Debug(msg, fields...)
}

// Debuga uses fmt.Sprint to construct and log a message at debug level.
// This call is a wrapper around [Sugaredlogger.Debug](https://godoc.org/go.uber.org/zap#Sugaredlogger.Debug)
func Debuga(args ...interface{}) {
	sugar.Debug(args...)
}

// Debugf uses fmt.Sprintf to construct and log a message at debug level.
// This call is a wrapper around [Sugaredlogger.Debugf](https://godoc.org/go.uber.org/zap#Sugaredlogger.Debugf)
func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

// Debugw logs a message at debug level with some additional context.
// This call is a wrapper around [Sugaredlogger.Debugw](https://godoc.org/go.uber.org/zap#Sugaredlogger.Debugw)
func Debugw(msg string, keysAndValues ...interface{}) {
	sugar.Debugw(msg, keysAndValues...)
}

// DebugEnabled returns whether output of messages at the debug level is currently enabled.
func DebugEnabled() bool {
	return logger.Core().Enabled(zap.DebugLevel)
}

// Error outputs a message at error level.
// This call is a wrapper around [logger.Error](https://godoc.org/go.uber.org/zap#logger.Error)
func Error(msg string, fields ...zapcore.Field) {
	logger.Error(msg, fields...)
}

// Errora uses fmt.Sprint to construct and log a message at error level.
// This call is a wrapper around [Sugaredlogger.Error](https://godoc.org/go.uber.org/zap#Sugaredlogger.Error)
func Errora(args ...interface{}) {
	sugar.Error(args...)
}

// Errorf uses fmt.Sprintf to construct and log a message at error level.
// This call is a wrapper around [Sugaredlogger.Errorf](https://godoc.org/go.uber.org/zap#Sugaredlogger.Errorf)
func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

// Errorw logs a message at error level with some additional context.
// This call is a wrapper around [Sugaredlogger.Errorw](https://godoc.org/go.uber.org/zap#Sugaredlogger.Errorw)
func Errorw(msg string, keysAndValues ...interface{}) {
	sugar.Errorw(msg, keysAndValues...)
}

// ErrorEnabled returns whether output of messages at the error level is currently enabled.
func ErrorEnabled() bool {
	return logger.Core().Enabled(zap.ErrorLevel)
}

// Warn outputs a message at warn level.
// This call is a wrapper around [logger.Warn](https://godoc.org/go.uber.org/zap#logger.Warn)
func Warn(msg string, fields ...zapcore.Field) {
	logger.Warn(msg, fields...)
}

// Warna uses fmt.Sprint to construct and log a message at warn level.
// This call is a wrapper around [Sugaredlogger.Warn](https://godoc.org/go.uber.org/zap#Sugaredlogger.Warn)
func Warna(args ...interface{}) {
	sugar.Warn(args...)
}

// Warnf uses fmt.Sprintf to construct and log a message at warn level.
// This call is a wrapper around [Sugaredlogger.Warnf](https://godoc.org/go.uber.org/zap#Sugaredlogger.Warnf)
func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

// Warnw logs a message at warn level with some additional context.
// This call is a wrapper around [Sugaredlogger.Warnw](https://godoc.org/go.uber.org/zap#Sugaredlogger.Warnw)
func Warnw(msg string, keysAndValues ...interface{}) {
	sugar.Warnw(msg, keysAndValues...)
}

// WarnEnabled returns whether output of messages at the warn level is currently enabled.
func WarnEnabled() bool {
	return logger.Core().Enabled(zap.WarnLevel)
}

// Info outputs a message at information level.
// This call is a wrapper around [logger.Info](https://godoc.org/go.uber.org/zap#logger.Info)
func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}

// Infoa uses fmt.Sprint to construct and log a message at info level.
// This call is a wrapper around [Sugaredlogger.Info](https://godoc.org/go.uber.org/zap#Sugaredlogger.Info)
func Infoa(args ...interface{}) {
	sugar.Info(args...)
}

// Infof uses fmt.Sprintf to construct and log a message at info level.
// This call is a wrapper around [Sugaredlogger.Infof](https://godoc.org/go.uber.org/zap#Sugaredlogger.Infof)
func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

// Infow logs a message at info level with some additional context.
// This call is a wrapper around [Sugaredlogger.Infow](https://godoc.org/go.uber.org/zap#Sugaredlogger.Infow)
func Infow(msg string, keysAndValues ...interface{}) {
	sugar.Infow(msg, keysAndValues...)
}

// InfoEnabled returns whether output of messages at the info level is currently enabled.
func InfoEnabled() bool {
	return logger.Core().Enabled(zap.InfoLevel)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
// This call is a wrapper around [logger.With](https://godoc.org/go.uber.org/zap#logger.With)
func With(fields ...zapcore.Field) *zap.Logger {
	return logger.With(fields...)
}

// Sync flushes any buffered log entries.
// Processes should normally take care to call Sync before exiting.
// This call is a wrapper around [logger.Sync](https://godoc.org/go.uber.org/zap#logger.Sync)
func Sync() {
	logger.Sync()
}
