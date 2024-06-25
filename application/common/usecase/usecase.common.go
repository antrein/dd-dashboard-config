package usecase

import (
	"antrein/dd-dashboard-config/application/common/repository"
	"antrein/dd-dashboard-config/internal/usecase/configuration"
	"antrein/dd-dashboard-config/internal/usecase/project"
	"antrein/dd-dashboard-config/model/config"
)

type CommonUsecase struct {
	ProjectUsecase *project.Usecase
	ConfigUsecase  *configuration.Usecase
}

func NewCommonUsecase(cfg *config.Config, repo *repository.CommonRepository) (*CommonUsecase, error) {
	configUsecase := configuration.New(cfg, repo.ConfigRepo, repo.InfraRepo)
	projectUsecase := project.New(cfg, repo.ProjectRepo, repo.InfraRepo)

	commonUC := CommonUsecase{
		ProjectUsecase: projectUsecase,
		ConfigUsecase:  configUsecase,
	}
	return &commonUC, nil
}
