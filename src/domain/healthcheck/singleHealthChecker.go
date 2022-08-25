package healthcheck

type SingleHealthChecker interface {
	Check() error
}
