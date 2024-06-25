package repository

import (
	"antrein/dd-dashboard-config/application/common/resource"
	"antrein/dd-dashboard-config/internal/repository/configuration"
	"antrein/dd-dashboard-config/internal/repository/infra"
	"antrein/dd-dashboard-config/internal/repository/project"
	"antrein/dd-dashboard-config/model/config"
)

type CommonRepository struct {
	ProjectRepo *project.Repository
	ConfigRepo  *configuration.Repository
	InfraRepo   *infra.Repository
}

func NewCommonRepository(cfg *config.Config, rsc *resource.CommonResource) (*CommonRepository, error) {
	projectRepo := project.New(cfg, rsc.Db)
	infraRepo := infra.New(cfg)
	configRepo := configuration.New(cfg, rsc.Db, infraRepo)

	commonRepo := CommonRepository{
		ProjectRepo: projectRepo,
		ConfigRepo:  configRepo,
		InfraRepo:   infraRepo,
	}
	return &commonRepo, nil
}
