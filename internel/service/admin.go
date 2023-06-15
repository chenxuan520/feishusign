package service

type AdminService struct {
}

var DefaultAdminService *AdminService = nil

func (a *AdminService) AdminLogin(code string) (string, error) {
	//TODO finish it
	//step 0 get user message

	//step 1 judge user part and if root

	//step 2 create jwt

	return "", nil
}

func NewAdminService() *AdminService {
	DefaultAdminService = &AdminService{}
	return DefaultAdminService
}
