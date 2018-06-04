package v1alpha1

import "fmt"

const (
	deploymentNameFmt     = "%s-tor-daemon"
	serviceNameFmt        = "%s-tor-svc"
	roleNameFmt           = "%s-tor-role"
	serviceAccountNameFmt = "%s-tor-serviceaccount"
)

func (s *OnionServiceSpec) GetVersion() int {
	v := 3
	if s.Version == 2 {
		v = 2
	}
	return v
}

func (s *OnionService) DeploymentName() string {
	return fmt.Sprintf(deploymentNameFmt, s.Name)
}

func (s *OnionService) ServiceName() string {
	return fmt.Sprintf(serviceNameFmt, s.Name)
}

func (s *OnionService) RoleName() string {
	return fmt.Sprintf(roleNameFmt, s.Name)
}

func (s *OnionService) ServiceAccountName() string {
	return fmt.Sprintf(serviceAccountNameFmt, s.Name)
}
