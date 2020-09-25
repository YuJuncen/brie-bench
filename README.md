# BRIE Benching

This repository contains the workload and some utils to run those cases.

**(WIP)**

## usage

```bash
exec component [--dry-run] [args passing to main.go...]
```

This command will make request to the default API server(172.16.5.110:8000), 
then the API server would create a cluster and run the workload. 

The server would response with a `cluster_id`, which can be used 
by other utils then.

### binaries (WIP)

- `get_stdout`: Get the `stdout` of the test from the `minio` s3 endpoint. 
(Need run `mc alias set minio http://172.16.4.18:30812 YOURACCESSKEY YOURSECRETKEY` firstly for the default API server).
- `get_artifact cluster_id`: Get the artifacts of the cluster, available only 
after cluster status become "DONE". 
- `get_cluster cluster_id`: Get the cluster info.
- `rebuild_metrics`: Rebuild the grafana metrics. (WIP)

### Parameters

For now, `BR` is the only supported component. Below lists other command line parameters.

- `-hash`: The commit hash of the component.
- `-workload`: The workload to run.
 


