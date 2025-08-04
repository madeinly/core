package fatal

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// once protects against double-opening the file.
var (
	logFile *os.File
	logger  *log.Logger
	inited  bool
)

// Init MUST be called once at program start-up.
// It pre-creates ./logs/error.log and keeps it open for the lifetime
// of the process.  Any problem here is a programmer error, not a
// runtime error – we panic so the program never starts.
func Init() {
	if inited {
		panic("fatal.Init called twice")
	}
	inited = true

	exe, err := os.Executable()
	if err != nil {
		panic("cannot locate executable: " + err.Error())
	}
	logDir := filepath.Join(filepath.Dir(exe), "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		panic("cannot create log dir: " + err.Error())
	}

	logPath := filepath.Join(logDir, "error.log")
	f, err := os.OpenFile(logPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		panic("cannot open error.log: " + err.Error())
	}
	logFile = f
	logger = log.New(logFile, "", log.LstdFlags)
}

// OnErr logs and exits gracefully.
// It never returns.
func OnErr(err error, msg string, args ...any) {
	if err == nil {
		return
	}

	if !inited {
		// Programmer error – we cannot even log, so panic.
		panic("fatal.OnErr called before fatal.Init")
	}

	// Build the formatted message.
	logger.Printf("FATAL: "+msg, args...)
	logger.Printf("goroutine %d\n", runtime.NumGoroutine())

	// Flush and close the log file explicitly.
	_ = logFile.Close()

	// Exit *after* every defer in the calling stack has run.
	os.Exit(1)
}
