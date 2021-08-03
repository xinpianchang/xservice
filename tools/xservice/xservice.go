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
	about = `
                          _
                         (_)
__  _____  ___ _ ____   ___  ___ ___
\ \/ / __|/ _ \ '__\ \ / / |/ __/ _ \
 >  <\__ \  __/ |   \ V /| | (_|  __/
/_/\_\___/\___|_|    \_/ |_|\___\___|
%37s

Another excellent & extensible micro service framework

xservice toolset more documentation refer
 https://github.com/xinpianchang/xservice

`

	rootCmd = &cobra.Command{
		Use:     "xservice",
		Short:   "xservice toolset",
		Long:    fmt.Sprintf(about, xservice.Version),
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
