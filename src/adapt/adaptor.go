package adapt

import (
    "fmt"
    "net/http"
    "config"
    "scm"
)

type Adaptor interface {
    Adapt(*http.Request, *config.Config) (*scm.Push, error)
}

type UnknownRepoError struct {
    msg         string
}

func (e UnknownRepoError) Error() string {
    return fmt.Sprintf("adapt: cannot pull from %s", e.msg)
}
