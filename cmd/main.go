package main

import (
	"fmt"

	"github.com/KazuyaMatsunaga/Go-VideoGameInformation-Scraping/pkg/repository"
	"github.com/KazuyaMatsunaga/Go-VideoGameInformation-Scraping/pkg/service"

	PkgRepo "github.com/KazuyaMatsunaga/Go-VideoGame-Package-Search/pkg/repository"
	PkgSvc "github.com/KazuyaMatsunaga/Go-VideoGame-Package-Search/pkg/service"
)

func main() {
	repo := repository.NewDetailClient()
	s := service.NewDetailService(repo)

	details := s.Detail()

	repoPkg := PkgRepo.NewPackageClient()
	sP := PkgSvc.NewPackageService(repoPkg)

	fmt.Printf("%v\n", sP.Package(details))
}