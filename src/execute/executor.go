package execute

type Executor interface {
    Execute(string, string) error
    Instance(string) (string, error)
}
