package common

import (
	"os"
	"syscall"
)

// ShutdownSignals returns all the signals that are being watched for to shut down services.
func ShutdownSignals() []os.Signal {
	return []os.Signal{
		syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL,
	}
}
