package main

import (
    "os"
    "flag"
    "log"
    "net/http"
    "config"
    "notify"
    "scm"
    "adapt"
    "deploy"
)

type handlers map[string] func(http.ResponseWriter, *http.Request)

func (hs handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if f, ok := hs[r.Method]; ok {
        f(w, r)
    } else {
        http.Error(w, "Bow ties are cool.", http.StatusMethodNotAllowed)
    }
}

func protocol(t adapt.Adaptor, p scm.Puller, d deploy.Deployer,
              n notify.Notifier,
              c *config.Config) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parse payload
        m, err := t.Adapt(r, c)
        if err != nil {
            if _, ok := err.(*adapt.UnknownRepoError); ok {
                http.NotFound(w, r)
            } else {
                http.Error(w, err.Error(), http.StatusBadRequest)
            }
            log.Println(err)
            return
        }

        // Pull changes
        if out, err := p.Pull(m); err != nil {
            http.Error(w, err.Error(), http.StatusConflict)
            log.Println(err)
            return
        } else {
            w.WriteHeader(http.StatusOK)
            log.SetOutput(out)
            // Additional deployment steps
            d.Deploy(m)
            log.SetOutput(os.Stdout)
            // Notify owner
            n.Notify(m, out)
        }
    }
}

func main() {
    hostname, _ := os.Hostname()
    var configFile = flag.String("config", "config.json", "configuration file")
    flag.Parse()

    c := config.Read(*configFile)
    p := &scm.Git {}
    n := &notify.Notifiers { &notify.Mail { Hostname: hostname } }
    d := &deploy.Deployers { "Wordpress": &deploy.Wordpress {},
                             "Composer": &deploy.Composer {} }

    if c.Slack != nil {
        *n = append(*n, &notify.Slack { Webhook: c.Slack.Webhook,
                                        Username: c.Slack.Username,
                                        Icon: c.Slack.Icon,
                                        Channel: c.Slack.Channel })
    }

    http.Handle("/g", handlers {"POST": protocol(&adapt.Github {}, p, d, n, c)})
    http.Handle("/b", handlers {"POST": protocol(&adapt.Bitbucket {},p,d,n, c)})
    http.ListenAndServe(c.Listen, nil)
}
