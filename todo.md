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

 - [x] supervisor
 - [x] worker
 - [x] suture
 - [x] multi-cpu
 - [x] buffer increase
 - [x] optmize buffer so work is pulled only when needed?
 - [x] cleanup
 - [x] spawn failure test
 - [x] record failures
 - [x] record current state when accepting work
 - [x] introduce concept of job 'claim'
   - [x] keep track of current claims with channels (as recommended by golang)
   - [x] keep track of failures and update claims with info
 - [x] use metrics (though this impl really should be a raw channel select
       or crdt positive-negative counter; atomic will still work yet not
       as elegant and integration with influxdb is provided)
 - [x] some tests
 - [x] re-organize
 - [x] incorporate previous config file work
 - [x] add more options to yaml file
 - [x] add work options
 - [x] add work yield
 - [x] ~~abstract work~~ it is abstract enough for now
 - [x] godoc => run `make install` so that app is in pkg, then run godoc -http=:8000
 - [x] finish etcd func to get config (just need to know what config info needed)
 - [x] fix embedded etcd config, seems to be a bug with yaml parser and embedded structs, name collisions?
 - [x] enable registration of services with etcd via worker
 - [x] set etcd service registration ttl to 20min
 - [x] stats in abstract supervisor
 - [ ] create generic worker work stream to let supervisor know generic work is done
       this will also tie the worker to the supervisor
 - [ ] easier to use New constructor for abstract supervisor
 - [ ] create SupportSupervisor
 - [ ] collapse supportsupervisor and supervisor (uuid, service tag, mixin)
 - [ ] create worker to continually poll for registration
 - [ ] logging to a service where admin can see it?
       => logging service? Metrics DB? Would need live access to it
 - [ ] http://dave.cheney.net/2014/09/28/using-build-to-switch-between-debug-and-release
 - [ ] rename stats -> metrics
 - [ ] Supervisor settings from etcd
 - [ ] Supervisor write to TaskQueues
 - [ ] Supervisor write stats to InfluxDB
 - [ ] handle git work
 - [ ] save all data to influxdb
 - [ ] fleet test deploy
 - [ ] accept no more work
 - [ ] event when no more claims