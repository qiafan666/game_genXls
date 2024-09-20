package lib

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"game_genXls/sxls"

	"github.com/tealeg/xlsx"
	"go.mongodb.org/mongo-driver/mongo"
)

func (g *Generate) WriteMongo(client *mongo.Client, readPath, savePath, dbName string) error {
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
			err2 := g.write2Mongo(client, dbName, sheet.Name, data)
			if err2 != nil {
				return err2
			}
		}
	}
	return nil
}

func (g *Generate) write2Mongo(client *mongo.Client, dbName, sheetName string, data []any) error {
	if len(data) == 0 {
		return nil
	}
	collection := client.Database(dbName).Collection(sheetName)
	_, err := collection.InsertMany(context.TODO(), data)
	return err
}
