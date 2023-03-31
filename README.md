# QQ-ChatGPT-Bot
实现openai qq对话功能
* 支持连续对话(感谢[@CsterKuroi](https://github.com/CsterKuroi))
* 支持预设(感谢[@lvyonghuan](https://github.com/lvyonghuan))

![image.png](https://s2.loli.net/2023/03/27/6VJEKkDsA8dIBzL.png)

## 如何使用
### 前置工作
1. 前往https://beta.openai.com/account/api-keys 获取api_key
2. 大陆用户安装[Clash](https://github.com/Dreamacro/clash/releases) Linux参照https://zhuanlan.zhihu.com/p/396272999
### 正式开始
1. 下载[go-cqhttp](https://github.com/Mrs4s/go-cqhttp/releases)
2. 下载[QQ-ChatGPT-Bot](https://github.com/SuInk/QQ-ChatGPT-Bot/releases)
### Windows
* 双击go-cqhttp可执行文件，按照提示登录QQ,选择2正向WebSocket
* 双击QQ-ChatGPT-Bot可执行文件，将openai的api_key 填入`config.cfg`中，再次运行
* 如果要使用角色预设功能，则请在`config.cfg`中的identity下填写想要bot扮演的角色的信息。同时，请将openai配置下的model更换成“text-davinci-003”。
* 如果要使用连续对话，请在`config.cfg`中的context下进行设置。如果要启用角色预设，则不支持连续对话。
### Linux
```bash
./go-cqhttp*
# 按照提示操作,选择2正向websocket，将本地登录过的`sesssion.token`复制进服务器，防止tx风控
./QQ-ChatGPT*
# 在config.cfg填入openai的api_key
# 其它配置参考windows的说明
# 关掉窗口，运行：
nohup ./go-cqhttp* &
nohup ./QQ-ChatGPT* &
```
### 手动运行
```bash
先运行go-cqhttp
git clone git@github.com:SuInk/QQ-ChatGPT-Bot.git
cd QQ-ChatGPT-Bot
go run main.go
# 然后根据提示信息修改config.cfg文件
# 再次执行: 
go run main.go
```
## 配置文件
### cq-http配置文件
```yaml
# config.yaml cqhttp配置文件
servers:
  # 添加方式，同一连接方式可添加多个，具体配置说明请查看文档
  #- http: # http 通信
  #- ws:   # 正向 Websocket
  #- ws-reverse: # 反向 Websocket
  #- pprof: #性能分析服务器
  # 正向WS设置
  - ws:
      # 正向WS服务器监听地址
      address: 0.0.0.0:8080
      middlewares:
        <<: *default # 引用默认中间件
```
### QQ-ChatGPT-Bot配置文件
```bash
...
[openai]
# 你的 OpenAI API Key, 可以在 https://beta.openai.com/account/api-keys 获取
api_key = "sk-xxxxxx" ## 必填
# openai是否走代理，默认关闭
use_proxy = false ## 中国大陆地区需开启
# Clash默认代理地址 Linux使用Clah参照https://zhuanlan.zhihu.com/p/396272999
proxy_url = "http://127.0.0.1:7890"
...
```
### 对话指令
* 在启用连续对话的情景下，聊天中输入`/clean`将清除之前的对话记录。