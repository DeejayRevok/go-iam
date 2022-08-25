package healthcheck

type HealthChecker struct {
	singleHealthCheckers []SingleHealthChecker
}

func (checker *HealthChecker) Check() error {
	for _, singleHealthChecker := range checker.singleHealthCheckers {
		if healthCheckErr := singleHealthChecker.Check(); healthCheckErr != nil {
			return healthCheckErr
		}
	}
	return nil
}

func NewHealthChecker(singleHealthCheckers []SingleHealthChecker) *HealthChecker {
	checker := HealthChecker{
		singleHealthCheckers: singleHealthCheckers,
	}
	return &checker
}
