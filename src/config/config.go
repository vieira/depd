package config

import (
    "os"
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
    Repositories    map[string] *Repository
}

func Read(filename string) *Config {
    var c Config

    if data, err := ioutil.ReadFile(filename); err != nil {
        log.Panicln(err)
    } else if err := json.Unmarshal(data, &c); err != nil {
        log.Panicln(err)
    }

    return &c
}

func (c *Config) Write(filename string) {
    b, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        log.Println("error:", err)
    }
    f, err := os.Create(filename)
    if err != nil {
        panic(err)
    }
    defer func() {
        if err := f.Close(); err != nil {
            panic(err)
        }
    }()

    if _, err := f.Write(b); err != nil {
        panic(err)
    }
}
