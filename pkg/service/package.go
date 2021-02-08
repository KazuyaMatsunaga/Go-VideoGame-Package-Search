package service

import (
	"log"
	PkgRepo "github.com/KazuyaMatsunaga/Go-VideoGame-Package-Search/pkg/repository"

	"github.com/KazuyaMatsunaga/Go-VideoGameInformation-Scraping/pkg/model"
)

type PackageService struct {
	repo PkgRepo.SearchRepository
}

func NewPackageService(repo PkgRepo.SearchRepository) *PackageService {
	return &PackageService {
		repo: repo,
	}
}

func (sP *PackageService) Package(datails []model.Detail) []model.Detail {
	detailData, errorList := sP.repo.Search(datails)
	if len(errorList) != 0 {
		for _, err := range errorList {
			log.Println(err)
		}
	}

	return detailData.([]model.Detail)
}