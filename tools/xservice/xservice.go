package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/xinpianchang/xservice"
	"github.com/xinpianchang/xservice/tools/xservice/generator"
	"github.com/xinpianchang/xservice/tools/xservice/gogen"
	"github.com/xinpianchang/xservice/tools/xservice/model"
)

var (
	rootCmd = &cobra.Command{
		Use:     "xservice",
		Short:   "xservice toolset",
		Version: xservice.Version,
	}
)

func main() {
	rootCmd.AddCommand(
		gogen.NewCmd,
		model.ModelCmd,
		generator.StatusMapGeneratorCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
