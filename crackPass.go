package main
 
import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)
 
var passpath  string= "G:/et_game/pass.txt"   //密码集 路径
var rarpath string= "D:/GoZip/test2.rar"    	// rar 文件路径
 
var password = make(chan string)   //创建管道，接收密码
var isOver = make(chan bool) //判断是否退出
 
func main() {
	go passtxt(passpath)
 
 
Loop:
	for{
		select {
		case rarpass:= <-password :
			// cmdshell(rarpath,rarpass)
			go cmdshell(rarpath,rarpass)
		case  <-time.After(time.Second * time.Duration(1)) :
			break Loop
		case <- isOver:
			break Loop
		}
	}
 
}
 
func cmdshell(rarpath string,pass string){
		cmd := exec.Command("unrar", "e","-p"+pass,rarpath,"D:/GoZip")  //解压出来保存 D/GoZip 上
		_, err := cmd.Output()

		fmt.Println(pass)
		if  err == nil {    
			fmt.Printf("密码为：%s \n",pass)
			isOver<-true  // 成功后退出
		}
}
 
func passtxt(passpath string) {
	fp, _ := os.OpenFile(passpath , os.O_RDONLY, 6)
	defer fp.Close()
 
	// 创建文件的缓存区
	r := bufio.NewReader(fp)
	for {
		pass, err2 := r.ReadBytes('\n')
		if err2 == io.EOF {      //文件末尾
			break
		}
		pass = pass[:len(pass)-2]   // 去除末尾 /n
		password <- string(pass)
	}
}
