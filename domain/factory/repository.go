package factory

import "github.com/BaianoDeve/aster/domain/repository"

type RepositoryFactory interface {
	CreateTransactionRepository() repository.TransactionRepository
}
