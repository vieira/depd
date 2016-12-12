package config

import (
    "encoding/json"
    "os"
    "os/signal"
    "syscall"
    "strings"
    "flag"
    "log"
    "fmt"
    "net"
)

func (c *Config) ListenAndConfigure(filename string) {
    sock, err := net.Listen("unix", "/tmp/depd.sock")
    if err != nil {
        log.Panicln("ListenAndConfigure:", err)
    }

    s := make(chan os.Signal, 1)
    signal.Notify(s, os.Interrupt, os.Kill, syscall.SIGTERM)
    go func() {
        <-s
        sock.Close()
        os.Exit(0)
    }()

    for {
        conn, err := sock.Accept()

        if err != nil {
            log.Println("ListenAndConfigure:", err)
        }

        c.handleConnection(conn, filename)
    }
}

func (c *Config) handleConnection(conn net.Conn, filename string) {
    for {
        conn.Write([]byte("> "))
        buf := make([]byte, 512)
        sz, err := conn.Read(buf)
        cmd := string(buf[0:sz])
        if err != nil {
            log.Println("handleConnection:", err)
            conn.Close()
            return
        }
        c.handleCommand(cmd, filename, conn)
    }
}

func (c *Config) handleCommand(cmd, filename string, conn net.Conn) {
    if argv := strings.Fields(cmd); len(argv) > 1 {
        user := argv[0]
        op := argv[1]
        switch op {
        case "add":
            flags := flag.NewFlagSet(op, flag.ContinueOnError)
            var url = flags.String(
                "url",
                "",
                "url of the repository (required)",
            )
            var ref = flags.String(
                "ref",
                "refs/heads/master",
                "branch deploy",
            )
            var path = flags.String(
                "path",
                "app",
                "path relative to home",
            )
            var email = flags.String(
                "email",
                user + "@yubo.be",
                "notify",
            )

            flags.SetOutput(conn)
            if err := flags.Parse(argv[2:]); err != nil {
                return
            }

            if *url == "" {
                flags.PrintDefaults()
                return
            }

            c.AddRepo(user, *url, *ref, *path, *email)
            c.Write(filename)
        case "list":
            for url, repo := range c.Repositories {
                if repo.Owner.Name == user {
                    b, err := json.MarshalIndent(repo, "", "  ")
                    if err != nil {
                        fmt.Println("error:", err)
                    }
                    conn.Write(append([]byte(url), '\n'))
                    conn.Write(append(b, '\n'))
                }
            }
        case "remove":
            flags := flag.NewFlagSet(op, flag.ContinueOnError)
            var url = flags.String(
                "url",
                "",
                "url of the repository (required)",
            )

            flags.SetOutput(conn)
            if err := flags.Parse(argv[2:]); err != nil {
                return
            }

            if *url == "" {
                flags.PrintDefaults()
                return
            }
            delete(c.Repositories, *url)
            c.Write(filename)

        default:
            conn.Write([]byte("add | remove | list\n"))
        }
    }
}

func (c *Config) AddRepo(user, url, ref, path, email string) {
    for repoUrl, repo := range c.Repositories {
        if repo.Owner.Name == user && repoUrl == url {
            repo.Ref = ref
            repo.Path = path
            repo.Owner.Email = email
            c.Repositories[url] = repo
            return
        }
    }
    // if new, create and append
    var repo Repository
    var owner Owner
    repo.Ref = ref
    repo.Path = path
    owner.Name = user
    owner.Email = email
    repo.Owner = &owner;
    c.Repositories[url] = &repo
}
