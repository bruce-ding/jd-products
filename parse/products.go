package parse

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type Product struct {
	ItemId    string
	Title string
	Url string
}

type Page struct {
	Page int
	Url  string
}

// 获取分页
func GetPages(url string) []Page {
	return ParsePages(url)
}

// 分析分页
func ParsePages(baseUrl string) (pages []Page) {
    for i := 0; i < 3; i++ {
    	pages = append(pages, Page{Page: i + 1, Url: baseUrl + strconv.Itoa(i+1)})
		lastItemsBaseUrl := "https://search.jd.com/s_new.php?keyword=%E5%A5%B6%E7%B2%89&enc=utf-8&qrst=1&rt=1&stop=1&vt=2&wq=%E5%A5%B6%E7%B2%89&stock=1&page="
		timestamp := strconv.FormatInt(time.Now().Unix(),10)
		lastItemsUrl := lastItemsBaseUrl + strconv.Itoa(2*(i+1)) + "&s=" + strconv.Itoa(49*(i+1)-20) + "&scrolling=y&log_id=" + timestamp
		pages = append(pages, Page{Page: i + 1, Url: lastItemsUrl})
	}
    return pages
}

// 分析数据
func ParseProducts(doc *goquery.Document) (products []Product) {
	doc.Find("li.gl-item > div.gl-i-wrap > div.p-name > a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		m, err := regexp.Match(`\/\/item.jd.com\/\d+`, []byte(href))
		if err != nil {
			fmt.Println(err)
		}
		itemUrl := href
		if m {
			itemUrl = "https:" + href
		} else {
			resp, err := http.Get(href)
			if err != nil {
				fmt.Printf("http.Get => %v", err.Error())
			}
			itemUrl = resp.Request.URL.String()
		}

		r, _ := regexp.Compile(`\d+`)
		itemId := r.FindString(itemUrl)

		title := s.Find("em").Eq(0).Text()

		product := Product{
			ItemId: itemId,
			Url: itemUrl,
			Title:  title,
		}

		products = append(products, product)
	})

	return products
}
