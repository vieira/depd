package notify

import (
    "log"
    "bytes"
    "net/http"
    "encoding/json"
    "scm"
)

type Slack struct {
    Webhook string
    Channel string `json:"channel"`
    Username string `json:"username"`
    Icon string `json:"icon_emoji"`
    Text string `json:"text"`
}

func (s *Slack) Notify(r *scm.Push, msg *bytes.Buffer) {
    s.Text = "Deployed revision " + r.After[:7] + " to " + r.Repository.Name

    for _, commit := range r.Commits {
        s.Text += "\n" + commit.Id[:7] + ": " + commit.Message +
            " (" + commit.Author.Email + ")"
    }

    json, err := json.Marshal(&s)

    if err != nil {
        log.Println(err)
        return
    }

    buf := bytes.NewBufferString("payload=")
    buf.Write(json)

    resp, err := http.Post(s.Webhook, "application/x-www-form-urlencoded", buf)

    if err != nil {
        log.Println(err)
        return
    }

    if resp.StatusCode != http.StatusOK {
        log.Println(resp.Status)
    }
}
