# BRIE Benching

This repository contains the workload and some utils to run those cases.

**(WIP)**

## usage

```bash
bin/exec component [--dry-run] [args passing to main.go...] [-- [args passing to the component...]] 
```

This command will make request to the default API server(172.16.5.110:8000), 
then the API server would create a cluster and run the workload. 

The server would response with a `cluster_id`, which can be used 
by other utils then.

### binaries (WIP)

- `get_file [cluster_id] [filename]`: Get the output file of the test from the `minio` s3 endpoint. 
(Need run `mc alias set minio http://172.16.4.18:30812 YOURACCESSKEY YOURSECRETKEY` firstly for the default API server).
- `get_cluster [cluster_id]`: Get the cluster info.
- `rebuild_metrics [cluster_id]`: Rebuild the grafana metrics. (WIP)

The cluster ID can be a dot(`.`), which means get the last requested cluster.  
The cluster ID can be absent, and then you can select one from recent requests.

### Parameters

For now, `BR` and `Dumpling` are supported component. Below lists other command line parameters.

- `--hash`: The commit hash of the component.
- `--workload`: The workload to run.
- `--repo`: The repository the component built from.
- `--`: All flags after this would be passed to the component directly, aware this might be overridden by the framework.  
 
### Component workload

For now, for all components, the workload is a snapshot of the database. It can be represented by many forms: including
BR backup, CSV or SQL files. 

Each component should do several cases (steps) with this workload, as described below.

#### BR

Input: a BR backup instance.

1. BR will restore the backup to the cluster.
2. BR will backup the cluster to another place. (can be skipped by `--br.skip-backup`)

#### dumpling

Input: a BR backup instance.

1. The framework will restore the backup to the cluster.
2. Dumpling will dump the cluster to CSV file. (can be skipped by `--dumpling.skip-csv`)
3. Dumpling will dump the cluster to CSV file. (can be skipped by `--dumpling.skip-sql`)

