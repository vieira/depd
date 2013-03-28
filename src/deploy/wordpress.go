package deploy

import (
    "io"
    "os"
    "log"
    "bytes"
    "strings"
    "io/ioutil"
    "os/exec"
    "net/http"
    "encoding/json"
    "scm"
)

type Wordpress struct {
    data []byte
    Version string
    Address string
}

func (w *Wordpress) Deploy(r *scm.Push, out *bytes.Buffer) {
    log.Println("Wordpress deployment detected")

    if err := json.Unmarshal(w.data, &w); err != nil {
        log.Println(err)
        return
    }
    if !w.isUpdated(r.Repository.Path) {
        err := w.download()
        if err == nil {
            err = w.untar()
        }
        if err == nil {
            w.move(r.Repository.Path)
        }
    } else {
        log.Println("Update not required")
    }
    w.clear(r.Repository.Path)
    log.Println("Deployment completed :)")
}

func (w *Wordpress) Data(data []byte) {
    w.data = data
}

func (w *Wordpress) download() error {
    out, err := os.Create(os.TempDir() + "/wp.tar.gz")
    defer out.Close()
    if err != nil {
        log.Println(err)
        return err
    }
    log.Println("Getting wordpress...")
    resp, err := http.Get("http://wordpress.org/wordpress-"+w.Version+".tar.gz")
    defer resp.Body.Close()
    if err != nil {
        log.Println(err)
        return err
    }
    n, err := io.Copy(out, resp.Body)
    if err != nil {
        log.Println(err)
        return err
    }
    log.Printf("Done, downloaded %2.1fMB", float64(n)/(1024*1024))
    return nil
}

func (w *Wordpress) untar() error {
    cmd := exec.Command("tar", "xvfz", "wp.tar.gz")
    cmd.Dir = os.TempDir()
    err := cmd.Run()
    if err != nil {
        log.Println(err)
    }
    os.Remove(os.TempDir() + "/wp.tar.gz")
    return err
}

func (w *Wordpress) move(path string) {
    tmp := os.TempDir() + "/wordpress"
    path = path[:len(path)-len("/wp-content")]
    if err := os.RemoveAll(tmp + "/wp-content"); err != nil {
        log.Println(err)
        return
    }

    if err := os.Remove(tmp + "/wp-config-sample.php"); err != nil {
        log.Println(err)
        return
    }

    if err := os.Remove(tmp + "/readme.html"); err != nil {
        log.Println(err)
        return
    }

    if err := os.Remove(tmp + "/license.txt"); err != nil {
        log.Println(err)
        return
    }

    if err := os.Rename(path + "/wp-content", tmp + "/wp-content"); err != nil {
        log.Println(err)
        return
    }

    if err := os.Rename(path+"/wp-config.php", tmp+"/wp-config.php"); err!=nil {
        log.Println(err)
        return
    }

    os.Rename(path+"/db-config.php", tmp+"/db-config.php") // hyperdb config

    if err := os.RemoveAll(path); err != nil {
        log.Println(err)
        return
    }

    if err := os.Rename(tmp, path); err != nil {
        log.Println(err)
        return
    }
}

func (w *Wordpress) isUpdated(path string) bool {
    v := "$wp_version = '"
    path = path[:len(path)-len("/wp-content")]
    b, err := ioutil.ReadFile(path + "/wp-includes/version.php")
    if err != nil {
        log.Println(err)
        return true // avoid update
    }
    s := string(b)
    s = s[strings.Index(s, v) + len(v):]
    s = s[:strings.Index(s, "\n") - 2]

    return s == w.Version
}

func (w *Wordpress) clear(path string) {
    resp, err := http.Get("http://" + w.Address + "/wp-content/clear-cache.php")

    if err != nil {
        log.Println(err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        log.Println("Cache successfully invalidated.")
    } else {
        log.Println("Cache was not invalidated.")
    }

    if err := os.Remove(path + "/clear-cache.php"); err != nil {
        log.Println("File clear-cache.php not found in wp-content/")
    }
}
