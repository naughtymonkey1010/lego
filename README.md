#lego  - happy use go 
go 脚手架

#### 特性
- 封装对接一点常用服务,降低开发使用成本
- 集成viper配置管理
- 集成logrus日志管理, 支持多日志配置, 支持yd日志格式 便于k8s收集
- 集成gin http server
- 集成grpc server
- 集成 gocron 定时任务调度 支持秒级别定时和指定时间定时
- 集成 redis, codis(自开发) redis 客户端 
- 集成 zookeeper 客户端
- 集成 mongo 客户端
- 集成 httplib(来源beego) http请求组件
- 集成 swagger ui
- 接管信号 支持http grpc graceful shutdown
- 集成 ydmetric(自开发) mon打点
- 集成 dingding robot机器人(自开发), 安全加签模式 支持发送text、link、markdown 类型消息
- 脚手架核心只依赖配置管理, 日志, gin, grpc 模块, 包尽量小, 其他模块以组件形式提供

项目使用demo参考
参考demo:  https://git.yidian-inc.com:8021/image/lego-demo