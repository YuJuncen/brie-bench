### 快速上手

这个仓库可以做什么：
- 向 api_server 发起请求，自动启动集群并且完成相应 workload 的测试
- 自由选择想要的组件版本
- 自由为组件提供命令行参数

首先 clone 仓库：
```bash
git clone https://github.com/yujuncen/brie-bench.git
cd brie-bench
```
然后构建虚拟环境并安装依赖：

```bash
make ctl-build
```

> 这一步会依赖 python3。

这时候会询问环境相关的参数，默认值是 Junpeng 老师现在搭建的测试集群，如果仅仅是测试用，所有选项留空，bucket 选择 &lt;not changed&gt; 就好。

这时候一切已经准备妥当，使用 bin/bench 这个小脚本即可启动测试，例如，启动一个 br 的测试：
```
bin/bench exec br \
    --workload-name tpcc30-br \
    --tidb.version 4.0.5 \
    --hash 7d62d5ff67e38edb82d5224d5904462bfce5dcf9
```
exec 的第二个参数是组件的名字，目前支持 br，lightning，dumpling。workload-name 项指定了 workload 的名字，默认情况下，这个程序会在之前指定的 bucket 中寻找相应的 workload。
--tidb.version 指定了 tidb 的版本，--hash 则指定了组件版本。更加详细的参数说明可以在下面的参考索引章节找到。
成功了之后，应该会获得集群的 ID。
启动之后，可以通过 get_cluster 获得相关的信息，例如：
```bash
bin/bench get_cluster . grafana
```
其中，’.’ 标志指代我们上一次启动的集群，grafana 指代获得 grafana 监控地址。

也可以直接获得集群的运行状况等相关信息，在省略最后的资源类型的时候，get_cluster 默认返回集群的基本信息：
```bash
bin/bench get_cluster . ｜ jq .status
```

可以通过 get_file 获得测试相应的产物：
brie-bench$ bin/bench get_file .
因为没有提供文件名，你可以在所有文件中选择。

{backup,restore}.log 是备份和恢复的日志。
report.md 是测试的报告。
stdout.log 是测试脚本的输出。

IP 地址对应的文件夹是集群中相应组件的日志。
如果需要和其他命令行工具联动，可以直接提供文件名，例如：
```bash
bin/bench get_file . restore.log | grep "summary"
```

至此，一个测试已经完成了。
