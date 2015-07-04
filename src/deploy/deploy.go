package deploy

import (
    "log"
    "encoding/json"
    "io/ioutil"
    "scm"
)

type Deployers map[string] Deployer

type Deployer interface {
    Deploy(*scm.Push)
    Data([]byte)
}

type manifest struct {
    Type string
    Settings json.RawMessage
}

func (ds Deployers) Deploy(r *scm.Push) {
    var ms []manifest
    file := r.Repository.Path + "/.depd.json"
    if data, err := ioutil.ReadFile(file); err != nil {
        log.Println(err)
    } else if err := json.Unmarshal(data, &ms); err != nil {
        log.Println(err)
    } else {
        for _, m := range ms {
            if d, ok := ds[m.Type]; ok {
                d.Data(m.Settings)
                d.Deploy(r)
            } else {
                log.Println("Skipping unknown deployment type " + m.Type)
            }
        }
    }
}

func (ds Deployers) Data(_ []byte) {
    panic("do not call me")
}
