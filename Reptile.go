package main

import (
	"fmt"
	"net/http"
	"strconv"
	"regexp"
	"strings"
	"os"
)



func main() {
	var start, end int
	fmt.Println("请输入起始页（ >=1）:")
	fmt.Scan(&start)
	fmt.Println("请输入终止页（>=起始页）:")
	fmt.Scan(&end)

	DoWork(start, end) //工作函数
}

//工作函数
func DoWork(start, end int) {
	fmt.Printf("准备爬取第%d页到第%d页的网址\n", start, end)
	page := make(chan int)
	for i:= start; i <= end; i++{
		// 定义一个函数，爬取页面
		go SpiderPape(i, page)
	}
	
	for i := start; i <= end; i++ {
		fmt.Printf("第%d个页面查取完成\n", <-page)
	}
}


func SpiderPape(i int, page chan<- int){
	//明确爬取的url
	url := "https://www.pengfu.com/xiaohua_" + strconv.Itoa(i) + ".html"
	fmt.Printf("正在爬取第%d个网页：%s\n",i, url)

	// 开始爬取页面内容
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println("HttpGet err = ", err)
		page<- i
		return
	}
	// fmt.Println(result)
	// 解释表达式
	re := regexp.MustCompile(`<h1 class="dp-b"><a href="(?s:(.*?))</a>`)
	if re == nil{
		fmt.Println("regexp.MustCompile err")
		page<- i
		return
	}
	
	joyUrls := re.FindAllStringSubmatch(result, -1)
	fmt.Println("test-----joyUrls---------- = ", joyUrls)
	fileTitle := make([]string, 0)
	fileContent := make([]string, 0)
	//fmt.Println(joyUrls)
	//取网址
	//第一个返回下标,第二个返回内容
	for _, data := range joyUrls {
		// 开始爬取每一个段子
		title, content, err := SpiderOneJoy(data[1])
		if err != nil {
			fmt.Println("SpiderOneJoy err", err)
			continue
		}
		fileTitle = append(fileTitle, title)
		fileContent = append(fileContent,content)

	}

	StoreJoyToFile(i, fileTitle, fileContent)
	page <- i // 写内容
}


func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}

	defer resp.Body.Close()

	// 读取网页内容
	buf := make([]byte, 1024 * 4)
	for{
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		result += string(buf[:n]) // 累加读取的内容
	}
	return
}


func SpiderOneJoy(url string) (title, content string, err error) {
	result, err1 := HttpGet(url)
	if err != nil {
		err = err1
		return
	}

	// 去关键信息
	rel := regexp.MustCompile(`<h1 class="dp-b"><a href="(?s:(.*?))</a>`)
	if rel == nil{
		err = fmt.Errorf("%s", "egexp.MustCompile err")
		return
	}

	// 取内容
	tmpTitle := rel.FindAllStringSubmatch(result, 1)
	for _, data := range  tmpTitle {
		title = data[1]
		//strings.Replace(title, "r", "", -1)
		//strings.Replace(title, "\n", "", -1)
		//strings.Replace(title, " ", "", -1)
		strings.Replace(title, "\t", "", -1)
		break
	}

	//取内容
	re2 := regexp.MustCompile(`<div class="content-img clearfix pt10 relative">(?s:(.*?))</div>`)
	if re2 == nil{
		err = fmt.Errorf("%s", "egexp.MustCompile err2")
		return
	}

	// 取内容
	tmpContent := re2.FindAllStringSubmatch(result, -1)
	for _, data := range  tmpContent {
		content = data[1]
		strings.Replace(content, "r", "", -1)
		strings.Replace(content, "\n", "", -1)
		strings.Replace(content, "<br />", "", -1)
		strings.Replace(content, "\t", "", -1)
		break
	}
	return

}


// 内容写入文件
func StoreJoyToFile(i int, fileTitle, fileContent []string) {
	// 新建文件
	f, err := os.Create(strconv.Itoa(i) + ".txt")
	if err != nil {
		fmt.Println("os.Create err", err)
		return
	}
	defer f.Close()

	//写内容
	n := len(fileTitle)
	for i := 0; i < n; i++{
		// 写标题
		f.WriteString(fileTitle[i] + "\n")
		// 写内容
		f.WriteString(fileContent[i] + "\n")
		f.WriteString("\r\n=================================\r\n")
	}

}












