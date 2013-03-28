package adapt

import (
    "log"
    "net/http"
    "encoding/json"
    "scm"
    "config"
)

type Github struct {
    After       string
    Ref         string
    Repository  *scm.Repository
    Commits     []scm.Commit
}

func (r *Github) Adapt(p *http.Request, c *config.Config) (*scm.Push, error) {
    if err := json.Unmarshal([]byte(p.FormValue("payload")), r); err != nil {
        log.Println(err)
        return nil, err
    }

    var m scm.Push
    if repo, ok := c.Repositories[r.Repository.Url]; ok && repo.Ref == r.Ref {
        m.Repository = r.Repository
        m.Repository.Path = repo.Path
        m.After = r.After
        m.Ref = r.Ref
        m.Commits = r.Commits
    } else {
        return nil, UnknownRepoError {
            msg: "unknown repo " + r.Repository.Name + " or ref " + r.Ref,
        }
    }
    return &m, nil
}
