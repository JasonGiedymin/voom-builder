Builder - Worker
----------------

## Install

It is recommended to install gpm and run it. GPM is simplier to use than
godep.

```
brew install gpm
gpm
```

## Legacy Install

### Requires ETCD

Project requires ETCD which can be `go get` and requires the etcd client
to be running (for testing).

```bash
go get github.com/coreos/go-etcd/etcd
```

Then:
 - Install etcd (on a non dev machine must use go 1.2) `brew install etcd`
 - Install etcdctl (command line tool) `brew install etcdctl`

 To run:
 ```bash
 $> etcd
 ```

### Google OAUTH Libs

```bash
go get code.google.com/p/goauth2/oauth
```

## Arch TODO

Workers use `etcd` to store info such as:
  - [ ] configuration: `/app/config`
  - [ ] worker list: `/app/workers`
  - [ ] stats: `/app/stats`
  - [ ] status/heartbeat: `/app/status`

Workers read from voom `pending-jobs` pull queue.
Workers write to voom `

Workers process these job types:
  - BuildJob (job includes timing info, start/end times): `/app/buildjobs`
  - StatsJob (realtime avg times): `/app/statsjobs` ?
  - StatusJob (status of nodes, cpu/mem info etc...): `/app/statusjobs` ?

## Wiki

### Flow Diagram

Below are various things to do with graphiz.

- [shapes](http://www.graphviz.org/content/node-shapes)

In the wiki contains a flow diagram. To build it execute:

```bash
dot flow.dot -Kfdp -Tps2 -o flow.ps
open flow.ps
```

For jpeg do:
```bash
dot flow.dot -Kfdp -Tjpeg -o flow.jpg
open flow.jpg
```


## Sysadmin Info

### Credentials

Access between workers and the TaskQueue is guarded by OAuth tokens.

P12 password for oauth credentials are: `notasecret`

### SSH Keys

Make sure to have the proper permissions.

1. protect the directory: `chmod 700 ~/.ssh/`
1. protect the files within: `chmod 600 ~/.ssh/*`

### Google Compute Setup (GCE)

The project id is the system identifier, not the project name.
In this case the main project should be `voom-registry-service`. To setup
gcloud to use this project type:

```bash
$ gcloud config set project voom-registry-service
```

Common commands:

Note that `gcutil` is deprecated. Use gcloud instead.

create instance: 
```bash
gcutil --project=voom-registry-service addinstance --image=projects/coreos-cloud/global/images/coreos-stable-410-0-0-v20140902 --persistent_boot_disk --zone=us-central1-f --machine_type=n1-standard-1 --metadata_from_file=user-data:cloud-config.yml core1
```

create micro instance
```bash
gcutil --zone us-central1-f --project=voom-registry-service addinstance --image=projects/coreos-cloud/global/images/coreos-stable-410-0-0-v20140902 --persistent_boot_disk --zone=us-central1-f --machine_type=f1-micro --metadata_from_file=user-data:cloud-config.yml --scopes https://www.googleapis.com/auth/taskqueue coreworker1
```

_recommended_ `gcloud` method
```bash
gcloud compute instances create coreworker1 --project=voom-registry-service --zone us-central1-f --image=coreos --machine-type f1-micro --boot-disk-type pd-standard --metadata-from-file user-data=cloud-config.yml --scopes https://www.googleapis.com/auth/taskqueue
```

update metadata
```bash
gcloud compute instances add-metadata core1 --metadata-from-file user-data=cloud-config.yml
```

delete instance and disk: 
```bash
gcloud compute instances delete core1` && `gcloud compute disks delete core1
```

### Common Etcd api calls

1. Etcd version: `curl -L http://127.0.0.1:4001/version`
