package git

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"time"
)

var (
    IGNORE []string
    GIT_REPOS []string
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

// Expects an absolute path to a directory and checks if it contains
// a .git dir indicating that this path is a git repository.
func isGitRepo(dirPath string) (bool, error) {
    entries, err := os.ReadDir(dirPath)
    if err != nil {
        return false, err
    }

    for _, entry := range entries {
        if entry.IsDir() && entry.Name() == ".git" {
            return true, nil 
        }
    }
    return false, nil 
}

func recurseInSearchOfGit(dirPath string) (error) {
    entries, err := os.ReadDir(dirPath)
    if err != nil {
        return err
    }

    path, err := filepath.Abs(dirPath)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        if entry.IsDir() {
            entryPath := filepath.Join(path, entry.Name())
            isGitRepo, err := isGitRepo(entryPath)
            if err != nil {
                return err
            }
            if isGitRepo {
                GIT_REPOS = append(GIT_REPOS, entryPath)
            } else {
                if !slices.Contains(IGNORE, dirPath) {
                    recurseInSearchOfGit(entryPath)
                }
            }
        }
    }
    return nil
}

// AuditDir takes a string [target] and searches for git repos within the specified path.
// Generates and outputs a report of the status of found git repositories.
func AuditDir(target string, logPath string, ignore []string)  {
    // Initialize Logger
    if logPath == "" {
        slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
    } else if logPath == "stderr" {
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

    // Set globals
    IGNORE = ignore

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

    path, err := filepath.Abs(target)
    if err != nil {
        slog.Error(err.Error())
        return
    }
    if !slices.Contains(IGNORE, target) {
        entries, err := os.ReadDir(target)
        if err != nil {
            slog.Error(err.Error())
            return
        }
        for _, entry := range entries {
            if entry.IsDir() {
                dirPath := filepath.Join(path, entry.Name())
                isGitRepo, err := isGitRepo(dirPath)
                if err != nil {
                    slog.Error(err.Error())
                    return
                }
                if isGitRepo {
                    GIT_REPOS = append(GIT_REPOS, dirPath)
                } else {
                    // recurse
                    if !slices.Contains(IGNORE, dirPath) {
                        recurseInSearchOfGit(dirPath)
                    }
                }
            }
        }
    }

    fmt.Println("Found the following repos:")
    for _, element := range GIT_REPOS {
        fmt.Println(element)
    }

    
    fmt.Println("Succesfully completed")
    slog.Info("Concluded audit", slog.Any("target", target), slog.Any("user", username), slog.Any("audit_status", "completed"))
}

