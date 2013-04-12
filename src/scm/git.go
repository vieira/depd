package scm

import (
    "os"
    "log"
    "os/exec"
    "bytes"
)

type Git struct {}

func (g Git) Pull(r *Push) (*bytes.Buffer, error) {
    var out bytes.Buffer
    log.SetOutput(&out)
    log.Println("Receiving push...")
    log.SetOutput(os.Stdout)

    cmd := exec.Command("git", "reset", "--hard")
    cmd.Dir = r.Repository.Path
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        return nil, err
    }

    cmd = exec.Command("git", "pull", "origin", r.Ref)
    cmd.Dir = r.Repository.Path
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        return nil, err
    }

    cmd = exec.Command("git", "submodule", "update", "--init")
    cmd.Dir = r.Repository.Path
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        return nil, err
    }
    out.WriteString("\n")
    return &out, nil
}
