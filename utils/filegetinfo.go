package utils

import (
	"fmt"
	"os"
)

func Readdir(rootpath string, path string, RandomRoute string) string {
	dir, err := os.Open(rootpath + path) // 使用当前目录
	if err != nil {
		fmt.Println("Error opening directory:", err)
		return `<tr><td><a href="/">/</a></td><td>不可操作</td><td>不可操作</td><td>不可操作</td></tr>`
	}
	defer dir.Close() // 确保在函数结束时关闭目录

	// 读取目录内容
	files, err := dir.ReadDir(-1) // 0 表示读取所有文件和子目录
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return `<tr><td><a href="/">/</a></td><td>不可操作</td><td>不可操作</td><td>不可操作</td></tr>`
	}

	res := `<td><a href="../" style="color: rgb(242, 237, 139);">Parent directory</a></td><td>不可操作</td><td>不可操作</td><td>不可操作</td>`

	// 打印目录中的文件信息
	for _, file := range files {
		fileInfo, _ := os.Stat(rootpath + path + file.Name())
		var fsize string
		var fname = fileInfo.Name()
		if file.IsDir() {
			fname += "/"
		}
		if fileInfo.Size() > 1024*1024*1024 {
			fsize = fmt.Sprintf("%.2fGB", float32(fileInfo.Size())/1024/1024/1024)
		} else if fileInfo.Size() > 1024*1024 {
			fsize = fmt.Sprintf("%.2fMB", float32(fileInfo.Size())/1024/1024)
		} else {
			fsize = fmt.Sprintf("%.2fKB", float32(fileInfo.Size())/1024)
		}
		filetime := fileInfo.ModTime().Format("2006/01/02&emsp;&emsp;&emsp;&emsp;15:04:05")
		a := ""
		if file.IsDir() {
			a = `<tr><td><a style="color: rgb(150, 242, 248);" href="` + RandomRoute + "/file/" + path + fname + `">` + fname + `</a></td><td>` + fsize + `</td><td>` + filetime + `</td><td><a href="` + RandomRoute + `/api/delete?file=/` + path + fname + `" style="color: red;">` + `删除</a>` + `</td></tr>`
		} else {
			a = `<tr><td><a style="color: rgb(210, 136, 253);" href="` + RandomRoute + "/file/" + path + fname + `">` + fname + `</a></td><td>` + fsize + `</td><td>` + filetime + `</td><td><a href="` + RandomRoute + `/api/delete?file=/` + path + fname + `" style="color: red;">` + `删除</a>` + `</td></tr>`
		}
		res += a

	}
	return res
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsDir()
}
