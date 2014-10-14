package actors

type Worker interface {
    Name() string
    Serve()
    Stop()
}
