# BRIE Benching

This repository contains the workload and some utils to run those cases.

**(WIP)**

## usage

```bash
bin/exec component [--dry-run] [args passing to main.go...]
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

For now, `BR` is the only supported component. Below lists other command line parameters.

- `--hash`: The commit hash of the component.
- `--workload`: The workload to run.
- `--repo`: The repository the component built from.
 


