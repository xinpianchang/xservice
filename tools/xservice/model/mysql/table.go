package mysql

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/xinpianchang/xservice/pkg/stringx"
)

type (
	Config struct {
		Dir         string
		Pkg         string
		Filter      *regexp.Regexp
		Gormcomment bool
	}

	MySQLGenerator struct {
		Config *Config
		Tables []*Table
	}

	Table struct {
		Name      string
		Comment   string
		Fields    []*Field
		Statement *jen.Statement
	}

	Field struct {
		TableName     string
		ColumnName    string
		ColumnDefault string
		IsNullable    bool
		DataType      string
		ColumnType    string
		ColumnKey     string
		Extra         string
		ColumnComment string
		GoType        string         `gorm:"-"`
		Statement     *jen.Statement `gorm:"-"`
	}
)

var (
	typeMysqlDic = map[string]string{
		"smallint":            "int16",
		"smallint unsigned":   "uint16",
		"int":                 "int",
		"int unsigned":        "uint",
		"bigint":              "int64",
		"bigint unsigned":     "uint64",
		"varchar":             "string",
		"char":                "string",
		"date":                "time.Time",
		"datetime":            "time.Time",
		"bit":                 "int8",
		"bit(1)":              "int8",
		"tinyint":             "int8",
		"tinyint unsigned":    "uint8",
		"tinyint(1)":          "int8",
		"tinyint(1) unsigned": "uint8",
		"json":                "string",
		"text":                "string",
		"timestamp":           "time.Time",
		"double":              "float64",
		"mediumtext":          "string",
		"longtext":            "string",
		"float":               "float32",
		"tinytext":            "string",
		"enum":                "string",
		"time":                "time.Time",
		"blob":                "[]byte",
		"tinyblob":            "[]byte",
	}

	// typeMysqlMatch regexp match types
	typeMysqlMatch = [][]string{
		{`^(tinyint)[(]\d+[)] unsigned`, "uint8"},
		{`^(tinyint)[(]\d+[)]`, "int8"},
		{`^(smallint)[(]\d+[)] unsigned`, "uint16"},
		{`^(smallint)[(]\d+[)]`, "int16"},
		{`^(int)[(]\d+[)] unsigned`, "uint"},
		{`^(int)[(]\d+[)]`, "int"},
		{`^(bit)[(]\d+[)]`, "int8"},
		{`^(bigint)[(]\d+[)] unsigned`, "uint64"},
		{`^(bigint)[(]\d+[)]`, "int64"},
		{`^(char)[(]\d+[)]`, "string"},
		{`^(enum)[(](.)+[)]`, "string"},
		{`^(set)[(](.)+[)]`, "string"},
		{`^(varchar)[(]\d+[)]`, "string"},
		{`^(varbinary)[(]\d+[)]`, "[]byte"},
		{`^(binary)[(]\d+[)]`, "[]byte"},
		{`^(tinyblob)[(]\d+[)]`, "[]byte"},
		{`^(decimal)[(]\d+,\d+[)]`, "float64"},
		{`^(mediumint)[(]\d+[)]`, "string"},
		{`^(double)[(]\d+,\d+[)]`, "float64"},
		{`^(float)[(]\d+,\d+[)]`, "float64"},
		{`^(datetime)[(]\d+[)]`, "time.Time"},
		{`^(timestamp)[(]\d+[)]`, "time.Time"},
	}

	linebreak = regexp.MustCompile(`[\n\r]+`)
)

func NewMySQLGenerator(config *Config) *MySQLGenerator {
	return &MySQLGenerator{
		Config: config,
		Tables: make([]*Table, 0, 512),
	}
}

func (t *MySQLGenerator) Gen(dsn string) error {
	if err := t.parse(dsn); err != nil {
		return err
	}

	c := jen.NewFile(t.Config.Pkg)
	c.HeaderComment("auto generated file DO NOT EDIT")
	c.Line()

	for _, table := range t.Tables {
		c.Add(table.Statement).Line()
	}

	file := filepath.Join(t.Config.Dir, "model.gen.go")
	err := c.Save(file)
	if err == nil {
		fmt.Println("generage model:", file)
	}
	return err
}

