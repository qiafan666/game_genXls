package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"game_genXls/sxls"

	"github.com/tealeg/xlsx"
)

func (g *Generate) GenJson(readPath, savePath string) error {
	if readPath == "" || savePath == "" {
		return fmt.Errorf("GenStruct readPath or savePath is nil")
	}
	g.savePath = savePath
	g.typeNames = make(map[string]string, excelCount)
	files, err := os.ReadDir(readPath)
	if err != nil {
		return fmt.Errorf("GenJson ReadDir is err:%v", err)
	}

	for _, file := range files {
		if path.Ext(file.Name()) != ".xlsx" {
			continue
		}
		wb, err := xlsx.OpenFile(readPath + "/" + file.Name())
		if err != nil {
			return fmt.Errorf("GenJson file[%s] is err[%v]", file.Name(), err)
		}

		if file.Name() == userDefStructXlsName {
			continue
		}
		if wb.Sheets == nil {
			return fmt.Errorf("GenJson file[%s] has not sheet", file.Name())
		}
		for _, sheet := range wb.Sheets {
			// 默认表名
			if strings.HasPrefix(sheet.Name, "Sheet") {
				continue
			}
			// 忽略首字母非英文大写字段
			if c := sheet.Name[0]; c < 'A' || c > 'Z' {
				continue
			}
			sheetTyp, ok := sxls.AllValues[sheet.Name]
			if !ok {
				continue
			}
			data, err1 := g.ParseData(sheet, sheetTyp, file.Name())
			if err1 != nil {
				return fmt.Errorf("GenJson file[%s] is err[%v]", file.Name(), err1)
			}
			_ = g.writeJsonFile(g.savePath+"/"+sheet.Name+".json", data)
		}
	}
	return nil
}

func (g *Generate) writeJsonFile(savePath string, data []any) error {
	fw, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return fmt.Errorf("writeJsonFile OpenFile is err:%v", err)
	}
	defer fw.Close()
	_ = fw.Truncate(0)

	return json.NewEncoder(fw).Encode(data)
}
