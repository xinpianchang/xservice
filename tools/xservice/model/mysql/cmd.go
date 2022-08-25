package mysql

import (
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/xinpianchang/xservice/pkg/log"
)

var (
	MySQLCmd = &cobra.Command{
		Use:                   "mysql",
		Short:                 "generete model from mysql datasource and basic CRUD base on GORM",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			datasource := viper.GetString("datasource")
			filter := viper.GetString("filter")
			dir := viper.GetString("dir")
			pkg := viper.GetString("pkg")
			gormcomment := viper.GetBool("gormcomment")

			x := regexp.MustCompile(filter)

			config := &Config{
				Dir:         dir,
				Pkg:         pkg,
				Filter:      x,
				Gormcomment: gormcomment,
			}

			err := NewMySQLGenerator(config).Gen(datasource)
			if err != nil {
				log.Error("generate error", zap.Error(err))
			}
		},
	}
)

func init() {
	pf := MySQLCmd.PersistentFlags()
	pf.Bool("gormcomment", false, "enable gorm comment")
	_ = viper.BindPFlag("gormcomment", pf.Lookup("gormcomment"))
}
