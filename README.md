漫画采集

------
## 采集说明
### 支持站点
- 漫画160 https://www.mh160.com

### 开发说明
- 1.下载本项目到 $GOPTAH/src 下：如 ```git clone https://github.com/skiy/comicFetch```
- 2.安装依赖库：```go get -v github.com/skiy/comicFetch``` (若本程序已开源到GitHub,那么可直接执行此命令行即可)
- 3.运行程序 ```go run main.go```
- **注意：**需要安装 golang.org 的 x 扩展，若无法下载，可通过命令行下载：```git clone https://github.com/golang/net  $GOPATH/src/golang.org/x/net```

### 使用说明
- 新增漫画：
向 redis db1 添加数据（格式） ```set newbooks '[{"id":25510,"flag":"mh160"},{"id":11106,"flag":"mh160"},{"id":11105,"flag":"mh160"}]'```
或向 newbooks.json 文件添加数据 (格式) ```[{"id":31512,"flag":"mh160"},{"id":11106,"flag":"mh160"},{"id":11105,"flag":"mh160"}]```

- 程序使用
使用 mysql 方式需要导入 comic.sql 到 MySQL数据库
使用 sqlite 方式，需要使用到 comic.db
并按照具体要求配置好 conf.ini 中的信息

### TODO
- 数据库、缓存 可配置 √
- 抓取图片资源至本地 √
- GoWeb 方式展示(不再依赖PHP) (b
- 支持 sqlite √