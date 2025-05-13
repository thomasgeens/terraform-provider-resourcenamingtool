// Copyright (c) Thomas Geens

package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// LogDebug logs a debug message to both tflog and standard logger
// This centralizes the logging pattern used throughout the provider
func logDebug(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	tflog.Debug(ctx, msg)
	log.Println("[DEBUG]", msg)
}

// LogInfo logs an info message to both tflog and standard logger
func logInfo(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	tflog.Info(ctx, msg)
	log.Println("[INFO]", msg)
}

// LogWarn logs a warning message to both tflog and standard logger
//
//lint:ignore U1000 Function is currently unused, but may be used in the future
func logWarn(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	tflog.Warn(ctx, msg)
	log.Println("[WARN]", msg)
}

// LogError logs an error message to both tflog and standard logger
//
//lint:ignore U1000 Function is currently unused, but may be used in the future
func logError(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	tflog.Error(ctx, msg)
	log.Println("[ERROR]", msg)
}

// LogTrace logs a trace message to both tflog and standard logger
//
//lint:ignore U1000 Function is currently unused, but may be used in the future
func logTrace(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	tflog.Trace(ctx, msg)
	log.Println("[TRACE]", msg)
}

// LogDebugWithFields logs a debug message with structured fields to both tflog and standard logger
func logDebugWithFields(ctx context.Context, msg string, fields map[string]interface{}) {
	tflog.Debug(ctx, msg, fields)

	// For standard logger, format fields as key-value pairs
	fieldStr := ""
	for k, v := range fields {
		fieldStr += fmt.Sprintf("%s=%v ", k, v)
	}
	log.Printf("[DEBUG] %s %s", msg, fieldStr)
}

// LogInfoWithFields logs an info message with structured fields to both tflog and standard logger
func logInfoWithFields(ctx context.Context, msg string, fields map[string]interface{}) {
	tflog.Info(ctx, msg, fields)

	// For standard logger, format fields as key-value pairs
	fieldStr := ""
	for k, v := range fields {
		fieldStr += fmt.Sprintf("%s=%v ", k, v)
	}
	log.Printf("[INFO] %s %s", msg, fieldStr)
}

// LogErrorWithFields logs an error message with structured fields to both tflog and standard logger
func logErrorWithFields(ctx context.Context, msg string, fields map[string]interface{}) {
	tflog.Error(ctx, msg, fields)

	// For standard logger, format fields as key-value pairs
	fieldStr := ""
	for k, v := range fields {
		fieldStr += fmt.Sprintf("%s=%v ", k, v)
	}
	log.Printf("[ERROR] %s %s", msg, fieldStr)
}
