package model

type ResponseRakuten struct {
	Items []struct {
		Item struct {
			PostageFlag    int    `json:"postageFlag"`
			Hardware       string `json:"hardware"`
			SalesDate      string `json:"salesDate"`
			SmallImageURL  string `json:"smallImageUrl"`
			Label          string `json:"label"`
			DiscountPrice  int    `json:"discountPrice"`
			ItemPrice      int    `json:"itemPrice"`
			LimitedFlag    int    `json:"limitedFlag"`
			Title          string `json:"title"`
			BooksGenreID   string `json:"booksGenreId"`
			AffiliateURL   string `json:"affiliateUrl"`
			ItemCaption    string `json:"itemCaption"`
			ListPrice      int    `json:"listPrice"`
			ReviewCount    int    `json:"reviewCount"`
			LargeImageURL  string `json:"largeImageUrl"`
			MakerCode      string `json:"makerCode"`
			Jan            string `json:"jan"`
			MediumImageURL string `json:"mediumImageUrl"`
			ReviewAverage  string `json:"reviewAverage"`
			TitleKana      string `json:"titleKana"`
			DiscountRate   int    `json:"discountRate"`
			Availability   string `json:"availability"`
			ItemURL        string `json:"itemUrl"`
		} `json:"Item"`
	} `json:"Items"`
	PageCount        int           `json:"pageCount"`
	Hits             int           `json:"hits"`
	Last             int           `json:"last"`
	Count            int           `json:"count"`
	Page             int           `json:"page"`
	Carrier          int           `json:"carrier"`
	GenreInformation []interface{} `json:"GenreInformation"`
	First            int           `json:"first"`
}