package model

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/xinpianchang/xservice/tools/xservice/model/mysql"
)

var (
	ModelCmd = &cobra.Command{
		Use:   "model",
		Short: "generete model from datasource and basic CRUD base on GORM",
	}
)

func init() {
	ModelCmd.AddCommand(
		mysql.MySQLCmd,
	)

	pf := ModelCmd.PersistentFlags()
	pf.StringP("datasource", "d", "", "datasource, valid golang SQL DSN, e.g. root:123456@(127.0.0.1:3306)/test")
	pf.StringP("filter", "f", "", "filter table via regex")
	pf.String("dir", "internal/model", "generate go model files to dir")
	pf.String("pkg", "model", "model package name")
	viper.BindPFlag("datasource", pf.Lookup("datasource"))
	viper.BindPFlag("filter", pf.Lookup("filter"))
	viper.BindPFlag("dir", pf.Lookup("dir"))
	viper.BindPFlag("pkg", pf.Lookup("pkg"))
}
