package common

import (
    "fmt"
)

func LogEntryf(entity string, format string, v ...interface{}) string {
    logFormat := fmt.Sprintf("[%s] - %s", entity, format)
    return fmt.Sprintf(logFormat, v...)
}

func LogEntry(entity string, msg string) string {
    return fmt.Sprintf("[%s] - %s", entity, msg)
}
