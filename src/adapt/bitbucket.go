package adapt

import (
    "log"
    "encoding/json"
    "net/http"
    "scm"
    "config"
)

type Bitbucket struct {
    Canon_url       string
    Commits         []struct {
        Branch          string
        Raw_node        string
    }
    Repository      struct {
        Name            string
        Absolute_url    string
    }
}

func (r *Bitbucket) Adapt(p *http.Request, c *config.Config) (*scm.Push, error){
    if err := json.Unmarshal([]byte(p.FormValue("payload")), r); err != nil {
        log.Println(err)
        return nil, err
    }

    var m scm.Push
    path := r.Repository.Absolute_url
    m.Repository = &scm.Repository {
        Url: r.Canon_url + path[:len(path)-1], // remove trailing slash
        Name: r.Repository.Name,
    }

    // search if any commit was made to the branch specified in the config
    if repo, ok := c.Repositories[m.Repository.Url]; ok {
        m.Repository.Path = repo.Path
        m.Repository.Owner = &scm.Owner {
            Name: repo.Owner.Name,
            Email: repo.Owner.Email,
        }
        for _, commit := range r.Commits {
            if commit.Branch == repo.Ref {
                m.Ref = commit.Branch
                m.After = commit.Raw_node
                return &m, nil
            }
        }
    }
    return nil, UnknownRepoError {
        msg: "unknown repo " + path + " or no commit to relevant ref",
    }
}
