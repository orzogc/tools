package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const baiduURL = "https://hanyu.baidu.com/s?wd=%s&ptype=zici"

// 爬取百度汉语的网页
func crawlBaidu(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

// 获取字的拼音
func getPinyin(docBody *goquery.Selection) string {
	// 字的拼音
	var pinyin string
	docBody.Find("div#pinyin.pronounce").Each(func(i int, s *goquery.Selection) {
		s.Find("b").Each(func(i int, p *goquery.Selection) {
			pinyin += p.Text() + " "
		})
	})

	if len(pinyin) != 0 {
		// 多余的空格是为了 GoldenDict 的排版
		pinyin = "拼音：" + pinyin + "\n \n"
	}

	return pinyin
}

// 获取字词的基本释义
func getBaseDef(docBody *goquery.Selection) string {
	// 字词的基本释义
	baseDefinition := docBody.First().Text()

	if len(baseDefinition) != 0 {
		// 去掉多余的空格
		baseDefinition = strings.Replace(baseDefinition, "  ", "", -1)
		// 去掉多余的空行
		baseDefinition = "基本释义：\n" + regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(strings.TrimSpace(baseDefinition), "\n")
	}

	return baseDefinition
}

// 获取字词的详细释义
func getDetaiDef(docBody *goquery.Selection) string {
	// 字词的详细释义
	var detailDefinition string
	// 多音字不同音的划分
	docBody.Eq(1).Find("dl").Each(func(i int, s *goquery.Selection) {
		s.Children().Each(func(i int, p *goquery.Selection) {
			tag := goquery.NodeName(p)
			if tag == "dt" {
				// 多音字的某个拼音
				detailDefinition += p.Text() + "\n"
			} else if tag == "dd" {
				// 字的释义的不同类型
				p.Children().Each(func(i int, q *goquery.Selection) {
					tag := goquery.NodeName(q)
					if tag == "p" {
						// 字的释义的类型
						detailDefinition += q.Text() + "\n"
					} else if tag == "ol" {
						// 字词的不同释义
						q.Find("li").Each(func(i int, r *goquery.Selection) {
							// 字词不同释义的序数
							detailDefinition += strconv.Itoa(i+1) + "."
							// 字词的具体释义
							r.Find("p").Each(func(i int, o *goquery.Selection) {
								text := o.Text()
								if len(text) != 0 {
									detailDefinition += text + "\n"
								}
							})
						})
						detailDefinition += " \n"
					}
				})
			}
		})
	})

	if len(detailDefinition) != 0 {
		// 去掉多余的空行
		detailDefinition = "\n \n" + "详细释义：\n" + strings.TrimRight(detailDefinition, " \n")
	}

	return detailDefinition
}

// 获取字词的释义
func getWords(url string) string {
	doc := crawlBaidu(url)
	docBody := doc.Find("div#content-panel")

	pinyin := getPinyin(docBody)

	docBody = docBody.Find("div.tab-content")
	baseDefinition := getBaseDef(docBody)
	detailDefinition := getDetaiDef(docBody)

	return pinyin + baseDefinition + detailDefinition
}

func main() {
	var words string
	// 只接受查询一个字词
	if len(os.Args) > 1 {
		words = os.Args[1]
	} else {
		return
	}

	searchURL := fmt.Sprintf(baiduURL, url.QueryEscape(words))
	result := getWords(searchURL)
	fmt.Print(result)
}
