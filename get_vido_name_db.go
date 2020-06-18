package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func main() {
  var start, end int
  fmt.Println("请输入起始页（ >=1）:")
  fmt.Scan(&start)
  fmt.Println("请输入终止页（>=起始页）:")
  fmt.Scan(&end)
  DoWork(start, end)
}

func DoWork(start, end int) {
	fmt.Printf("准备爬取第%d页到第%d页的网址\n", start, end)
	page := make(chan int)
	for i:= start ; i <= end; i++{
		// 定义一个函数，爬取页面
		go SpiderPape(i - 1, page)
	}
	
	for i := start; i <= end; i++ {
		fmt.Printf("豆瓣电影第%d个页面查取完成\n", <-page)
	}

	
}



func SpiderPape(i int, page chan<- int){
	request, err := http.NewRequest("GET", "https://movie.douban.com/top250?start=" + strconv.Itoa(i  *25) + "&filter=", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := http.Client{}
	//添加请求头，模仿浏览器
	//可以接受什么格式的数据
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	//可以接受的语言
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	request.Header.Add("Cache-Control", "max-age=0")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Host", "movie.douban.com")
	request.Header.Add("Pragma", "no-cache")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.25 Safari/537.36 Core/1.70.3756.400 QQBrowser/10.5.4039.400")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	html := string(htmlBytes)
	//编号
	id := regexp.MustCompile(` <a href="https://movie.douban.com/subject/(.*?)/" class="">`)
	idSlice := id.FindAllStringSubmatch(html, -1)
    //fmt.Println(len(idSlice))

	//片名
	name := regexp.MustCompile(` <img width="100" alt="(.*?)"`)
	nameSlice := name.FindAllStringSubmatch(html, -1)
	//fmt.Println(len(nameSlice))

	//评分
	ratReg := regexp.MustCompile(`<span class="rating_num" property="v:average">(.*?)</span>`)
	ratRegSlice := ratReg.FindAllStringSubmatch(html, -1)
	//fmt.Println(len(ratRegSlice))

	//评价人数,投票vote
	voteReg := regexp.MustCompile(`<span>(.*?)人评价</span>`)
	voteRegSlice := voteReg.FindAllStringSubmatch(html, -1)
	//fmt.Println(len(voteRegSlice))

	//描述
	desc := regexp.MustCompile(` <span class="inq">(.*?)</span>`)
	descSlice := desc.FindAllStringSubmatch(html, -1)
	//fmt.Println(len(descSlice))

	//封面
	image := regexp.MustCompile(`src="(.*?)" class=""`)
	imageSlice := image.FindAllStringSubmatch(html, -1)

	dbWork(idSlice, nameSlice, ratRegSlice, voteRegSlice, descSlice, imageSlice)
	//fmt.Println(len(imageSlice))

	// fmt.Println("电影编号   电影名字   评分   评价人数    描述    封面  ")
	// for i := 0; i < len(nameSlice); i++ {
	// 	fmt.Printf("%s,%s,%s,%s,%s,%s\n",
	// 		idSlice[i][1],
	// 		nameSlice[i][1],
	// 		ratRegSlice[i][1],
	// 		voteRegSlice[i][1],
	// 		descSlice[i][1],
	// 		imageSlice[i][1])
	// }
	//username mysql账号

    
	page <- i + 1
}


func dbWork(idSlice, nameSlice, ratRegSlice, voteRegSlice, descSlice, imageSlice [][]string){
	var userName = "root"
    //password mysql密码
    var passWord = "110120"
    //ip   mysql数据库的IP
    var ip = "127.0.0.1"
    //port  mysql数据库的端口
    var port = "3306"
    // database 需要连接的数据库名称
    var database = "test"
    var dataSourceName string
    dataSourceName =  userName + ":" + passWord +"@tcp(" + ip + ":" + port + ")/"+database+"?charset=utf8"
    databases,err:=sql.Open("mysql",dataSourceName)


	// databases,err:=sql.Open("mysql","root:110120(127.0.0.1:3306)/test?charset=utf8")
	defer databases.Close()
	if err!=nil {
		log.Fatal(err)
		return
	}
	for i:=0; i<len(idSlice); i++ {
		_,err:=databases.Exec("replace into test_movie(id, name, rate, vote, movie_desc, movie_url) " + "values" +" (?,?,?,?,?,?)",
		idSlice[i][1],
		nameSlice[i][1],
		ratRegSlice[i][1],
		voteRegSlice[i][1],
		descSlice[i][1],
		imageSlice[i][1])
		if err!=nil {
			log.Fatal(err)
			break
		}
		//每采一页空一格
		// fmt.Println()
	}

   
}


