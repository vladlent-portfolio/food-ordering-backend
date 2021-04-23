package order

type Service struct {
	repo *Repository
}

func ProvideService(repo *Repository) *Service {
	return &Service{repo}
}
