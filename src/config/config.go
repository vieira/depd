package config

import (
    "log"
    "io/ioutil"
    "encoding/json"
)

type Repository struct {
    Ref             string
    Path            string
    Owner           *Owner
}

type Owner struct {
    Name            string
    Email           string
}

type Slack struct {
    Username        string
    Webhook         string
    Icon            string
    Channel         string
}

type Config struct {
    Listen          string
    Slack           *Slack
    Repositories    map[string] Repository
}

func Read(fileName string) *Config {
    var c Config

    if data, err := ioutil.ReadFile(fileName); err != nil {
        log.Panicln(err)
    } else if err := json.Unmarshal(data, &c); err != nil {
        log.Panicln(err)
    }

    return &c
}
