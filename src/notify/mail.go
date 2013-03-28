package notify

import (
    "log"
    "bytes"
    "net/smtp"
    "scm"
)

type Mail struct {
    Hostname string
}

func (m *Mail) Notify(r *scm.Push, msg *bytes.Buffer) {
    var mail bytes.Buffer
    repo, rev := r.Repository.Name, r.After
    owner, addr := r.Repository.Owner.Name, r.Repository.Owner.Email

    c, err := smtp.Dial("localhost:25")

    if err != nil {
        log.Println(err)
        return
    }

    c.Mail("root@" + m.Hostname); c.Rcpt("root@" + m.Hostname); c.Rcpt(addr)

    mail.WriteString("Subject: [" + repo + "] Deployment " + rev[:7] + "\n")
    mail.WriteString("From: Deploy Daemon <root@" + m.Hostname + ">\n")
    mail.WriteString("To: " + owner + " <" + addr + ">\n")

    mail.Write(msg.Bytes())

    wc, err := c.Data()
    if err != nil {
        log.Println(err)
        return
    }
    defer wc.Close()
    if _, err = mail.WriteTo(wc); err != nil {
        log.Println(err)
    }
}
