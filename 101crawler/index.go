/*
*
project goal:
1. fetch and render any given website
2. parse and download all image in current site
3. (optional) parse further links, go in and do 2.3. again
*/
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const DownloadBasePath = "./images/"

// 获取网站上爬取的数据
// htmlContent 是上面的 html 页面信息，selector 是我们第一步获取的 selector
func GetHttpHtmlContent(url string, selector string, sel interface{}) (string, error) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	//初始化参数，先传一个空的数据
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, _ := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	// 执行一个空task, 用提前创建Chrome实例
	_ = chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	// 创建一个上下文，超时时间为40s  此时间可做更改  调整等待页面加载时间
	timeoutCtx, toCancel := context.WithTimeout(chromeCtx, 40*time.Second)
	defer toCancel()

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector),
		chromedp.OuterHTML(sel, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return htmlContent, nil
}

// 得到具体的数据
// 就是对上面的 html 进行解析，提取我们想要的数据
// 这个 seletor 是在解析后的页面上，定位要哪部分数据
func GetSpecialData(htmlContent string, selector string) ([]string, error) {
	var output []string
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal(err)
		return output, err
	}

	dom.Find(selector).Each(func(i int, selection *goquery.Selection) {
		if curSrc := selection.AttrOr("src", ""); curSrc != "" {
			output = append(output, curSrc)
		}
	})
	return output, nil
}

func DownloadImageFromUrl(url string) error {
	filepaths := strings.Split(url, "/")
	fileName := filepaths[len(filepaths)-1]

	//Get the response bytes from the url
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("received non 200 response code: %s", url)
	}

	//Create a empty file
	os.Mkdir(DownloadBasePath, os.ModePerm)
	file, err := os.Create(DownloadBasePath + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("failed to copy file: %s -> %s", url, fileName)
	}

	fmt.Printf("success: %s\n", url)

	return nil
}

func main() {
	link := "https://shubibubi.itch.io/cozy-fishing"
	visibleSelector := "#wrapper"
	sel := "document.querySelector('body')"
	selector := "img"

	htmlContent, error := GetHttpHtmlContent(link, visibleSelector, sel)

	if error != nil {
		return
	}

	elem, error := GetSpecialData(htmlContent, selector)

	if error != nil {
		return
	}

	if len(elem) <= 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(elem))

	for _, src := range elem {
		go func(_src string) {
			defer wg.Done()

			log.Printf("new src: %s", _src)
			DownloadImageFromUrl(_src)
		}(src)
	}

	wg.Wait()
	fmt.Println("all elems are finished")
}
