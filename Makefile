e .PHONY: fmt

savePath := sxls
readPath := /Users/gongguowei/go/src/ning.com/game_genXls/excel
mainPath := /Users/gongguowei/go/src/ning.com/mainproject/config/sxls
gitPath := /Users/gongguowei/go/src/ning.com/Design

fmt:
	go fmt ./sxls/struct.go

git_pull:
	cd $(gitPath) && git pull

gen_struct:
	go run main.go --savePath $(savePath) --readPath $(readPath) --genStruct 1

gen_mongo:
	go run main.go --savePath $(savePath) --readPath $(readPath) --genStruct 2 --mongoUri "mongodb://10.80.10.109:27017" --mongoDb "config_dev"

cp:
	cp -R $(savePath)/struct.go $(mainPath)
##&& cp -R $(savePath)/json $(mainPath)

struct: gen_struct fmt
all: git_pull gen_struct gen_mongo fmt cp