func (t *MySQLGenerator) parse(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return err
	}

	var tb []map[string]interface{}
	err = db.Raw("select table_name, table_comment from information_schema.tables where table_schema = database()").Find(&tb).Error
	if err != nil {
		return err
	}

	names := make([]string, 0, len(tb))
	for _, item := range tb {
		name := item["table_name"].(string)
		comment := item["table_comment"].(string)
		table := &Table{Name: name, Comment: comment}
		if t.Config.Filter != nil {
			if t.Config.Filter.MatchString(name) {
				t.Tables = append(t.Tables, table)
				names = append(names, name)
			}
		} else {
			t.Tables = append(t.Tables, table)
			names = append(names, name)
		}
	}

	var fields = make([]*Field, 0, 1024*8)
	err = db.Raw(`
		select
			table_name,
			column_name,
			column_default,
			lower(is_nullable) = 'yes' is_nullable,
			data_type,
			lower(column_type) column_type,
			column_key,
			extra,
			column_comment
		from
			information_schema.columns
		where
			table_schema = database()
			and table_name in(?)
	`, names).Scan(&fields).Error
	if err != nil {
		return err
	}

	for _, field := range fields {
		if v, ok := typeMysqlDic[field.ColumnType]; ok {
			field.GoType = v
		} else {
			for _, v := range typeMysqlMatch {
				if ok, _ := regexp.MatchString(v[0], field.ColumnType); ok {
					field.GoType = v[1]
					break
				}
			}
		}

		if field.GoType == "" {
			panic(fmt.Sprintf("unknown type: %s", field.ColumnType))
		}

		field.Statement = t.fieldStatement(field)
	}

	for _, table := range t.Tables {
		for _, field := range fields {
			if strings.EqualFold(table.Name, field.TableName) {
				if table.Fields == nil {
					table.Fields = make([]*Field, 0, 100)
				}
				table.Fields = append(table.Fields, field)
			}
		}

		table.Statement = t.tableStatement(table)
	}

	return nil
}

func (t *MySQLGenerator) tableStatement(table *Table) *jen.Statement {
	typeName := stringx.CamelCase(table.Name)
	typeName = strings.ReplaceAll(typeName, "-", "")

	c := jen.Comment(fmt.Sprintf("%s table: %s", typeName, table.Name)).Line()
	if table.Comment != "" {
		c.Comment(table.Comment).Line()
	}

	c.Type().Id(typeName).Struct(jen.Do(func(c *jen.Statement) {
		for _, field := range table.Fields {
			c.Add(field.Statement).Line()
		}
		t.removeStatement(c, 1)
	})).Line()

	// table name
	c.Commentf("TableName set table of %v, refer: https://gorm.io/docs/conventions.html", table.Name).Line()
	c.Func().Params(jen.Id(typeName)).Id("TableName").Call().String().Block(
		jen.Return(jen.Lit(fmt.Sprint(table.Name))),
	).Line().Line()

	// BeforeUpdate Hook
	fields := make([]*Field, 0, 2)
	for _, field := range table.Fields {
		if field.GoType == "time.Time" {
			for _, it := range []string{"updated_at", "last_updated_at", "last_changed_at"} {
				if strings.EqualFold(field.ColumnName, it) {
					fields = append(fields, field)
					break
				}
			}
		}
	}
	if len(fields) > 0 {
		c.Comment("BeforeUpdate use for check field has changed, refer: https://gorm.io/docs/update.html").Line()
		c.Comment("The Changed method only works with methods Update, Updates,").Line()
		c.Comment("and it only checks if the updating value from Update / Updates equals the model value, will return true if it is changed and not omitted").Line()
		c.Func().Params(jen.Id("t").Id(typeName)).Id("BeforeUpdate").Params(jen.Id("tx").Op("*").Id("gorm.DB")).Id("error").Block(
			jen.If(jen.Id("tx.Statement.Changed").Call()).Block(
				jen.Do(func(c *jen.Statement) {
					for _, field := range fields {
						c.Id("tx.Statement.SetColumn").Call(jen.Lit(stringx.CamelCase(field.ColumnName)), jen.Id("time.Now").Call()).Line()
					}
					t.removeStatement(c, 1)
				}),
			),
			jen.Return(jen.Nil()),
		)
		c.Line().Line()
	}

	// model basic
	c.Add(t.tableDefaultModel(table))
	c.Line()

	return c
}

