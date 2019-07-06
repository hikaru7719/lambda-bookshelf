package finder

import (
	"encoding/csv"
	"fmt"
	"github.com/hikaru7719/lambda-bookshelf/domain"
	"io"
	"strings"
)

type CSV struct {
	Reader io.ReadCloser
}

func (c *CSV) Find(searchWord string) ([]domain.Book, error) {
	bookSlice := make([]domain.Book, 0, 10)
	csvReader := csv.NewReader(c.Reader)
	record, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, row := range record {
		lowerTitle := strings.ToLower(row[2])
		lowerSearchWord := strings.ToLower(searchWord)
		if strings.Contains(lowerTitle, lowerSearchWord) && row[2] != "" && searchWord != "" {
			newBook := domain.Book{ISBN: row[0], Title: row[2], Author: row[3], Publisher: row[4]}
			bookSlice = append(bookSlice, newBook)
		}
	}
	return bookSlice, nil
}

func (c *CSV) Close() {
	c.Reader.Close()
}
