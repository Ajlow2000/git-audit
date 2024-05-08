package git

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

// isDirectory determines if a file represented
// by `path` is a directory or not
func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}


func validateTarget(target string) error {
    if target == "" {
        userHome, err := os.UserHomeDir()
        if err != nil {
            return err
        }
        target = userHome
    } else {
        validDir, err := isDirectory(target)
        if err != nil {
            return err
        }
        if !validDir {
            return errors.New("'" + target + "' is not a valid direcotry")
        }
    }
    
    return nil
}

func validateLogPath(logPath string) (string, error) {
    if isDir, _ := isDirectory(logPath); isDir {
        standardFileName := time.Now().UTC().Format("2006-01-02")
        standardFileName = standardFileName + ".json"
        logPath = filepath.Join(logPath, standardFileName)
    } else { // assume logPath is meant to be a file
        base := filepath.Base(logPath)
        if filepath.Ext(base) != ".json" {
            return logPath, errors.New("logPath (" + logPath + ") must be of type '.json'")
        }
    }
    return logPath, nil
}

// AuditDir takes a string [target] and searches for git repos within the specified path.
// Generates and outputs a report of the status of found git repositories.
//
// By default, logs will print to stderr and lead to very cluttered output. It is suggested
// to either specify a [logPath] or throw the stderr away. Ie: `audit-dir $HOME 2> /dev/null`
func AuditDir(target string, logPath string)  {
    // Initialize Logger
    if logPath == "stderr" {
        slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))
    } else {
        logPath, err := validateLogPath(logPath)
        if err != nil {
            fmt.Println(err.Error())
            return
        }
        logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            fmt.Println("An error occured when opening '" + logPath + "' in append mode")
            fmt.Println("\t" + err.Error())
            return
        }
        defer logFile.Close()
        fmt.Println("Initialized logging to: " + logPath)
        slog.SetDefault(slog.New(slog.NewJSONHandler(logFile, nil)))
    }

    // Validate target
    if target == "" {
        var err error
        target, err = os.UserHomeDir()
        if err != nil {
            slog.Error(err.Error())
        }
    } else {
        err := validateTarget(target)
        if err != nil {
            slog.Error(err.Error())
            return
        }
    }

    // Begin Audit
    currentUser, err := user.Current()
    if err != nil {
        slog.Error(err.Error())
    }
    username := currentUser.Username
    slog.Info("Beginning an audit of git repositories", slog.Any("target", target), slog.Any("user", username), slog.Any("audit_status", "initialized"))
    fmt.Println("Beginning an audit of git repositories in: " + target)
    
    fmt.Println("Succesfully completed")
    slog.Info("Concluded audit", slog.Any("target", target), slog.Any("user", username), slog.Any("audit_status", "completed"))
}