func (t *MySQLGenerator) tableDefaultModel(table *Table) *jen.Statement {
	typeName := stringx.LowerCamelCase(fmt.Sprint("default_", table.Name, "Model"))
	typeName = strings.ReplaceAll(typeName, "-", "")
	modelName := strings.ReplaceAll(stringx.CamelCase(table.Name), "-", "")

	c := jen.Commentf("%s default %sModel implements with basic operation", typeName, stringx.CamelCase(table.Name)).Line()
	c.Type().Id(typeName).Struct(
		jen.Id("tx").Op("*").Qual("gorm.io/gorm", "DB"),
	).Line()

	// new model
	newModelFn := stringx.CamelCase(fmt.Sprint("New", modelName, "Model"))
	c.Commentf("%s create new op Model", newModelFn).Line()
	c.Func().Id(newModelFn).Params(jen.Id("tx *gorm.DB")).Op("*").Id(typeName).Block(
		jen.Return(jen.Op("&").Id(typeName).Block(jen.Id("tx").Op(":").Id("tx.Model").Call(jen.Op("&").Id(modelName).Block()).Op(","))),
	).Line()

	// Model
	c.Comment("Model for update tx model").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Model").Params(jen.Id("value interface{}")).Op("*").Id(typeName).Block(
		jen.Return(jen.Id(newModelFn).Call(jen.Id("t.tx.Model(value)"))),
	).Line()

	// Scopes
	c.Comment("Scopes pass current database connection to arguments `func(DB) DB`, which could be used to add conditions dynamically").Line()
	c.Comment("refer: https://gorm.io/docs/advanced_query.html#Scopes").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Scopes").Params(jen.Id("fn ...func(tx *gorm.DB) *gorm.DB")).Op("*").Id(typeName).Block(
		jen.Return(jen.Id(newModelFn).Call(jen.Id("t.tx.Scopes(fn...)"))),
	).Line()

	// Clauses
	c.Comment("Clauses add clauses").Line()
	c.Comment("refer: https://gorm.io/docs/sql_builder.html#Clauses").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Clauses").Params(jen.Id("conds").Op("...").Qual("gorm.io/gorm/clause", "Expression")).Op("*").Id(typeName).Block(
		jen.Return(jen.Id(newModelFn).Call(jen.Id("t.tx.Clauses(conds...)"))),
	).Line()

	// Where
	c.Comment("Where").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Where").Params(jen.Id("query interface{}, args ...interface{}")).Op("*").Id(typeName).Block(
		jen.Return(jen.Id(newModelFn).Call(jen.Id("t.tx.Where(query, args...)"))),
	).Line()

	// Order
	c.Comment("Order").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Order").Params(jen.Id("value interface{}")).Op("*").Id(typeName).Block(
		jen.Return(jen.Id(newModelFn).Call(jen.Id("t.tx.Order(value)"))),
	).Line()

	// Distinct
	c.Comment("Distinct").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Distinct").Params(jen.Id("args interface{}")).Op("*").Id(typeName).Block(
		jen.Return(jen.Id(newModelFn).Call(jen.Id("t.tx.Distinct(args)"))),
	).Line()

	// Limit
	c.Comment("Limit").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Limit").Params(jen.Id("limit int")).Op("*").Id(typeName).Block(
		jen.Return(jen.Id(newModelFn).Call(jen.Id("t.tx.Limit(limit)"))),
	).Line()

	// Offset
	c.Comment("Offset").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Offset").Params(jen.Id("offset int")).Op("*").Id(typeName).Block(
		jen.Return(jen.Id(newModelFn).Call(jen.Id("t.tx.Offset(offset)"))),
	).Line()

	// Select
	c.Comment("Select").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Select").Params(jen.Id("query interface{}, args ...interface{}")).Op("*").Id(typeName).Block(
		jen.Return(jen.Id(newModelFn).Call(jen.Id("t.tx.Select(query, args...)"))),
	).Line()

	// Scan
	c.Comment("Scan scan value to a struct").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Scan").Params(jen.Id("dest interface{}")).Id("error").Block(
		jen.Return(jen.Id("t.tx.Scan(dest).Error")),
	).Line()

	// Pluck
	c.Comment("Pluck Query single column from database and scan into a slice").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Pluck").Params(jen.Id("column string, dest interface{}")).Id("error").Block(
		jen.Return(jen.Id("t.tx.Pluck(column, dest).Error")),
	).Line()

	// Count
	c.Comment("Count Get matched records count").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Count").Params().Id("(int64, error)").Block(
		jen.Id("var v int64"),
		jen.Id("err := t.tx.Count(&v).Error"),
		jen.Return(jen.Id("v, err")),
	).Line()

	// CountMust
	c.Comment("CountMust Get matched records count if error occurs just log it").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("CountMust").Params(jen.Id("ctx").Qual("context", "Context")).Id("int64").Block(
		jen.Id("var v int64"),
		jen.Id("err := t.tx.Count(&v).Error"),
		jen.If(jen.Id("err != nil")).Block(
			jen.Qual("github.com/xinpianchang/xservice/pkg/log", "For").Id("(ctx).Error").Call(jen.Lit("count"), jen.Qual("go.uber.org/zap", "Error").Call(jen.Id("err"))),
		),
		jen.Return(jen.Id("v")),
	).Line()

	// Create
	c.Comment("Create insert the value into database").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Create").Params(jen.Id("data *").Id(modelName)).Id("error").Block(
		jen.Return(jen.Id("t.tx.Create(data).Error")),
	).Line()

	// CreateAll
	c.Comment("CreateAll Batch Insert").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("CreateAll").Params(jen.Id("data []*").Id(modelName)).Id("error").Block(
		jen.Return(jen.Id("t.tx.Create(data).Error")),
	).Line()

	// CreateInBatches
	c.Comment("CreateInBatches").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("CreateInBatches").Params(jen.Id("data []*").Id(modelName).Id(", batchSize int")).Id("error").Block(
		jen.Return(jen.Id("t.tx.CreateInBatches(data, batchSize).Error")),
	).Line()

	// CreateInBatchesMap
	c.Comment("CreateInBatchesMap").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("CreateInBatchesMap").Params(jen.Id("data []map[string]interface{}, batchSize int")).Id("error").Block(
		jen.Return(jen.Id("t.tx.CreateInBatches(data, batchSize).Error")),
	).Line()

	// CreateMap
	c.Comment("CreateMap Create From Map").Line()
	c.Comment("common usage eg. Create From SQL Expression/Context Valuer refer https://gorm.io/docs/create.html#Create-From-SQL-Expression-Context-Valuer").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("CreateMap").Params(jen.Id("data []map[string]interface{}")).Id("error").Block(
		jen.Return(jen.Id("t.tx.Create(data).Error")),
	).Line()

	// Save
	c.Comment("Save update value in database, if the value doesn't have primary key, will insert it").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Save").Params(jen.Id("data *").Id(modelName)).Id("error").Block(
		jen.Return(jen.Id("t.tx.Save(data).Error")),
	).Line()

	// SaveAll
	c.Comment("SaveAll update all items in database, if the value doesn't have primary key, will insert it").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("SaveAll").Params(jen.Id("data []*").Id(modelName)).Id("error").Block(
		jen.Return(jen.Id("t.tx.Save(data).Error")),
	).Line()

	// Update
	c.Comment("Update attributes with `struct`, will only update non-zero fields").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Update").Params(jen.Id("data *").Id(modelName)).Id("error").Block(
		jen.Return(jen.Id("t.tx.Updates(data).Error")),
	).Line()

	// UpdateForce
	c.Comment("UpdateForce force update include zero value fields").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("UpdateForce").Params(jen.Id("data *").Id(modelName)).Id("error").Block(
		jen.Return(jen.Id(`t.tx.Model(data).Select("*").Updates(data).Error`)),
	).Line()

	// UpdateMap
	c.Comment("UpdateMap update attributes with `map`").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("UpdateMap").Params(jen.Id("data map[string]interface{}")).Id("error").Block(
		jen.Return(jen.Id("t.tx.Updates(data).Error")),
	).Line()

	// UpdateColumn
	c.Comment("UpdateColumn update only one column").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("UpdateColumn").Params(jen.Id("column string, value interface{}")).Id("error").Block(
		jen.Return(jen.Id("t.tx.UpdateColumn(column, value).Error")),
	).Line()

	// Get
	c.Comment("Get retrieving by primary key").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Get").Params(jen.Id("id interface{}")).Op("(*").Id(modelName).Id(",error").Op(")").Block(
		jen.Id("var data").Id(modelName),
		jen.Id("err := t.tx.Find(&data, id).Error"),
		jen.Return(jen.Id("&data, err")),
	).Line()

	c.Comment("Find").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Find").Params(jen.Id("conds ...interface{}")).Op("([]*").Id(modelName).Id(",error").Op(")").Block(
		jen.Id("var data []*").Id(modelName),
		jen.Id("err := t.tx.Find(&data, conds...).Error"),
		jen.Return(jen.Id("data, err")),
	).Line()

	c.Comment("First").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("First").Params(jen.Id("conds ...interface{}")).Op("(*").Id(modelName).Id(",error").Op(")").Block(
		jen.Id("var data ").Id(modelName),
		jen.Id("err := t.tx.First(&data, conds...).Error"),
		jen.Return(jen.Id("&data, err")),
	).Line()

	c.Comment("Last").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Last").Params(jen.Id("conds ...interface{}")).Op("(*").Id(modelName).Id(",error").Op(")").Block(
		jen.Id("var data ").Id(modelName),
		jen.Id("err := t.tx.Last(&data, conds...).Error"),
		jen.Return(jen.Id("&data, err")),
	).Line()

	c.Comment("Take").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Take").Params(jen.Id("conds ...interface{}")).Op("(*").Id(modelName).Id(",error").Op(")").Block(
		jen.Id("var data ").Id(modelName),
		jen.Id("err := t.tx.Take(&data, conds...).Error"),
		jen.Return(jen.Id("&data, err")),
	).Line()

	c.Comment("FindOne, avoid 'record not found' error").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("FindOne").Params(jen.Id("conds ...interface{}")).Op("(*").Id(modelName).Id(",error").Op(")").Block(
		jen.Id("var data ").Id(modelName),
		jen.Id("err := t.tx.First(&data, conds...).Error"),
		jen.If(jen.Id("err != nil &&").Qual("errors", "Is(err, gorm.ErrRecordNotFound)")).Block(
			jen.Return(jen.Id("nil, nil")),
		),
		jen.Return(jen.Id("&data, err")),
	).Line()

	c.Comment("Delete").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("Delete").Params(jen.Id("conds ...interface{}")).Id("error").Block(
		jen.Return(jen.Id("t.tx.Delete(&").Id(modelName).Id("{}, conds...).Error")),
	).Line()

	c.Comment("DeletePermanently").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("DeletePermanently").Params(jen.Id("conds ...interface{}")).Id("error").Block(
		jen.Return(jen.Id("t.tx.Unscoped().Delete(&").Id(modelName).Id("{}, conds...).Error")),
	).Line()

	fields := make([]string, 0, len(table.Fields))
	for _, field := range table.Fields {
		if field.isPrimary() {
			continue
		}
		fields = append(fields, fmt.Sprintf(`"%s"`, field.ColumnName))
	}
	c.Comment("GetFieldNames get all field names but exclude primary key field").Line()
	c.Func().Params(jen.Id("t").Op("*").Id(typeName)).Id("GetFieldNames").Params().Id("[]string").Block(
		jen.Return(jen.Id("[]string{" + strings.Join(fields, ", ") + "}")),
	).Line()

	return c
}

