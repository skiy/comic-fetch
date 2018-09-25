漫画采集

------
## 采集说明
### 支持站点
- 漫画160 https://www.mh160.com

### 使用说明
- 新增漫画：
向 redis db1 添加数据（格式） ```set newbooks '[{"id":25510,"flag":"mh160"},{"id":11106,"flag":"mh160"},{"id":11105,"flag":"mh160"}]'```
或向 newbooks.json 文件添加数据 (格式) ```[{"id":25510,"flag":"mh160"},{"id":11106,"flag":"mh160"},{"id":11105,"flag":"mh160"}]```

- 更新漫画
使用 ```go run main.go``` 运行本项目

### TODO
- 数据库、缓存 可配置 √
- 抓取图片资源至本地
- GoWeb 方式展示(不再依赖PHP) (b
- 支持 sqlite √