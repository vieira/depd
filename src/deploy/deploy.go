package deploy

import (
    "log"
    "os"
    "encoding/json"
    "io/ioutil"
    "bytes"
    "scm"
)

type Deployers map[string] Deployer

type Deployer interface {
    Deploy(*scm.Push, *bytes.Buffer)
    Data([]byte)
}

type manifest struct {
    Type        string
}

func (ds Deployers) Deploy(r *scm.Push, out *bytes.Buffer) {
    var m manifest
    log.SetOutput(out)
    file := r.Repository.Path + "/.depd.json"
    if data, err := ioutil.ReadFile(file); err != nil {
        log.Println(err)
    } else if err := json.Unmarshal(data, &m); err != nil {
        log.Println(err)
    } else if d, ok := ds[m.Type]; ok {
        d.Data(data)
        d.Deploy(r, out)
    }
    log.SetOutput(os.Stdout)
}

func (ds Deployers) Data(_ []byte) {
    panic("do not call me")
}
