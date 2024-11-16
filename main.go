package main

import (
	"crypto/md5"
	"embed"
	"encoding/hex"
	"fmt"
	"ghserve/conf"
	"ghserve/utils"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed front/file.html
var file_html string

//go:embed front/layui
var f embed.FS

func main() {
	var cfg = conf.Init_config()
	rroute := ""
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	hash := md5.Sum([]byte(utils.Random_str()))
	staticroute := "/" + hex.EncodeToString(hash[:])[:]
	if cfg.RandomRoute {
		rroute = "/" + utils.Random_str()
	}
	engine := gin.Default()
	var GRoute gin.IRoutes
	if cfg.Auth {
		GRoute = engine.Group(rroute).Use(gin.BasicAuth(gin.Accounts{
			cfg.User: cfg.Password,
		}))
	} else {
		GRoute = engine.Group(rroute)
	}
	GRoute.GET("/", func(context *gin.Context) {
		context.Redirect(302, rroute+"/file")
	})
	GRoute.StaticFS(staticroute, http.FS(f))
	GRoute.GET("/file/*filepath", func(context *gin.Context) {
		filepath := context.Request.RequestURI
		if strings.HasPrefix(filepath, rroute+"/file") {
			filepath = filepath[len(rroute)+5:]
			filepath, err := url.QueryUnescape(filepath)
			if err != nil {
				fmt.Printf("解码出错: %s\n", err)
				return
			}
			filepath = filepath[1:]
			fmt.Println(cfg.RootPath + filepath)
			if utils.IsDir(cfg.RootPath + filepath) {
				if !strings.HasSuffix(filepath, "/") && len(filepath) != 0 {
					context.Redirect(302, context.Request.RequestURI+"/")
					return
				}
				a := strings.ReplaceAll(file_html, "{{staticroute}}", rroute+staticroute)
				a = strings.ReplaceAll(a, "{{rroute}}", rroute)
				a = strings.ReplaceAll(a, "{{fileinfo}}", utils.Readdir(cfg.RootPath, filepath, rroute))
				a = strings.ReplaceAll(a, "{{uploadpath}}", "/"+filepath)
				context.Writer.WriteString(a)
			} else {
				context.File(cfg.RootPath + filepath)
			}
		}
	})

	GRoute.GET("/api/delete", func(context *gin.Context) {
		file := context.Query("file")
		fmt.Println(file)
		if !strings.HasPrefix(file, "/") || strings.Contains(file, "/../") {
			context.JSON(500, gin.H{
				"message": "File Path Error!",
			})
			return
		}
		file = file[1:]
		fmt.Println(file)
		err := os.Remove(cfg.RootPath + file)
		if err != nil {
			context.JSON(500, gin.H{
				"message": err,
			})
		} else {
			fmt.Println(rroute + "/file/" + file)
			context.Redirect(302, rroute+"/file/"+file[:strings.LastIndex(file, "/")])
		}
	})

	GRoute.POST("/api/upload", func(context *gin.Context) {
		path := context.Query("path")
		if !strings.HasPrefix(path, "/") || strings.Contains(path, "/../") {
			context.JSON(500, gin.H{
				"message": "File Path Error!",
			})
			return
		}
		path = path[1:]
		file, err := context.FormFile("file")
		if err != nil {
			context.JSON(500, gin.H{
				"msg": err.Error(),
			})
			return
		}
		fmt.Println(cfg.RootPath + path + file.Filename)
		context.SaveUploadedFile(file, cfg.RootPath+path+file.Filename)
		context.JSON(200, gin.H{
			"msg": "上传成功",
		})
	})

	// 直接运行 engine 服务
	fmt.Printf("Server Listening on %s:%s%s\n", cfg.Host, cfg.Port, rroute)
	err := engine.Run(cfg.Host + ":" + cfg.Port)
	if err != nil {
		return
	}
}
