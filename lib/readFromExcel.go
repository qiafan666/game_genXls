package lib

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"unsafe"

	"github.com/tealeg/xlsx"
)

// ReadFromExcel 从excel读取配置数据
func (g *Generate) ReadFromExcel(readPath, fileName string, typ any) ([]any, error) {
	wb, err := xlsx.OpenFile(readPath + "/" + fileName)
	if err != nil {
		return nil, fmt.Errorf("GenStruct xlsx OpenFile is err :%v", err)
	}
	return g.ParseData(wb.Sheets[0], typ, fileName)
}

func (g *Generate) ParseData(sheet *xlsx.Sheet, typ any, fileName string) ([]any, error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 2048)
			l := runtime.Stack(buf, false)
			for i := range buf[:l] {
				if buf[i] == '\n' {
					buf[i] = ' '
				}
			}
			stack := *(*string)(unsafe.Pointer(&buf))
			_, _ = fmt.Fprintf(os.Stderr, "%v: %s\n", r, stack)
			fmt.Printf("解析数据错误，请检查excel文件:%s\n", fileName)
			os.Exit(2)
		}
	}()
	conf := newConf(reflect.TypeOf(typ).Elem(), sheet.Name)
	if err := conf.parseGoStructField(); err != nil {
		return nil, err
	}
	lines := g.readAll(sheet)

	if err := conf.parseColumnMeta(lines[2], lines[1]); err != nil {
		return nil, err
	}
	return conf.parseConfData(lines[lineNumber:])
}

func (g *Generate) readAll(sheet *xlsx.Sheet) [][]string {
	sheetData := make([][]string, 0)
	// 遍历行
	for i := 0; i < sheet.MaxRow; i++ {
		// 判断某一列的第一列是否为空
		if sheet.Cell(i, 0).Value == "" {
			continue
		}
		cellData := make([]string, 0)
		// 遍历列
		for j := 0; j < sheet.MaxCol; j++ {
			cellData = append(cellData, sheet.Cell(i, j).Value)
		}
		sheetData = append(sheetData, cellData)
	}
	return sheetData
}
