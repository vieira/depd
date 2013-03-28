package scm

type Push struct {
    Ref         string
    After       string
    Commits     []Commit
    Repository  *Repository
}

type Commit struct {
    Id          string
    Message     string
    Author      *Owner
}

type Repository struct {
    Name        string
    Url         string
    Path        string
    Owner       *Owner
}

type Owner struct {
    Name        string
    Email       string
}
