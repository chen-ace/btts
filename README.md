# Batch TTS
利用阿里云智能语音交互-语音生成服务，批量生成语音文件。

## 编译
### 安装 Go
要编译此项目，您需要先安装 Go。您可以从 [Go 官方网站](https://golang.org/dl/) 下载并安装最新版本的 Go。

本项目要求 Go 的版本不低于 1.20。

### 编译本机可用的命令
安装 Go 后，您可以使用以下命令来编译此项目：

```bash
go build -o btts
```
这将在当前目录下生成一个名为 `btts` 的可执行文件。

### 编译其他平台
本项目的Makefile自带了交叉编译的支持， 您也可以使用 `make` 命令来编译项目：

```bash
make
```

## 使用方法

本应用接受几个命令行参数：

- `-i, --input`：指定输入剧本文件的路径。
- `-o, --output`：指定保存音频文件的目录。
- `-h, --help`：打印帮助信息。
- `-c, --clean`：清空输出目录。
- `-r, --random`：随机选择n个项目进行测试。如果设置为0，将生成所有项目。

以下是如何使用这些选项的示例：

```bash
$ ./btts -i /path/to/input -o /path/to/output -r 10
```

## 输入文件格式
本项目暂时接受csv格式的文件作为输入。csv文件的格式如下：

```csv
项目,镜头,发言角色,序号,文本,发言人,语调,语速,音量,音频格式
测试,1,旁白,1,在数字世界中，每个字符、每行代码都承载着特殊的意义。而在这个世界里，字符串扮演着一个不可或缺的角色。但你有没有想过，为什么在像JavaScript这样的语言中，字符串是不可变的？,zhida,0,25,69,wav
测试,1,虎,1,小龙啊，你从事前端开发6/7年了，你可知道javascript中字符串为什么是不可变的吗？,aida,0,25,69,wav****
```

## 参考信息
- [发言人取值参考](https://help.aliyun.com/document_detail/155645.html?spm=a2c4g.11186623.6.540.6)
- 音量，取值范围：0～100。默认值：50。
- 语速，取值范围：-500～500，默认值：0。 [-500, 0, 500] 对应的语速倍速区间为 [0.5, 1.0, 2.0]。1倍速是指模型默认输出的合成语速，语速会依据每一个发音人略有不同，大概每秒钟4个字左右。


- 语调，取值范围：-500～500，默认值：0。
- 音频格式，支持的音频格式：wav、mp3。默认值：wav。

## 配置文件
配置文件的格式如下：
```json
{
  "ACCESS_KEY": "xxx",
  "ACCESS_SECRET": "xxx",
  "APP_KEY": "xxx"
}
```
或者指定阿里云URL：
```json
{
  "URL": "https://xxx.com",
  "ACCESS_KEY": "xxx",
  "ACCESS_SECRET": "xxx",
  "APP_KEY": "xxx"
}
```
配置文件请放置于`~/.btts/`目录下,文件名为`config.json`。第一次启动程序时，程序会询问用户，并自动创建该文件，您也可以手动创建。