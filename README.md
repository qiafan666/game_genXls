# 服务端导表工具

## Getting started

修改 makefile 中的 各个路径，然后运行 `make` 命令。

## 配置表结构
前四行为表头，后面为数据，例子可看excel中的称号表实例。

## 数据类型

byte = int8
int = int32
long = int64
float = float64
string = string
List<> 代表数组，如  List<int> 代表 int32 数组。

## 自定义数据结构
自定义类型需要在 struct_自定义数据结构表中配好，然后再引用。