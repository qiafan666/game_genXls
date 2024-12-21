package lib

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

type tsvFieldParser struct { //
	parse func(string, reflect.Value) error //
}

// go 代码支持的类型
var goTypes = map[reflect.Kind]tsvFieldParser{
	reflect.Int8:    {parse: xlsParseInt},
	reflect.Int16:   {parse: xlsParseInt},
	reflect.Int32:   {parse: xlsParseInt},
	reflect.Int64:   {parse: xlsParseInt},
	reflect.Struct:  {parse: xlsParseJson},
	reflect.Map:     {parse: xlsParseJson},
	reflect.Bool:    {parse: xlsParseBool},
	reflect.Array:   {parse: xlsParseJson},
	reflect.Slice:   {parse: xlsParseJson},
	reflect.Float32: {parse: xlsParseFloat},
	reflect.Float64: {parse: xlsParseFloat},
	reflect.String:  {parse: xlsParseString},
}

type goFieldMeta struct {
	name      string       // 全小写
	fieldIdx  int          // struct field index
	columnIdx int          // xls column index
	goTyp     reflect.Type //
}

type conf struct {
	name   string                  // 配置表名
	typ    reflect.Type            // go struct type
	fields map[string]*goFieldMeta // go struct field信息
}

func newConf(typ reflect.Type, name string) *conf {
	return &conf{
		name:   name,
		typ:    typ,
		fields: map[string]*goFieldMeta{},
	}
}

// 解析go struct field
func (c *conf) parseGoStructField() error {
	for i := 0; i < c.typ.NumField(); i++ {
		f := c.typ.Field(i)
		// 忽略非大写字段
		if c := f.Name[0]; c < 'A' || c > 'Z' {
			continue
		}
		ft := f.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		// 字段类型检查
		if _, ok := goTypes[ft.Kind()]; !ok {
			return fmt.Errorf("sheet=%v.%v 不支持的类型=%v", c.name, f.Name, ft.Kind())
		}
		c.fields[f.Name] = &goFieldMeta{
			fieldIdx:  i,
			columnIdx: -1,
			name:      f.Name,
			goTyp:     f.Type,
		}
	}

	return nil
}

// 解析tsv行首信息
func (c *conf) parseColumnMeta(conditions, names []string) error {
	for i, cond := range conditions {

		//id不区分类型，同时保证整洁，统一为Id
		if capitalizeFirstLetter(names[i]) == "Id" {
			name := capitalizeFirstLetter(names[i])
			if f := c.fields[name]; f != nil {
				f.columnIdx = i
			}
			continue
		}

		// 服务器不解析
		if cond != "S" && cond != "A" {
			continue
		}

		name := capitalizeFirstLetter(names[i])
		if f := c.fields[name]; f != nil {
			f.columnIdx = i
		}
	}
	return nil
}

// 解析配置数据
func (c *conf) parseConfData(lines [][]string) ([]any, error) {
	if len(lines) == 0 {
		return nil, nil
	}
	ret := make([]any, 0, len(lines)-1)
	for idx, line := range lines {
		v, err := c.parseLine(line, idx)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}

func (c *conf) parseLine(line []string, lineNum int) (any, error) {
	rv := reflect.New(c.typ).Elem()
	rt := rv.Type()
	for _, f := range c.fields {
		if f.columnIdx == -1 {
			return nil, fmt.Errorf("conf %v parse %v typ:%v 表头格式错误（一般是空的列，有空格，小写等）", c.name, f.name, c.typ)
		}
		ft := rt.Field(f.fieldIdx)
		fv := rv.Field(f.fieldIdx)
		fk := ft.Type.Kind()
		row := line[f.columnIdx]
		err := goTypes[fk].parse(row, fv)
		if err != nil {
			return nil, fmt.Errorf("sheet=%v 解析字段=%v 错误=%v 行数=%v 类型=%v", c.name, f.name, err, lineNum+5, c.typ)
		}
	}
	return rv.Addr().Interface(), nil
}

func xlsParseInt(s string, v reflect.Value) error {
	if strings.TrimLeft(s, " ") == "" {
		v.SetInt(0)
		return nil
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("%v 必须为 int", s)
	}

	switch v.Kind() {
	case reflect.Int8:
		if i > math.MaxInt8 {
			return fmt.Errorf("%v overflow", s)
		}
	case reflect.Int16:
		if i > math.MaxInt16 {
			return fmt.Errorf("%v overflow", s)
		}
	case reflect.Int32:
		if i > math.MaxInt32 {
			return fmt.Errorf("%v overflow", s)
		}
	default:
	}

	v.SetInt(i)
	return nil
}

func xlsParseString(s string, v reflect.Value) error {
	v.SetString(s)
	return nil
}

func xlsParseBool(s string, v reflect.Value) error {
	if s == "" {
		v.SetBool(false)
		return nil
	}
	i, err := strconv.ParseBool(s)
	if err != nil {
		return fmt.Errorf("%v 必须为 bool", s)
	}
	v.SetBool(i)
	return nil
}

func xlsParseFloat(s string, v reflect.Value) error {
	if s == "" {
		v.SetFloat(0)
		return nil
	}
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("%v 必须为 float64", s)
	}
	v.SetFloat(i)
	return nil
}

func xlsParseJson(s string, v reflect.Value) error {
	if s == "" {
		return nil
	}
	// v 是指针, 创建实体
	if v.Type().Kind() == reflect.Ptr {
		e := reflect.New(v.Type().Elem())
		v.Set(e)
		v = v.Elem()
	}
	// json 反序列化
	return json.Unmarshal([]byte(s), v.Addr().Interface())
}
