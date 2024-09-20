package lib

import (
	"encoding/json"
	"fmt"
	"os"

	"game_genXls/sxls"
)

func ReadFromJson(savePath string) {
	bytes, err := os.ReadFile(savePath + "/" + "xls.json")
	if err != nil {
		fmt.Printf("open json file error:%v", err)
	}
	var xls sxls.RawConfig
	err = json.Unmarshal(bytes, &xls)
	if err == nil {
		fmt.Println("success")
	} else {
		fmt.Println(err)
	}
}
