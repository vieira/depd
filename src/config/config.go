package config

import (
    "log"
    "io/ioutil"
    "encoding/json"
)

type Repository struct {
    Url             string
    Ref             string
    Path            string
    Owner           *Owner
}

type Owner struct {
    Name            string
    Email           string
}

type Config struct {
    Listen          string
    Repositories    map[string] Repository
}

func Read(fileName string) *Config {
    var c Config
    jc := struct {
        Listen string
        Repositories []Repository
    }{}

    if data, err := ioutil.ReadFile(fileName); err != nil {
        log.Panicln(err)
    } else if err := json.Unmarshal(data, &jc); err != nil {
        log.Panicln(err)
    }

    c.Listen = jc.Listen
    c.Repositories = make(map[string] Repository, len(jc.Repositories))
    for _, repo := range jc.Repositories {
        c.Repositories[repo.Url] = repo
    }

    return &c
}
