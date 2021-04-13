package dish

type API struct {
	Service *Service
}

func ProvideAPI(s *Service) *API {
	return &API{s}
}
