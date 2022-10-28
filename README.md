# lan

[![License](https://img.shields.io/github/license/Lanly109/lan)](LICENSE)
![Golang Version](https://img.shields.io/badge/Golang-1.18.2-blue)

综合[filter](https://github.com/Lanly109/filter-v3)等工具，集`清理无关文件`、`名单检查`、`md5生成`、`修改时间检查`、`查重`等为一体的小工具，便于收代码后等一系列校对操作。

因为是众多工具的集合，懒人福音，故取名为`lan`，意为`懒`。

## 流程

1. 召集考点负责人开会，各考点确定学校和试室的缩写，方便后面汇总。
2. 准备配置用的考生名单`namelist`，每个考点一份，包含考号和试室缩写。（csv, 两列）
3. 考点负责人学习本工具的使用
4. 考点负责人配置各试室的 `toml` 文件、考点的 `toml` 文件。

为避免跨系统带来的中文编码影响，请务必全程仅使用**英文**，包括`csv`的内容，`toml`配置文件。

如需使用中文，注意程序默认中文编码为`UTF-8`。

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
        │   └── expr.cpp
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

`windows`用户在文件所在目录下按住`shift`键及鼠标右键选择`在 xxxxxx 中打开`，即可打开终端

以下命令，`cmd终端`**不需要**前面的`./`

### 使用流程
```bash
# 生成配置文件，并修改
./lan config

# 清理无关文件
./lan clean
# 检查成员名单
./lan check
# 生成md5码文件，并下发
./lan md5
# 检查文件修改时间
./lan time
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
./lan config
``` 

该指令在当前工作路径下生成名为`config.toml`的配置文件，配置内容适用于`demo`

可接`--name=<name>`参数自定义配置文件名字

```bash
./lan config --name=myconfig.toml
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
Extensions : 支持的扩展名
``` 

该操作会将所有位于`SourcePath`文件夹里的形如`GD-xxxx/problem/problem.ext`的`ext`文件复制到`CodePath`文件夹。

有异常文件（如`代码名`不是`题目名`，`代码文件`的路径深度不正确会有警告信息）

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

该操作会将`NameList`中为`Room`的考号与`CodePath`中考号比对，给出缺少考号以及不应存在的考号列表。

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

### 修改时间检查

```bash 
./lan time
``` 

配置说明：

```toml
CodePath : 待检查的文件夹路径
StartTime : 比赛开始时间
EndTime : 比赛结束时间
AbnormalLog : 修改时间异常的学生清单
``` 

该操作会将`CodePath`的所有修改时间不在比赛时间内的文件列一份清单，保存在`AbnormalLog`中。

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
./lan clean 304 raw_304 --problems=expr,live,number,power --extensions=.cpp,.c,.pas
./lan check 304 --room=304 --namelist=namelist.csv
./lan md5 304
./lan time 304 --starttime="2021-11-20 08:30:00" --endtime="2021-11-17 13:00:00"
./lan moss 304 --problem=expr --language=cc --maxlimit=10 --numberresult=250 --userid=xxxx
``` 

各参数的优先级：`命令行参数` > `配置文件` > `默认值`
