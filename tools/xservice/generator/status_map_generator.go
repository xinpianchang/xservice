package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	StatusMapGeneratorCmd = &cobra.Command{
		Use:                   "statusmap status.yaml",
		Short:                 "generate status from yaml status mapping file",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			generateStatusMap(args[0])
		},
	}
)

func generateStatusMap(file string) {
	b, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	var statusMap map[int]string
	err = yaml.Unmarshal(b, &statusMap)
	if err != nil {
		panic(err)
	}

	var codes []int
	for k := range statusMap {
		codes = append(codes, k)
	}
	sort.Ints(codes)

	f := jen.NewFile("dto")
	f.HeaderComment("auto generated file DO NOT EDIT")
	f.HeaderComment(fmt.Sprintf("generate from file: %s ", file))
	dict := jen.Dict{}
	consts := make([]jen.Code, 0, len(codes))
	for i, code := range codes {
		codeN := fmt.Sprint(code)
		if code < 0 {
			codeN = fmt.Sprint("_", -code)
		}

		name := fmt.Sprint("StatusCode", codeN)
		if i == 0 {
			consts = append(consts, jen.Id(name).Qual("github.com/xinpianchang/xservice/pkg/responsex", "StatusCode").Op("=").Lit(code))
		} else {
			consts = append(consts, jen.Id(name).Op("=").Lit(code).Comment(statusMap[code]))
		}
		dict[jen.Id(name)] = jen.Lit(statusMap[code])
	}
	f.Add(jen.Const().Defs(consts...))
	f.Func().Id("init").Params().Block(
		jen.Comment("init status map"),
		jen.Id("responsex.SetStatusMap").Call(jen.Map(jen.Id("responsex.StatusCode")).String().Values(dict)),
	)
	fileabs, _ := filepath.Abs(file)
	target := filepath.Join(filepath.Dir(fileabs), "d_status_map.go")
	err = os.WriteFile(target, []byte(f.GoString()), 0600)
	if err != nil {
		panic(err)
	}

	target, _ = filepath.Abs(target)
	fmt.Println("generage statusmap:", target)
}
