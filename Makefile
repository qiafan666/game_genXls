e .PHONY: fmt

savePath := sxls
readPath := /Users/go/src/ning.com/game_genXls/excel
mainPath := /Users/go/src/ning.com/mainproject/config/sxls
gitPath := /Users/go/src/ning.com/Design

fmt:
	go fmt ./sxls/struct.go

git_pull:
	cd $(gitPath) && git pull

gen_struct:
	go run main.go --savePath $(savePath) --readPath $(readPath) --genType struct

gen_json:
	go run main.go --savePath $(savePath) --readPath $(readPath) --genType json

gen_mongo:
	go run main.go --savePath $(savePath) --readPath $(readPath) --genType mongo --mongoUri "mongodb://10.80.10.109:2701" --mongoDb "config_dev"

cp:
	cp -R $(savePath)/struct.go $(mainPath)
##&& cp -R $(savePath)/json $(mainPath)

struct: gen_struct fmt
json: gen_json fmt
mongo: git_pull gen_struct gen_mongo fmt cp
