package clients

import (
    "errors"
    "log"

    "github.com/JasonGiedymin/voom-builder/config"

    "github.com/coreos/go-etcd/etcd"
)

var (
    // errors
    ErrorNoEtcdClient      = errors.New("None or invalid etcd client given.")
    ErrorNoEtcdConnection  = errors.New("Could not contact etcd.")
    ErrorNoEtcdConfigFound = errors.New("Could not find config in etcd.")
    ErrorEtcdPathNotFound  = errors.New("Could not find the services path in etcd.")
)

const (
    // keys
    EtcdKey_ConfigPath = "/app/supervisor/config"
)

type EtcdPair struct {
    Key   string
    Value string
}

// == Etcd Wrapper ==
type Etcd struct { // compose etcd to make key access easier
    client     *etcd.Client
    configFile *config.EtcdConfig
}

// Get Supervisor Config
func (e *Etcd) Config() (*EtcdSupervisorConfig, error) {
    if r, err := e.client.Get(e.configFile.Paths.SupervisorConfig, true, true); err != nil {
        log.Println(ErrorNoEtcdConfigFound.Error() + " Error was: " + err.Error())
        return nil, err
    } else {
        //log.Printf("k:%s, v:%s, i:%d", r.Node.Key, r.Node.Value, r.Node.CreatedIndex)
        configData := ParseSupervisorConfig(r.Node.Value)
        log.Printf("Config: %v", configData)
        return &configData, nil
    }
}

// wrap for convience
func (e *Etcd) Set(data EtcdPair) (string, error) {
    key := data.Key
    value := data.Value

    resp, err := e.client.Set(key, value, e.configFile.Service_ttl)
    if err != nil {
        log.Fatal(err)
        return "", err
    }

    // log.Println(resp.Node)

    return resp.Node.Value, nil
}

func check(etcdClient *etcd.Client) error {
    if etcdClient == nil {
        return ErrorNoEtcdClient
    }

    if cluster := etcdClient.GetCluster(); len(cluster) == 0 {
        return ErrorNoEtcdConnection
    }

    if _, err := etcdClient.CreateDir("supervisors", 1); err != nil {
        return ErrorNoEtcdConnection
    }

    return nil
}

func NewEtcdClient(configFile *config.EtcdConfig) (*Etcd, error) {
    etcdClient := etcd.NewClient([]string{configFile.Location})

    err := check(etcdClient)
    if err != nil {
        return nil, err
    }

    consistency := configFile.Consistency
    if consistency != "" {
        etcdClient.SetConsistency(consistency)
    }

    return &Etcd{client: etcdClient, configFile: configFile}, nil
}
