package notify

import (
    "bytes"
    "scm"
)

type Notifiers []Notifier

type Notifier interface {
    Notify(*scm.Push, *bytes.Buffer)
}

func (ns Notifiers) Notify(r *scm.Push, out *bytes.Buffer) {
    for _, n := range ns {
        n.Notify(r, out)
    }
}
