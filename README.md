## golang web开发脚手架,集成了redis,mongo,influxdb数据库相关api

### 介绍
1. 采用　gin + xorm 自动生成增删改查api，可重写增删改查api函数　# golang的orm库中xorm其实要比gorm和beego/orm要好用,命名,用法都很清晰，地址 http://www.xorm.io/
2. 程序启动的默认配置都在setting包下面，setting优先从环境变量里面读取配置
3. docker-compose文件夹下面有常用的一些组件的docker-compose.yml以及相关配置，可以直接使用，或者自行更改配置

