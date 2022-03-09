package main

import (
	"awesomeProject5/pkg/setting"
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)


type Product struct {
	gorm.Model
	Origina  string
	Uid string
}

type Data struct {
	id string `json:"id"`
}


type Postinfo struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	Url    string `form:"url" json:"url" uri:"url" xml:"url" binding:"required"`
	Token string `form:"token" json:"token" uri:"token" xml:"token" binding:"required"`
	Authorization string `form:"authorization" json:"authorization" uri:"authorization" xml:"authorization" binding:"required"`
	Origina string `form:"origina" json:"origina" uri:"origina" xml:"origina" binding:"required"`
	Out string `form:"out" json:"out" uri:"out" xml:"out" binding:"required"`
}

func init()  {
	setting.Setup()
}


func main()  {
	r := gin.Default()
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
	)
	db, err := gorm.Open(
		mysql.New(
			mysql.Config{
				DSN: dsn,
				DefaultStringSize: 256,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
		})
	fmt.Println(db,err)
	db.AutoMigrate(&Product{})
	if db.Error != nil{
		fmt.Println("不好")
	}
	r.POST("/api/post", func(context *gin.Context) {
		fmt.Println("发机会")
		var form Postinfo
		if err := context.Bind(&form); err != nil{
			context.JSON(http.StatusBadRequest,gin.H{
				"error": err.Error(),
			})
			return
		}

		//文件下载
		imgPath := "file/"
		imgUrl := form.Url

		fileName := path.Base(imgUrl)


		res, err := http.Get(imgUrl)
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

		//重命名
		var filena = fmt.Sprintf(imgPath,fileName + "jpg")
		os.Rename(filena, fmt.Sprintf("file/" + form.Out))



		imagePostURL := "https://upload.qiangwe.com/upload?albumId=" + form.Token
		data := map[string]string{
			"albumId": form.Token,
		}
		client := http.Client{}
		bodyBuf := &bytes.Buffer{}
		bodyWrite := multipart.NewWriter(bodyBuf)
		files, err := os.Open(fmt.Sprintf("file/" + form.Out))
		defer files.Close()
		if err != nil {
			log.Println("err")
		}
		// file 为key
		fileWrite, err := bodyWrite.CreateFormFile("file",  fmt.Sprintf("file/" + form.Out))
		_, err = io.Copy(fileWrite, file)
		if err != nil {
			log.Println("err")
		}
		bodyWrite.Close() //要关闭，会将w.w.boundary刷写到w.writer中
		// 创建请求
		contentType := bodyWrite.FormDataContentType()
		for i,v := range data {
			_ = bodyWrite.WriteField(i,v)
		}
		req, err := http.NewRequest(http.MethodPost, imagePostURL, bodyBuf)
		if err != nil {
			log.Println("err")
		}
		// 设置头
		req.Header.Set("Content-Type", contentType)
		req.Header.Add("Authorization","Bearer " + form.Authorization)
		req.Header.Add("X-Correlation-Id",form.Token)
		req.Header.Add("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:97.0) Gecko/20100101 Firefox/97.0")
		resp, err := client.Do(req)
		if err != nil {
			log.Println("err")
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("err")
		}
		id := gjson.Get(string(b), "id")
		fmt.Println("last name:", id.String())

		product := Product{Origina: form.Origina, Uid: id.String(),}

		result := db.Create(&product) // 通过
		fmt.Println(result)



		context.JSON(200,gin.H{
			"data": "成功",
		})
	})
	r.Run(":8000")
}