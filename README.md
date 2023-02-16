# QQ-ChatGPT-Bot
实现openai qq对话功能，原生跨平台
![Screenshot_2023-02](https://s2.loli.net/2023/02/16/zJXgnOxRY1w4jZE.jpg)

## 如何使用
1. 下载[go-cqhttp](https://github.com/Mrs4s/go-cqhttp/releases)
2. 下载[QQ-ChatGPT-Bot](https://github.com/SuInk/QQ-ChatGPT-Bot/releases)
### Windows
* 双击go-cqhttp可执行文件，按照提示登录QQ
* 双击QQ-ChatGPT-Bot可执行文件，将openai的api_key 填入`config.cfg`中，再次运行
### Linux
```bash
./go-cqhttp*
# 按照提示操作，将本地登录过的`sesssion.token`复制进服务器，防止tx风控
./QQ-ChatGPT*
# 在config.cfg填入openai的api_key 
# 关掉窗口，运行：
nohup ./go-cqhttp* &
nohup ./QQ-ChatGPT* &
```
