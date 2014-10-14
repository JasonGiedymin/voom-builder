package common

import (
    "fmt"
)

func LogEntryf(entity string, format string, v ...interface{}) string {
    logFormat := fmt.Sprintf("[%s] - %s", entity, format)
    return fmt.Sprintf(logFormat, v...)
}
