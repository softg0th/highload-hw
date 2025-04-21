package infra

import "storage/internal/repository"

type Infra struct {
	Repo *repository.Repository
	Pool *Pool
}

func NewInfra(repo *repository.Repository, pool *Pool) *Infra {
	return &Infra{
		Repo: repo,
		Pool: pool,
	}
}
