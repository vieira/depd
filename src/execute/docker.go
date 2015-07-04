package execute

import (
    "bytes"
    "errors"
    "log"
    "strings"
    "os/exec"
)

type Docker struct {}

func (d *Docker) Execute(cmd string, cid string) error {
    cmdargs := strings.Fields(cmd)
    dx := exec.Command("docker", append([]string{"exec", cid}, cmdargs...)...)
    out, err := dx.CombinedOutput()
    log.Println(strings.TrimSuffix(string(out[:]), "\n"))
    return err
}

func (d *Docker) Instance(image string) (string, error) {
    var cidBuff bytes.Buffer
    var err error
    dps := exec.Command("docker", "ps")
    grp := exec.Command("grep", image)
    cut := exec.Command("cut", "-f1", "-d", " ")
    grp.Stdin, _ = dps.StdoutPipe()
    cut.Stdin, _ = grp.StdoutPipe()
    cut.Stdout = &cidBuff
    _ = cut.Start()
    _ = grp.Start()
    _ = dps.Run()
    _ = grp.Wait()
    _ = cut.Wait()
    cid := strings.TrimSuffix(cidBuff.String(), "\n")
    if cid == "" {
        err = errors.New("container not found")
    }
    return cid, err
}
