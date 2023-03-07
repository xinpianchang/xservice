package gogen

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/xinpianchang/xservice/v2/pkg/log"
	"github.com/xinpianchang/xservice/v2/tools/xservice/gogen/assets"
)

var (
	// NewCmd is the cobra command for new project
	NewCmd = &cobra.Command{
		Use:   "new",
		Short: "create new project",
		Run: func(cmd *cobra.Command, args []string) {
			module := strings.TrimSpace(viper.GetString("module"))
			if module == "" {
				module = "github.com/example/hello"
			}

			project := newProject(module)
			if project == nil {
				return
			}

			// check target directory
			target := viper.GetString("target")

			targetFileInfo, err := os.Stat(target)
			if errors.Is(err, os.ErrNotExist) {
				if err = os.MkdirAll(target, 0755); err != nil {
					log.Fatal("create target directory", zap.Error(err))
				}
			} else if !targetFileInfo.IsDir() {
				log.Error("target should be an empty directory")
				return
			}

			if files, err := os.ReadDir(target); err != nil {
				log.Error("list target", zap.Error(err))
			} else if len(files) > 0 {
				log.Error("target should be empty")
				return
			}

			err = fs.WalkDir(assets.ProjectFS, "project", func(src string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				path := strings.TrimPrefix(src, "project")
				if path == "" {
					return nil
				}

				tf := filepath.Join(target, path)

				if d.IsDir() {
					if err = os.MkdirAll(tf, 0755); err != nil {
						return err
					}
					return nil
				}

				if strings.HasSuffix(src, ".tpl") {
					data, _ := assets.ProjectFS.ReadFile(src)
					tp := template.Must(template.New(src).Parse(string(data)))
					file, err := os.OpenFile(strings.TrimSuffix(tf, ".tpl"), os.O_CREATE|os.O_WRONLY, 0655)
					if err != nil {
						return err
					}
					if err = tp.Execute(file, project); err != nil {
						return err
					}
				} else {
					file, err := os.OpenFile(tf, os.O_CREATE|os.O_WRONLY, 0655)
					if err != nil {
						return err
					}
					sf, _ := assets.ProjectFS.Open(src)
					_, err = io.Copy(file, sf)
					if err != nil {
						return err
					}
				}

				return nil
			})

			if err != nil {
				log.Error("walk", zap.Error(err))
			}
		},
	}
)

// Project is project information container
type Project struct {
	Module string
	Repo   string
	Name   string
}

func init() {
	pf := NewCmd.PersistentFlags()
	pf.StringP("target", "t", ".", "output directory")
	pf.StringP("module", "m", "", "module name")
	_ = viper.BindPFlag("target", pf.Lookup("target"))
	_ = viper.BindPFlag("module", pf.Lookup("module"))
}

func newProject(module string) *Project {
	match, err := regexp.MatchString(`^[a-z0-9\.\-_\/]+$`, module)
	if err != nil {
		log.Error("match module", zap.Error(err))
		return nil
	}
	if !match {
		log.Error("invalid module")
		return nil
	}

	project := &Project{Module: module}
	ms := strings.Split(module, "/")
	switch n := len(ms); {
	case n == 3:
		project.Repo = ms[1]
		project.Name = ms[2]
		return project
	case n == 1:
		project.Repo = ms[0]
		project.Name = ms[0]
		return project
	default:
		log.Error("invalid module", zap.String("example", "github.com/example/helloworld"))
		return nil
	}
}
