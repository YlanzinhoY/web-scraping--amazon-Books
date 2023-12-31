package scraping

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Book struct {
	Title string `json:"title"`
	Autor string `json:"author"`
	Price string `json:"price"`
	Link  string `json:"link"`
}

var books []Book

func Scraping() {

	pages := []string{
		"https://www.amazon.com.br/gp/most-wished-for/books/7842670011/ref=zg_mw_pg_1_books?ie=UTF8&pg=1",
		"https://www.amazon.com.br/gp/most-wished-for/books/7842670011/ref=zg_mw_pg_1_books?ie=UTF8&pg=2",
	}

	c := colly.NewCollector()

	c.OnHTML("div.a-column.a-span12.a-text-center._cDEzb_grid-column_2hIsc", func(h *colly.HTMLElement) {
		link := "amazon.com.br" + h.ChildAttr("a.a-link-normal", "href")
		title := h.ChildText("span div._cDEzb_p13n-sc-css-line-clamp-1_1Fn1y")
		price := h.ChildText("span._cDEzb_p13n-sc-price_3mJ9Z")
		author := h.ChildText("div.a-row.a-size-small div._cDEzb_p13n-sc-css-line-clamp-1_1Fn1y")

		if title != "" {
			books = append(books, Book{
				Title: title,
				Autor: author,
				Price: price,
				Link:  link,
			})
		}

		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println(err)
			}
		}()

		index, err := f.NewSheet("Books")
		if err != nil {
			fmt.Println(err)
			return
		}

		f.SetCellValue("Livros", "A1", "Titulo")
		f.SetCellValue("Livros", "B1", "Autor")
		f.SetCellValue("Livros", "C1", "Preco")
		f.SetCellValue("Livros", "D1", "Link")
		// Set active sheet of the workbook.
		f.SetActiveSheet(index)

		for i, values := range books {

			row := i + 2

			f.SetCellValue("Books", fmt.Sprintf("A%d", row), values.Title)
			f.SetCellValue("Books", fmt.Sprintf("B%d", row), values.Autor)
			f.SetCellValue("Books", fmt.Sprintf("C%d", row), values.Price)
			f.SetCellValue("Books", fmt.Sprintf("D%d", row), values.Link)
		}

		if err := f.SaveAs("Books.xlsx"); err != nil {
			fmt.Println(err)
		}

	})

	c.OnHTML("div._cDEzb_p13n-sc-css-line-clamp-2_EWgCb", func(h *colly.HTMLElement) {
		title := h.Text
		if title != "" {
			books = append(books, Book{
				Title: title,
			})
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	for _, url := range pages {
		c.Visit(url)
	}

	createJson()
}

func createJson() {
	data, err := json.MarshalIndent(books, " ", "")

	if err != nil {
		log.Fatal()
	}

	f, _ := os.Create("books.json")

	_, err = f.WriteString(string(data))

	if err != nil {
		panic("error")
	}
}
