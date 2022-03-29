# config 基于yaml的配置
> 项目配置文件管理

### 基于yaml.v3实现
通过启动参数 -c 设置配置文件路径，默认加载 ./conf/local.yml

#### 强制加载指定配置

```go
//通过启动参数实现
./program -c *.yml
config.Load("config file path")
```