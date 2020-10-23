### 参考索引

这一节描绘了 bin/bench 的几个 subcommands，以及它们的相关参数。

#### `exec`
basic usage:

```bash
	bin/bench exec $component_name [ test script args... ] -- [ component args... ]
```
其中，test script args 是传递给测试脚本的参数，支持如下：
- `--hash [str]`: 组件的 commit hash。
- `--workload-name [str]`: workload 的名字，这个参数必须被指定。
- `--workload-storage [str]`: workload 的地址，假如说 workload 没有保存在默认的 s3 bucket 中，可以使用这个指定完整的地址。使用 BRIE 风格的 storage 字符串。
- `--repo [str]`: 组件的 repo.
- `--pr`: 如果指定这个选项，测试将会分别用 `--hash` 和 `--repo` 指定的仓库和官方最新的 master 版本分别进行一次 bench，并且在 report 中输出两者之间的比较。
- `--import-to [str]`：如果指定这个选择，在测试完成后，api_server 还会将测试完成之后的数据库利用 br 做一份备份到指定的文件夹。利用这个选项和 lightning 可以将一些 csv 文件转化为 br 的备份文件。
- `--{tikv|pd|tidb}.version` 指定集群中对应组件的版本号。
- `--{tikv|pd|tidb}.hash` 指定集群中对应组件的 commit hash，当它与版本号同时出现的时候，版本号将会被忽略。
- `cluster-version` 指定集群版本号。

在 ‘--’ 之后的参数都会直接传递给相应的组件。需要注意测试框架本身也会传递一些参数，因此可能会有互相覆盖的窗口出现。

除此之外，还有一些组件特定的参数，这些参数也在 test script args 中：

```
（Lightning）
      --lightning.backend string the backend that lightning uses for benching (default "local")
      --lightning.index-concurrency uint   lightning.index-concurrency
      --lightning.io-concurrency uint      lightning.io-concurrency
      --lightning.range-concurrency uint   tikv-importer.range-concurrency
      --lightning.region-split-size uint   tikv-importer.region-split-size
      --lightning.send-kv-pairs uint       tikv-importer.send-kv-pairs
      --lightning.table-concurrency uint   lightning.table-concurrency

（BR）
      --br.skip-backup                     skip the backup step of br benching

（Dumpling）
      --dumpling.skip-csv                  skip dumpling to csv step in dumpling benching
      --dumpling.skip-sql                  skip dumpling to sql step in dumpling benching
```

例子：

```bash
bin/bench exec dumpling --workload-name tpcc30-br -- -F 256M
```

这个会执行 dumpling 的测试，workload 为 `tpcc30-br`，并且将 -F 256M 传递给 dumpling。
	
#### `get_cluster`

basic usage:

```bash
bin/bench get_cluster [$cluster_spec] [$type]
```

其中，`cluster_spec` 是集群的 ID。
如果留空，会提示在最近启动的集群中选择。
如果输入为一个点（`.`），等价于最近一次启动的集群。

type 可以是 info，resource 或者 metric。默认为 info。
info 会获得集群的基本信息。
resource 会获得集群中各个组件的信息。
metric 或者 grafana 会获得集群的 grafana 地址（这个 grafana 会在测试完成之后销毁）。

例子：
```bash
bin/get_cluster . grafana
```

这个命令会获得最近一次测试的 grafana 地址。

#### get_file

```bash
bin/bench get_file [$cluster_spec] [$file_name]
```

其中，cluster_spec 和 get_cluster 相同。

file_name 被忽略的时候，可以交互式地选择。
