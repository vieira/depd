package scm

import (
    "bytes"
)

type Puller interface {
    Pull(*Push) (*bytes.Buffer, error)
}
