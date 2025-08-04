package internal

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/madeinly/core/settings"
)

// RunChecks verifies all critical dependencies and permissions at startup.
// It will panic if any check fails, preventing the application from starting
// in a broken state.
func RunChecks(ch chan<- string) {
	ch <- "Running pre-flight checks..."

	if err := checkLogFile(); err != nil {
		ch <- fmt.Sprintf("Pre-flight check failed: %v", err)
		panic(0)
	}

	if err := checkEmail(); err != nil {
		ch <- fmt.Sprintf("Pre-flight check failed: %v", err)
		panic(0)
	}

	if err := checkHTTP(); err != nil {
		ch <- fmt.Sprintf("Pre-flight check failed: %v", err)
		panic(0)
	}

	if err := checkWorkspacePermissions(); err != nil {
		ch <- fmt.Sprintf("Pre-flight check failed: %v", err)
		panic(0)
	}

	ch <- "All pre-flight checks passed."
}

// checkLogFile ensures the log directory and file are writable.
func checkLogFile() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot locate executable: %w", err)
	}
	logDir := filepath.Join(filepath.Dir(exe), "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return fmt.Errorf("cannot create log directory: %w", err)
	}

	logPath := filepath.Join(logDir, "error.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("cannot open error.log for writing: %w", err)
	}
	return f.Close()
}

// checkEmail verifies the SMTP connection and credentials.
func checkEmail() error {
	appSettings := settings.GetSettings().Email
	if appSettings.Address == "" || appSettings.Port == "" || appSettings.User == "" {
		return fmt.Errorf("email settings are incomplete")
	}

	serverAddr := fmt.Sprintf("%s:%s", appSettings.Address, appSettings.Port)
	auth := smtp.PlainAuth("", appSettings.User, appSettings.Password, appSettings.Address)

	if appSettings.Encryption == "ssl/tls" {
		tlsconfig := &tls.Config{ServerName: appSettings.Address}
		conn, err := tls.Dial("tcp", serverAddr, tlsconfig)
		if err != nil {
			return fmt.Errorf("email (ssl/tls) connection failed: %w", err)
		}
		client, err := smtp.NewClient(conn, appSettings.Address)
		if err != nil {
			return err
		}
		defer client.Close()
		return client.Noop()
	} else {
		client, err := smtp.Dial(serverAddr)
		if err != nil {
			return fmt.Errorf("email (starttls) connection failed: %w", err)
		}
		defer client.Close()
		if ok, _ := client.Extension("STARTTLS"); ok {
			tlsconfig := &tls.Config{ServerName: appSettings.Address}
			if err = client.StartTLS(tlsconfig); err != nil {
				return fmt.Errorf("email (starttls) failed to upgrade: %w", err)
			}
		}
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("email auth failed: %w", err)
		}
		return client.Noop()
	}
}

// checkHTTP verifies that outbound HTTP requests can be made.
func checkHTTP() error {
	// We check against a known, reliable server.
	_, err := http.Get("https://www.google.com/generate_204")
	if err != nil {
		return fmt.Errorf("outbound http check failed: %w", err)
	}
	return nil
}

// checkWorkspacePermissions verifies that the application has the correct permissions.
func checkWorkspacePermissions() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot locate executable: %w", err)
	}
	// Check if the executable's directory is writable.
	if err := os.WriteFile(filepath.Join(filepath.Dir(exe), ".perm_check"), []byte("test"), 0o644); err != nil {
		return fmt.Errorf("workspace permission check failed: %w", err)
	}
	return os.Remove(filepath.Join(filepath.Dir(exe), ".perm_check"))
}
