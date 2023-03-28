# QQ-ChatGPT-Bot
实现openai qq对话功能
* 支持连续对话(感谢[@CsterKuroi](https://github.com/CsterKuroi))
* 支持预设(感谢[@lvyonghuan](https://github.com/lvyonghuan))
![image.png](https://s2.loli.net/2023/03/27/6VJEKkDsA8dIBzL.png)

## 如何使用
1. 下载[go-cqhttp](https://github.com/Mrs4s/go-cqhttp/releases)
2. 下载[QQ-ChatGPT-Bot](https://github.com/SuInk/QQ-ChatGPT-Bot/releases)
### Windows
* 双击go-cqhttp可执行文件，按照提示登录QQ,选择2正向WebSocket
* 双击QQ-ChatGPT-Bot可执行文件，将openai的api_key 填入`config.cfg`中，再次运行
* 如果要使用角色扮演功能，则请在config中的identity下填写想要bot扮演的角色的信息。同时，请将openai配置下的model更换成“text-davinci-003”。如不想使用角色扮演功能，请直接删掉identity下prompt里双引号内的内容。
### Linux
```bash
./go-cqhttp*
# 按照提示操作,选择2正向websocket，将本地登录过的`sesssion.token`复制进服务器，防止tx风控
./QQ-ChatGPT*
# 在config.cfg填入openai的api_key 
# 关掉窗口，运行：
nohup ./go-cqhttp* &
nohup ./QQ-ChatGPT* &
```
```bash
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
