package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-crawler/jd-products/model"
	"github.com/go-crawler/jd-products/parse"
)

var (
	BaseUrl = "https://search.jd.com/Search?keyword=%E5%A5%B6%E7%B2%89&enc=utf-8&qrst=1&rt=1&stop=1&vt=2&wq=%E5%A5%B6%E7%B2%89&stock=1&page="
)

// 新增数据
func Add(products []parse.Product) {
	for index, product := range products {
		if err := model.DB.Create(&product).Error; err != nil {
			log.Printf("db.Create index: %s, err : %v", index, err)
		}
	}
}

// 导出为XLSX格式
func ExportToExcel(products []parse.Product) {
	xlsx := excelize.NewFile()
	// Set value of a cell.
	for i, product := range products {
		index := strconv.Itoa(i + 1)
		xlsx.SetCellValue("Sheet1", "A"+index, product.ItemId)
		xlsx.SetCellValue("Sheet1", "B"+index, product.Title)
		xlsx.SetCellValue("Sheet1", "C"+index, product.Url)
	}
	// Save xlsx file by the given path.
	err := xlsx.SaveAs("./products.xlsx")
	if err != nil {
		fmt.Println(err)
	}
}

// 开始爬取
func Start() {
	var products []parse.Product

	pages := parse.GetPages(BaseUrl)
	for _, page := range pages {
		client := &http.Client{}
		request, err := http.NewRequest("GET", page.Url, nil)

		request.Header.Add("authority", "search.jd.com")
		request.Header.Add("method", "GET")
		request.Header.Add("path", "/s_new.php?keyword=%E5%A5%B6%E7%B2%89&enc=utf-8&qrst=1&rt=1&stop=1&vt=2&wq=%E5%A5%B6%E7%B2%89&stock=1&page=2&s=29&scrolling=y&log_id=1548672772.95744&tpl=1_M&show_items=6727645,1216716,1805141,831721,2058341,2571565,2165630,6514056,6514054,422539,3455692,725211,1805138,6514034,1950749,3849865,270669,1431731,4264348,1462644,2828306,4189699,5224868,2785708,804013,6727647,2563427,1365864,4007909,3335038")
		request.Header.Add("scheme", "https")
		request.Header.Add("referer", "https://search.jd.com/Search?keyword=%E5%A5%B6%E7%B2%89&enc=utf-8&wq=%E5%A5%B6%E7%B2%89&pvid=f77576002bb5452f90fd187f21f2b3f5")
		request.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")
		request.Header.Add("x-requested-with", "XMLHttpRequest")
		request.Header.Add("Cookie", "pinId=WRpDDyok09iy7T9bbrq9wA; 3AB9D23F7A4B3C9B=PEF4ZBOT5LT7YVZL4LGF6WQE2HHVKDOQ2SSPQYEP5R5TLB3KMWLGRI462MJ2RM5DVMG2ZUB5SS4VV2C256VEZ3RTHI; shshshfpa=0d991cfd-69a0-34c9-679c-68e5faada909-1548294856; __jdv=122270672|www.google.com|-|referral|-|1548294857497; areaId=1; shshshfpb=tNukS30SgEQRz5BFID%2FcUZQ%3D%3D; _gcl_au=1.1.733008858.1548294859; user-key=e8c7f93f-8078-4cf5-95c0-db27e0dcba85; cn=2; wlfstk_smdl=kkril7myi8tc0eikj49rjyrvi0ehkqn5; TrackID=1d3qHZsr8sa7ap0bx0H_jOQwmg-yCC7spgqvsjdW7DroFaPQXLYRdZkLM24nQPKjpGtSbC0A2VSDHY4WEKe-mmJ1Zm8mghRA7z5JVbu3K47ZSwnBctTL3uuniZMPBn7Rl; pin=dingxiaohu2009; unick=dingxiaohu2009; ceshi3.com=103; _tp=8WNRZtSwp4WIDsceuZDCkQ%3D%3D; _pst=dingxiaohu2009; ipLoc-djd=2-2822-51979-0.717816615; ipLocation=%u4e0a%u6d77; __tak=fb1f7cab86ad3dea183ee391df6778fb08963eb95d1aa5609c437b9b59fedee0d7c546827c92daec31d67572b9df147039639909b9408d8b33dc824573c665c9b5bc4eaf636399382a75b3954deb8c76; xtest=9696.cf6b6759; rkv=V0500; qrsc=3; __jda=122270672.1548294857496990674217.1548294857.1548660625.1548665183.3; __jdc=122270672; mt_xid=V2_52007VwMWV1RbWl8WTR9ZDWALGlNaXldZH08pDAJiAEZWCV5OWB1BEUAAblEQTlVfV1sDQU1aVTBRRQEKCFNcL0oYXA17AhJOXlFDWhZCG1QOZgMiUG1YYlkeTR1eA2QDGmJdXVRd; shshshfp=c55439e8a3a560121c9672696dcb6c6a; __jdb=122270672.27.1548294857496990674217|3.1548665183; shshshsID=710b94e8c75619cff655ca511b469195_26_1548672774622")

		if err != nil {
			panic(err)
		}
		//处理返回结果
		response, _ := client.Do(request)
		doc, err := goquery.NewDocumentFromReader(response.Body)
		defer response.Body.Close()

		if err != nil {
			log.Println(err)
		}
		products = append(products, parse.ParseProducts(doc)...)
	}
	//Add(products)

	ExportToExcel(products)

}

func main() {
	Start()

	defer model.DB.Close()
}
