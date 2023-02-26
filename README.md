# lan

[![License](https://img.shields.io/github/license/Lanly109/lan)](LICENSE)
![Golang Version](https://img.shields.io/badge/Golang-1.18.2-blue)

综合[filter](https://github.com/Lanly109/filter-v3)等工具，集`清理无关文件`、`名单检查`、`md5生成`、`修改时间检查`、`查重`等为一体的小工具，便于收代码后等一系列校对操作。

因为是众多工具的集合，懒人福音，故取名为`lan`，意为`懒`。

## 更新日志

### 2023.2.26

- `time`命令更名为`valid`命令，新增`SourceSizeLimit`参数，在检查时间戳的基础上追加代码文件大小检查
- 更新`demo`

### 2022.11.5

- 新增`已知缺席考生`特性，优化`clean`命令的结果描述
- 新增`gen`命令，原先`config`命令移至到`gen`的子命令下，同时新增`countdown`和`share`子命令，用于生成`倒计时html模板`和`打开共享文件夹源代码`文件
- 修复了重复执行`md5`命令时，输出文件内容异常的情况
- 更新`demo`

### 2022.10.28

- 修复了`moss`命令时相对路径解析错误的问题
- 修改了`md5`命令默认保存文件名

### 2022.10.15

- v0.1版发布

## 流程

1. 召集考点负责人开会，各考点确定学校和试室的缩写，方便后面汇总。
2. 准备配置用的考生名单`namelist`，每个考点一份，包含考号和试室缩写。（csv, 两列, 不包含说明头）
3. 考点负责人学习本工具的使用
4. 考点负责人配置各试室的 `toml` 文件、考点的 `toml` 文件。

为避免跨系统带来的中文编码影响，请务必全程仅使用**英文**，包括`csv`的内容，`toml`配置文件。

如需使用中文，注意程序默认中文编码为`UTF-8`。

### NameList制作教程

`NameList`应形如

```
GD-00018,402
GD-00032,304
GD-00062,402,0
GD-00077,304
GD-00081,304
GD-00111,304,0
GD-00128,402
GD-00139,304,0
GD-00153,402
GD-00192,304
GD-00291,402
```

第一列为考生号，第二列为试室号。应避免中文。

可用`excel`制作，保存时请选择`csv`格式，不要选择`csv UTF-8`格式。

> 后者格式实际为`UTF-8 BOM`，与前者相比，文件头会多处`EF BB BF`字段，影响本程序对该文本的读取。

第三列为可选列，在各考场签到完毕后，对于缺考考生，可在其考号后增加第三列，其值为`0`。本程序将认为该考生为`已知缺考考生`，在`check`命令中与`未知缺考考生`加以区别。

## 工具使用指南

供学习测试用的[demo](https://github.com/Lanly109/lan/releases/download/demo/demo.zip)

包括可执行文件如下方式放置，配置文件`config.toml`在下面有生成方法，`raw_304`为收取代码后的文件夹结构
```bash
.
├── config.toml
├── lan
├── namelist.csv
└── raw_304
    ├── GD-00032
    │   ├── expr
    │   ├── live
    │   ├── number
    │   └── power
    │       └── power.cpp
    ├── GD-00077
    └── GD-00081
        ├── expr
        ├── live
        │   └── live.cpp
        └── power
            ├── power
            └── power.cpp
``` 


### 安装

在[release](https://github.com/Lanly109/lan/releases)中选择适合自己系统和架构的可执行文件下载。

下面假设可执行文件名为`lan`

`linux`用户需给程序添加执行权限
```bash 
chmod +x lan
``` 

### 使用

请于终端内运行。

`Windows`用户在文件所在目录下按住`shift`键及鼠标右键选择`在 xxxxxx 中打开`，即可打开终端

以下命令，`cmd终端`**不需要**前面的`./`

`Windows`建议使用`Windows Terminal`作为终端模拟器，支持色彩显示

### 使用流程
```bash
# 生成配置文件，并修改
./lan gen config

# 清理无关文件
./lan clean
# 检查成员名单
./lan check
# 生成md5码文件，并下发
./lan md5
# 检查文件修改时间和文件大小
./lan valid
# 查重
./lan moss
```

### 帮助

```bash 
./lan help
``` 

该指令会显示帮助文案

### 生成配置文件

```bash 
./lan gen config
``` 

该指令在当前工作路径下生成名为`config.toml`的配置文件，配置内容适用于`demo`

可接`--name=<name>`参数自定义配置文件名字

```bash
./lan gen config --name=myconfig.toml
``` 

该指令在当前工作路径下生成名为`myconfig.toml`的配置文件，配置内容适用于`demo`

---

后续指令默认读取`config.toml`配置文件，如需更改，可接`--config=<name>`参数自定义读取配置文件

### 清理无关文件

```bash 
./lan clean
``` 

配置说明：

```toml
CodePath : 过滤后的文件夹路径
Problems : 比赛题目
SourcePath : 待清理文件夹
Extensions : 保留的扩展名
IgnoreExtensions : 忽略的扩展名
``` 

该操作会将所有位于`SourcePath`文件夹里的形如`GD-xxxx/problem/problem.ext`的`ext`文件复制到`CodePath`文件夹。

有异常文件（如`代码名`不是`题目名`，`代码文件`的路径深度不正确，有考生文件夹但**无**有效代码，有对应题目文件夹但**无**题目代码等会有警告信息）

---

后续指令均默认文件夹已经过清理，**无无关文件**（如`姓名.txt`,`*.pdf`,`*.in`,`*.out`,`*.ans`,`*.exe`)等

### 名单检查

```bash 
./lan check
``` 

配置说明：

```toml
CodePath : 待检查的文件夹路径
Room : 检查的试室， all 为全部
NameList : 名单文件
``` 

该操作会将`NameList`中为`Room`的考号与`CodePath`中考号比对，给出缺少考号以及不应存在的考号列表，同时会分别给出`已知缺考`(从签到表里得知缺考)和`未知缺考`的考号

### md5生成

```bash 
./lan md5
``` 

配置说明：

```toml
CodePath : 待生成md5的文件夹路径
Md5File : md5文件名，非必要不更改
``` 

该操作会将`CodePath`的所有文件生成一份`md5`表单，配合[checker](https://github.com/xfoxfu/checker)使用，

### 修改时间和文件大小检查

```bash 
./lan valid
``` 

配置说明：

```toml
CodePath : 待检查的文件夹路径
StartTime : 比赛开始时间
EndTime : 比赛结束时间
SourceSizeLimit : 文件最大大小（单位:byte）
AbnormalLog : 修改时间异常的学生清单
``` 

该操作会将`CodePath`的所有修改时间不在比赛时间内以及大小超出限制的文件列一份清单，保存在`AbnormalLog`中。

### 查重

```bash 
./lan moss
``` 

配置说明：

```toml
CodePath : 待检查的文件夹路径
ReviewProblem : 查重题目
ReviewUserID : 账号
ReviewComment : 查重注释
ReviewLanguage : 代码语言
ReviewMaxLimit : 当同样的代码出现文件数大于该次数时，认为不是抄袭代码
ReviewExperimental : 启用新特性检查
ReviewNumberResult : 显示的结果数
``` 

该操作会将`CodePath`的`ReviewProblem`代码进行查重
 
---

上述命令的参数均可在命令行设置，具体用法参见`help`指令

与上述配置文件配置等效的命令行如下

```bash
./lan clean 304 raw_304 --problems=expr,live,number,power --extensions=.cpp,.c,.pas --ignoreexts=.txt,.in,.out,.ans,.pdf,.exe
./lan check 304 --room=304 --namelist=namelist.csv
./lan md5 304
./lan valid 304 --starttime="2021-11-20 08:30:00" --endtime="2021-11-17 13:00:00" --sizelimit=102400
./lan moss 304 --problem=expr --language=cc --maxlimit=10 --numberresult=250 --userid=xxxx
``` 

各参数的优先级：`命令行参数` > `配置文件` > `默认值`

### 生成打开共享文件夹的脚本

```bash 
./lan gen share
``` 

该指令在当前工作路径下生成名为`sharing.cpp`的脚本代码，修改`27行`的`\\\\10.78.30.99\\b404`为对应路径（在`Windows`的资源管理器的路径输入`\\10.78.30.99\b404`即可访问到共享文件夹，注意`\`的转义），然后编译出可执行文件下发。

可接`--name=<name>`参数自定义脚本文件名字

### 生成比赛倒计时的网页

```bash 
./lan gen countdown
``` 

该指令在当前工作路径下生成名为`countdown.html`的html网页，修改`7,8,9`行对应参数，注意时间的格式。

注意修改文件后，浏览器要刷新才能生效。

可接`--name=<name>`参数自定义脚本文件名字
