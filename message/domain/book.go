package domain

type Book struct {
	ISBN      string
	Title     string
	Author    string
	Publisher string
}

func (b *Book) ToString() string {
	return "ISBN:" + b.ISBN + "\n" + "Title:" + b.Title + "\n" + "Author:" + b.Author + "\n" + "Publisher:" + b.Publisher + "\n"
}
