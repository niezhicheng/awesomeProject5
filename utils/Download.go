package utils

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func HttpDownload(url string)  {
	imgPath := "file/"
	fileName := path.Base(url)


	res, err := http.Get(url)
	if err != nil {
		fmt.Println("A error occurred!")
		return
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32 * 1024)


	file, err := os.Create(imgPath + fileName + ".jpg")
	if err != nil {
		panic(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	written, _ := io.Copy(writer, reader)
	fmt.Printf("Total length: %d", written)
}
