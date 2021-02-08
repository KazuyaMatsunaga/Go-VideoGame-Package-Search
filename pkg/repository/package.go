package repository

import (
	"os"
	"time"
	"encoding/json"

	"strings"
	"log"
	"net/http"
	"sync"
	"strconv"

	"github.com/KazuyaMatsunaga/Go-VideoGameInformation-Scraping/pkg/model"
	PkgModel "github.com/KazuyaMatsunaga/Go-VideoGame-Package-Search/pkg/model"

	"github.com/joho/godotenv"
)

type PackageClient struct{}

func NewPackageClient() SearchRepository {
	return &PackageClient{}
}

func (c *PackageClient) Search(i interface{}) (interface{}, []error) {
	var results interface{}
	errorList := make([]error, 0)
	
	switch i.(type) {
		case []model.Detail:
			results, errorList = GoruToPkgImgSearch(i.([]model.Detail))
			return results, errorList
		default:
			return nil, errorList
	}
}

func GoruToPkgImgSearch(details []model.Detail) (interface{}, []error) {
	detailsAddPkgImg := make([]model.Detail, 0)
	errorList := make([]error, 0)

	limitChCount := 3

	limitCh := make(chan struct{}, limitChCount)
	defer close(limitCh)

	var wg sync.WaitGroup
	wg.Add(limitChCount)

	detailCh := make(chan []model.Detail, limitChCount)
	defer close(detailCh)
	errorCh := make(chan []error, limitChCount)
	defer close(errorCh)

	// details = details[:10]

	lenDetails := len(details)

	countIndex := lenDetails / limitChCount

	firstIndex := 0
	lastIndex := countIndex

	for i := 0; i < limitChCount; i++ {
		if i != (limitChCount - 1) {
			detailsGoru := details[firstIndex:lastIndex - 1]

			limitCh <- struct{}{}
			go PkgImgSearch(limitCh, &wg, i, detailsGoru, detailCh, errorCh)

			firstIndex += countIndex
			lastIndex += countIndex
		} else {
			detailsGoru := details[firstIndex:lenDetails - 1]

			limitCh <- struct{}{}
			go PkgImgSearch(limitCh, &wg, i, detailsGoru, detailCh, errorCh)
		}
	}
	wg.Wait()

	L_detail:
		for {
			select{
				case detailsRes := <- detailCh:
					detailsAddPkgImg = append(detailsAddPkgImg, detailsRes...)
				default:
					break L_detail
			}
		}

	L_err:
		for {
			select{
				case errorRes := <- errorCh:
					errorList = append(errorList, errorRes...)
				default:
					break L_err
			}
		}

	return detailsAddPkgImg, errorList
}

func PkgImgSearch(limitCh chan struct{}, wg *sync.WaitGroup, count int, details []model.Detail, detailCh chan []model.Detail, errorCh chan []error){
	defer wg.Done()

	sleepTime := 60

	err := godotenv.Load("./.env")
    if err != nil {
        log.Fatal(err)
	}

	detailsAddPkgImg := make([]model.Detail, 0)
	errorList := make([]error, 0)

	rakutenIDEnvStr := []byte{}
	rakutenIDEnvStr = append(rakutenIDEnvStr, "RAKUTEN_APP_ID_"...)
	rakutenIDEnvStr = append(rakutenIDEnvStr, strconv.Itoa(count)...)

	baseURL := "https://app.rakuten.co.jp/services/api/BooksGame/Search/20170404?"

	for i, d := range details {
		if len(d.Price) == 0 || len(d.ReleaseDate) == 0 || len(d.Platform) == 0 || len(d.Genre) == 0 {
			continue
		}
		for _, dP := range d.Platform {
			if dP == "PS4" || dP == "Switch" {
				var pkgImgStrc model.PkgImg

				var itemURLStrc model.ItemURL

				reqURL := []byte(baseURL)
				reqURL = append(reqURL, "applicationId="...)
				reqURL = append(reqURL, os.Getenv(string(rakutenIDEnvStr))...)
				reqURL = append(reqURL, "&booksGenreId="...)
				reqURL = append(reqURL, "006"...)
				reqURL = append(reqURL, "&hardware="...)
				if dP == "PS4" {
					reqURL = append(reqURL, "PS4"...)
				} else if dP == "Switch" {
					reqURL = append(reqURL, "Nintendo%20Switch"...)
				}

				reqURL = append(reqURL, "&title="...)
				keywordQuery := []byte{}
				keywordQuery = append(keywordQuery, d.Title...)
				keywordQueryConvert := strings.Replace(string(keywordQuery), " ", "", -1)
				reqURL = append(reqURL, keywordQueryConvert...)

				log.Println(string(reqURL))
				
				res, err := http.Get(string(reqURL))
				if err != nil {
					log.Println(err)
					errorList = append(errorList, err)
					pkgImgStrc = model.PkgImg{dP, "No Image"}
					itemURLStrc = model.ItemURL{dP, "No Link"}
					d.ItemURL = append(d.ItemURL, itemURLStrc)
					d.PackageImg = append(d.PackageImg, pkgImgStrc)
					continue
				}

				defer res.Body.Close()

				var itemSet PkgModel.ResponseRakuten

				if err := json.NewDecoder(res.Body).Decode(&itemSet); err != nil {
					log.Println(err)
					errorList = append(errorList, err)
					pkgImgStrc = model.PkgImg{dP, "No Image"}
					itemURLStrc = model.ItemURL{dP, "No Link"}
					d.ItemURL = append(d.ItemURL, itemURLStrc)
					d.PackageImg = append(d.PackageImg, pkgImgStrc)
					continue
				}

				if len(itemSet.Items) != 0 {
					if itemSet.Items[0].Item.LargeImageURL != "" {
						imgURL := itemSet.Items[0].Item.LargeImageURL
						imgURL = strings.Replace(imgURL, "?_ex=200x200", "", -1)
						pkgImgStrc = model.PkgImg{dP, imgURL}
					} else {
						log.Println("No Image")
						pkgImgStrc = model.PkgImg{dP, "No Image"}
					}
				} else {
					// log.Println("No Image")
					pkgImgStrc = model.PkgImg{dP, "No Image"}
				}

				if len(itemSet.Items) != 0 {
					if itemSet.Items[0].Item.ItemURL != "" {
						itemURL := itemSet.Items[0].Item.ItemURL
						itemURLStrc = model.ItemURL{dP, itemURL}
					} else {
						itemURLStrc = model.ItemURL{dP, "No Link"}
					}
				} else {
					itemURLStrc = model.ItemURL{dP, "No Link"}
				}

				d.ItemURL = append(d.ItemURL, itemURLStrc)
				d.PackageImg = append(d.PackageImg, pkgImgStrc)
			}
		}

		detailsAddPkgImg = append(detailsAddPkgImg, d)
		log.Printf("Index %v of %v ： Search Sleep %vs\n", i, string(rakutenIDEnvStr), sleepTime)
		time.Sleep(time.Second * time.Duration(sleepTime))
	}

	detailCh <- detailsAddPkgImg
	errorCh <- errorList
	<-limitCh
}