func (t *MySQLGenerator) fieldTypeStatement(field *Field) *jen.Statement {
	comment := t.oneline(field.ColumnComment)
	idx := strings.Index(comment, "struct:")
	if idx >= 0 {
		idx := strings.Index(comment, "struct:")
		gs := comment[idx+7:]
		idx = strings.Index(gs, " ")
		if idx > 0 {
			gs = gs[:idx]
		}
		idx = strings.LastIndex(gs, "/")
		path := ""
		name := gs
		if idx > 0 {
			path = gs[:idx]
			name = gs[idx+1:]
		}
		return jen.Qual(path, name)
	}

	switch field.GoType {
	case "int":
		return jen.Int()
	case "uint":
		return jen.Uint()
	case "int8":
		return jen.Int8()
	case "uint8":
		return jen.Uint8()
	case "int16":
		return jen.Int16()
	case "uint16":
		return jen.Uint16()
	case "int32":
		return jen.Int32()
	case "uint32":
		return jen.Uint32()
	case "int64":
		return jen.Int64()
	case "uint64":
		return jen.Uint64()
	case "string":
		return jen.String()
	case "time.Time":
		if strings.EqualFold(field.ColumnName, "deleted_at") {
			return jen.Qual("gorm.io/gorm", "DeletedAt")
		} else {
			return jen.Qual("time", "Time")
		}
	case "float32":
		return jen.Float32()
	case "float64":
		return jen.Float64()
	case "[]byte":
		return jen.Op("[]").Byte()
	default:
		panic(fmt.Sprintf("unknown GoType %s", field.GoType))
	}
}

