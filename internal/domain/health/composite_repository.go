package health

import "context"

type compositeRepository struct {
	repositories []Repository
}

func NewCompositeRepository(repos ...Repository) Repository {
	return &compositeRepository{repositories: repos}
}

func (r *compositeRepository) CheckDB(ctx context.Context) error {
	for _, repo := range r.repositories {
		if err := repo.CheckDB(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *compositeRepository) CheckRedis(ctx context.Context) error {
	for _, repo := range r.repositories {
		if err := repo.CheckRedis(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *compositeRepository) CheckRabbitMQ(ctx context.Context) error {
	for _, repo := range r.repositories {
		if err := repo.CheckRabbitMQ(ctx); err != nil {
			return err
		}
	}
	return nil
}
