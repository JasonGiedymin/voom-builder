package clients

import (
    "encoding/json"
    "log"
)

// == Etcd Key Objects ==
type EtcdSupervisorConfig struct {
    TaskQueue string
}

func ParseSupervisorConfig(jsonData string) EtcdSupervisorConfig {
    var parsedData EtcdSupervisorConfig

    err := json.Unmarshal([]byte(jsonData), &parsedData)
    if err != nil {
        log.Fatalf("Error while trying to parse Etcd config string, data:[%s], error: %v", jsonData, err)
    }

    return parsedData
}