func (t *MySQLGenerator) fieldStatement(field *Field) *jen.Statement {
	c := jen.Id(stringx.CamelCase(field.ColumnName))

	if field.IsNullable {
		c.Op("*")
	}

	c.Add(t.fieldTypeStatement(field))

	// tags
	tag := make(map[string]string, 2)
	{
		v := fmt.Sprintf(`column:%s;type:%s`, field.ColumnName, field.ColumnType)
		if field.isPrimary() {
			v += ";primary_key"
		}
		if strings.EqualFold(field.Extra, "auto_increment") {
			v += fmt.Sprint(";autoIncrement", field.ColumnDefault)
		}
		if !field.IsNullable {
			v += ";not null"
		}
		if field.ColumnDefault != "" {
			v += fmt.Sprint(";default:", field.ColumnDefault)
		}
		if t.Config.Gormcomment && field.ColumnComment != "" {
			comment := t.oneline(field.ColumnComment)
			comment = strings.ReplaceAll(comment, ";", `\\;`)
			comment = strings.ReplaceAll(comment, ":", `\\:`)
			v += fmt.Sprint(";comment:", comment)
		}
		tag["gorm"] = v
	}

	{
		if !strings.Contains(field.ColumnComment, "json_hidden") {
			v := stringx.LowerCamelCase(field.ColumnName)
			if field.IsNullable {
				v += ",omitempty"
			}
			tag["json"] = v
		}
	}

	c.Tag(tag)

	if comment := t.oneline(field.ColumnComment); comment != "" {
		c.Comment(comment)
	}

	return c
}

func (t *MySQLGenerator) removeStatement(c *jen.Statement, count int) {
	*c = (*c)[:len(*c)-count]
}

func (t *MySQLGenerator) oneline(str string) string {
	return strings.TrimSpace(linebreak.ReplaceAllString(str, " "))
}

func (t Field) isPrimary() bool {
	return strings.ToUpper(t.ColumnKey) == "PRI"
}
