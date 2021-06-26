// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

var app = "{{.Name}}"

// for local machine build
func Build() error {
	dir := fmt.Sprintf("dist/%s-%s-%s", app, runtime.GOOS, runtime.GOARCH)
	target := fmt.Sprintf("%s/%s", dir, app)

	args := make([]string, 0, 10)
	args = append(args, "build", "-o", target)
	args = append(args, "-ldflags="+flags(), "cmd/main.go")

	if err := sh.Run(mg.GoCmd(), args...); err != nil {
		return err
	}

	_ = os.Mkdir(fmt.Sprintf("%s/logs", dir), 0750)
	_ = sh.Copy(fmt.Sprintf("%s/README.md", dir), "README.md")
	_ = sh.Run("tar", "-czf", fmt.Sprintf("%s.tar.gz", dir), "-C", "dist", filepath.Base(dir))

	return nil
}

func flags() string {
	timestamp := time.Now().Format(time.RFC3339)
	h := hash()
	m := "{{.Module}}"
	tpl := fmt.Sprintf(`-buildid %%s -extldflags "-static" -X "%s/version.Build=%%s" -X "%s/version.BuildAt=%%s"`, m, m)
	return fmt.Sprintf(tpl, h, h, timestamp)
}

// hash returns the git hash for the current repo or "" if none.
func hash() string {
	hash, _ := sh.Output("git", "rev-parse", "HEAD")
	if hash == "" {
		return "00000000"
	}
	return hash
}

// cleanup all build files
func Clean() {
	_ = sh.Rm("dist")
}
