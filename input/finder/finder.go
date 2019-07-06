package finder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Book struct {
	ISBN      string
	Title     string
	Author    string
	Publisher string
}

func (b *Book) ConvertCSV() string {
	return b.ISBN + ",1," + b.Title + "," + b.Author + "," + b.Publisher + ",,,\n"
}

type BookResponse struct {
	Kind       string `json:"kind"`
	TotalItems int    `json:"totalItems"`
	Items      []struct {
		Kind       string `json:"kind"`
		ID         string `json:"id"`
		Etag       string `json:"etag"`
		SelfLink   string `json:"selfLink"`
		VolumeInfo struct {
			Title               string   `json:"title"`
			Authors             []string `json:"authors"`
			Publisher           string   `json:"publisher"`
			PublishedDate       string   `json:"publishedDate"`
			Description         string   `json:"description"`
			IndustryIdentifiers []struct {
				Type       string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"industryIdentifiers"`
			ReadingModes struct {
				Text  bool `json:"text"`
				Image bool `json:"image"`
			} `json:"readingModes"`
			PageCount           int      `json:"pageCount"`
			PrintType           string   `json:"printType"`
			Categories          []string `json:"categories"`
			MaturityRating      string   `json:"maturityRating"`
			AllowAnonLogging    bool     `json:"allowAnonLogging"`
			ContentVersion      string   `json:"contentVersion"`
			PanelizationSummary struct {
				ContainsEpubBubbles  bool `json:"containsEpubBubbles"`
				ContainsImageBubbles bool `json:"containsImageBubbles"`
			} `json:"panelizationSummary"`
			ImageLinks struct {
				SmallThumbnail string `json:"smallThumbnail"`
				Thumbnail      string `json:"thumbnail"`
			} `json:"imageLinks"`
			Language            string `json:"language"`
			PreviewLink         string `json:"previewLink"`
			InfoLink            string `json:"infoLink"`
			CanonicalVolumeLink string `json:"canonicalVolumeLink"`
		} `json:"volumeInfo"`
		SaleInfo struct {
			Country     string `json:"country"`
			Saleability string `json:"saleability"`
			IsEbook     bool   `json:"isEbook"`
			ListPrice   struct {
				Amount       float64 `json:"amount"`
				CurrencyCode string  `json:"currencyCode"`
			} `json:"listPrice"`
			RetailPrice struct {
				Amount       float64 `json:"amount"`
				CurrencyCode string  `json:"currencyCode"`
			} `json:"retailPrice"`
			BuyLink string `json:"buyLink"`
			Offers  []struct {
				FinskyOfferType int `json:"finskyOfferType"`
				ListPrice       struct {
					AmountInMicros int64  `json:"amountInMicros"`
					CurrencyCode   string `json:"currencyCode"`
				} `json:"listPrice"`
				RetailPrice struct {
					AmountInMicros int64  `json:"amountInMicros"`
					CurrencyCode   string `json:"currencyCode"`
				} `json:"retailPrice"`
			} `json:"offers"`
		} `json:"saleInfo"`
		AccessInfo struct {
			Country                string `json:"country"`
			Viewability            string `json:"viewability"`
			Embeddable             bool   `json:"embeddable"`
			PublicDomain           bool   `json:"publicDomain"`
			TextToSpeechPermission string `json:"textToSpeechPermission"`
			Epub                   struct {
				IsAvailable  bool   `json:"isAvailable"`
				AcsTokenLink string `json:"acsTokenLink"`
			} `json:"epub"`
			Pdf struct {
				IsAvailable  bool   `json:"isAvailable"`
				AcsTokenLink string `json:"acsTokenLink"`
			} `json:"pdf"`
			WebReaderLink       string `json:"webReaderLink"`
			AccessViewStatus    string `json:"accessViewStatus"`
			QuoteSharingAllowed bool   `json:"quoteSharingAllowed"`
		} `json:"accessInfo"`
		SearchInfo struct {
			TextSnippet string `json:"textSnippet"`
		} `json:"searchInfo"`
	} `json:"items"`
}

func Find(isbn string) (*Book, error) {
	resp, err := http.Get("https://www.googleapis.com/books/v1/volumes?q=isbn:" + isbn)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	byteSlise, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var bookResponse BookResponse
	json.Unmarshal(byteSlise, &bookResponse)

	book := &Book{}
	var builder strings.Builder

	if len(bookResponse.Items[0].VolumeInfo.Authors) > 1 {
		for _, name := range bookResponse.Items[0].VolumeInfo.Authors {
			builder.WriteString(name + "／")
		}
		book.Author = builder.String()
	} else {
		book.Author = bookResponse.Items[0].VolumeInfo.Authors[0]
	}

	book.ISBN = isbn
	book.Publisher = bookResponse.Items[0].VolumeInfo.Publisher
	book.Title = bookResponse.Items[0].VolumeInfo.Title
	return book, nil
}
