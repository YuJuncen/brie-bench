### 开发文档

这个仓库分为两个部分：向大集群 API Server 发起请求的 `ctl` 和作为 workload 被执行的 `workload`，整体架构如下：

```
+--------+         +-----------------+
|        | request |                 |
|  ctl   +--------->  API Server     |
|        |         |                 |
+--------+         +--------+--------+
                            |
                            | pull and execute in docker
                            |
                            |
                   +-----------------+
                   |        |        |
                   | +------v------+ |
                   | |  Workload   | |
                   | |             | |
                   | +-------------+ |
                   |                 |
                   |  Docker Hub     |
                   |                 |
                   +-----------------+
```

Workload 部分已经设置了 GitHub Action，每次 push 的时候都会将其上传到 DockerHub。

#### ctl

ctl 部分会构造到 API Server 的请求，同时也有一些从 S3 拉取造物（artifacts）的脚本。
这部分使用 python 写成，主要代码在 `ctl/lib/` 中，这些文件的功能如下：

```
ctl/lib
├── cli_parse.py -- 命令行解析
├── exec.py -- exec 命令
├── get_cluster.py -- get_cluster 子命令
├── get_file.py -- get_file 子命令
├── init_config.py -- 构造配置文件相关的代码
└── saved_clusters.py -- 查询既往集群相关的代码
```

#### workload

workload 部分是要在 docker 环境中执行的脚本。

```
workload
├── components -- 各个组件的构建和运行脚本
│   ├── br.go
│   ├── component.go
│   ├── dumpling.go
│   ├── lightning_conf.go
│   └── lightning.go
├── config -- 命令行参数配置
│   ├── config.go
│   └── constants.go
├── main.go -- 入口
└── utils 
    ├── cluster.go -- bench 的上下文相关的操作
    ├── cmd.go -- 运行 bash commands 的工具
    ├── debug.go -- 一些保存调试信息相关的代码
    ├── git
    │   └── git.go -- git 命令的抽象 
    ├── http_wrapper.go -- http 相关的操作
    ├── metric
    │   ├── bench.go -- 生成报告相关的代码
    │   ├── bench_test.go
    │   ├── size.go -- 字节数量格式化输出的相关代码
    │   └── values.go -- 对 bench 量化数据的抽象
    ├── misc.go -- 一些杂物：socket addr parser 等等
    ├── pd
    │   └── client_ext.go -- 对 PD 的操作的抽象
    └── storage
        └── s3.go -- s3 临时存储 （BR 用）
```

### TODO

- more workloads?
- `/bench` for BRIE?
- bug fixing and useability improvement?