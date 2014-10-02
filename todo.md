TODO
---------

Legend:

    - [x] done
    - [~] current
    - [-] not doing

### v0.1.0

- [x] need work defined
- [x] need poller
- [-] need poison pill
  - [-] need way to stop the goroutine, select?
- [-] extract and modularize
- [-] need master state monitor
  - [-] needs to be synced with etcd?

### v0.1.1

- [ ] test influxdb docker
- [ ] put influxdb out on gce with fleet
- [ ] test builder using metrics lib
- [ ] test builder writing to influxdb using metrics lib
- [ ] top level coordination is missing between the supervisors
- [ ] lots of tests