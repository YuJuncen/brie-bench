# BRIE Benching

This repository contains the workload and some utils to run those cases.

**(WIP)**

## usage

Firstly, make a config file:

```bash
make ctl-build
```

This command will set up python venv and install requirements needed. 
There would also be a prompt asking you the s3 parameters. Current you can leave all unchanged.

> Hint:
> You can always change the config by `bin/bench create_config`.

```bash
bin/bench exec component [--dry-run] [args passing to main.go...] [-- [args passing to the component...]] 
```

This command will make request to the API server, then the API server would create a cluster and run the workload. 

The server would response with a `cluster_id`, which can be used by other utils then.

### bench sub-commands (WIP)

- `get_file [cluster_id] [filename]`: Get the output file of the test.
- `get_cluster [cluster_id] [(info|metrics)]`: Get the cluster info, or grafana address.
- `rebuild_metrics [cluster_id]`: Rebuild the grafana metrics. (WIP)

The cluster ID can be a dot(`.`), which means get the last requested cluster.  
The cluster ID can be absent, and then you can select one from recent requests.

### Parameters

For now, `BR` and `Dumpling` are supported component. Below lists other command line parameters.

- `--hash`: The commit hash of the component.
- `--workload-name`: The workload to run.
- `--workload-storage`: The storage for workload.
- `--repo`: The repository the component built from.
- `--`: All flags after this would be passed to the component directly, aware this might be overridden by the framework.  
 
### Component workload

For now, for all components, the workload is a snapshot of the database. It can be represented by many forms: including
BR backup, CSV or SQL files. 

Each component should do several cases (steps) with this workload, as described below.

The `input` could be specified by `--workload-storage` directly, or just provide `--workload-name` to let the framework
deriving where the workload is.

#### BR

Input: a BR backup instance.

1. BR will restore the backup to the cluster.
2. BR will backup the cluster to another place. (can be skipped by `--br.skip-backup`)

#### dumpling

Input: a BR backup instance.

1. The framework will restore the backup to the cluster.
2. Dumpling will dump the cluster to CSV file. (can be skipped by `--dumpling.skip-csv`)
3. Dumpling will dump the cluster to CSV file. (can be skipped by `--dumpling.skip-sql`)

#### Lightning

When testing tidb-lightning, you should provide which backend (local or TiDB) to use by `--lightning.backend`.

Input: a Folder of CSV or SQL.

1. Lightning will run that case.

There are some extra flags for configuring lightning, they would be mapped to 
variables in the lightning config file when specified. 

- lightning.index-concurrency     (by `--lightning.index-concurrency`)
- lightning.io-concurrency        (by `--lightning.io-concurrency`) 
- lightning.table-concurrency     (by `--lightning.table-concurrency`)
- tikv-importer.region-split-size (by `--lightning.region-split-size`)
- tikv-importer.send-kv-pairs     (by `--lightning.send-kv-pairs`)
- tikv-importer.range-concurrency (by `--lightning.range-concurrency`) 