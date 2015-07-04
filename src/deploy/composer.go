package deploy

import (
    "scm"
    "log"
    "execute"
    "encoding/json"
)

type Composer struct {
    executor execute.Executor
    data []byte
    Dev bool
    Profile bool
}

func (c *Composer) Deploy(r *scm.Push) {
    log.Println("Composer deployment detected")
    cmd := "composer install --prefer-dist -n -o -d=" + r.Repository.Path

    if err := json.Unmarshal(c.data, &c); err != nil {
        log.Println(err)
        return
    }

    if !c.Dev {
        cmd += " --no-dev"
    }

    if c.Profile {
        cmd += " --profile"
    }

    if cid, err := c.executor.Instance(r.Repository.Owner.Name); err != nil {
        log.Println(err)
        return
    } else {
        c.executor.Execute(cmd, cid)
    }
}

func (c *Composer) Data(data []byte) {
    c.data = data
    c.executor = &execute.Docker {}
}
