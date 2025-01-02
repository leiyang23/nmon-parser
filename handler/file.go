package handler

import (
	"encoding/base64"
	"fmt"
	"github.com/leiyang23/nmon-parser/core"
	"github.com/leiyang23/nmon-parser/resource"
	"github.com/leiyang23/nmon-parser/util"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

const CookieFile = "uploadedFilePath"

func FileHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	content, err := resource.Resource.ReadFile(strings.TrimPrefix(path, "/"))
	if err != nil {
		http.Error(w, "file not exist", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if strings.ToLower(filepath.Ext(path)) == "xlsx" {
		w.Header().Set("Content-Disposition", "attachment;filename=data.xlsx")
	}

	_, err = w.Write(content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ExportExcelHandler(w http.ResponseWriter, r *http.Request) {

	fp, err := r.Cookie(CookieFile)
	if err != nil || fp == nil {
		http.Error(w, "excelFile not exist", http.StatusBadRequest)
		return
	}

	// 上传的文件通过 Cookie 来传递保存路径
	decoded, err := base64.StdEncoding.DecodeString(fp.Value)
	if err != nil {
		http.Error(w, "excelFile path error "+err.Error(), http.StatusBadRequest)
		return
	}
	srcFile := string(decoded)

	// 解析 nmon 文件，写入 Excel
	startTime := time.Now()
	excelFile, _ := filepath.Abs("./uploads/" + filepath.Base(srcFile) + ".xlsx")
	_ = os.Remove(excelFile)

	AAAData, err := core.GetAAA(srcFile)
	if err != nil {
		http.Error(w, util.WrapErr(err, "save excelFile failed").Error(), http.StatusInternalServerError)
		return
	}
	err = core.AddAAAToExcel(excelFile, AAAData)
	if err != nil {
		http.Error(w, util.WrapErr(err, "save excelFile failed").Error(), http.StatusInternalServerError)
		return
	}

	BBBPData, err := core.GetBBBP(srcFile)
	if err != nil {
		http.Error(w, util.WrapErr(err, "save excelFile failed").Error(), http.StatusInternalServerError)
		return
	}
	err = core.AddBBBPToExcel(excelFile, BBBPData)
	if err != nil {
		http.Error(w, util.WrapErr(err, "save excelFile failed").Error(), http.StatusInternalServerError)
		return
	}

	categories, err := core.GetAllCategory(srcFile)
	if err != nil {
		http.Error(w, util.WrapErr(err, "parse category error").Error(), http.StatusBadRequest)
		return
	}
	data, err := core.GetCategoryAllData(srcFile)
	if err != nil {
		http.Error(w, util.WrapErr(err, "parse category data error").Error(), http.StatusBadRequest)
		return
	}
	err = core.AddToExcel(excelFile, categories, data)
	if err != nil {
		http.Error(w, util.WrapErr(err, "save excel excelFile failed").Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("生成文件耗时：", time.Since(startTime).Seconds())

	// 清理掉大量无用的缓存对象，更快地释放内存
	debug.FreeOSMemory()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", url.QueryEscape(filepath.Base(excelFile))))

	startTime = time.Now()
	f, err := os.OpenFile(excelFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		http.Error(w, util.WrapErr(err, "read excel excelFile error").Error(), http.StatusInternalServerError)
		return
	}
	defer util.Close(f)

	_, err = io.Copy(w, f)
	fmt.Println("传输文件耗时：", time.Since(startTime).Seconds())
	if err != nil {
		http.Error(w, util.WrapErr(err, "trans excelFile error").Error(), http.StatusInternalServerError)
		return
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 检查请求方法
	if r.Method != http.MethodPost {
		http.Error(w, "只允许 POST 请求", http.StatusMethodNotAllowed)
		return
	}

	// 2. 解析表单，最大 200 MB
	if err := r.ParseMultipartForm(1024 * 1024 * 200); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3. 获取上传的文件
	file, header, err := r.FormFile("file") // "file" 是表单中文件字段的名称
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer util.Close(file)

	// 4. 创建保存文件的目录
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. 创建新文件
	localFile := filepath.Join(uploadDir, header.Filename)
	localFile, _ = filepath.Abs(localFile)
	dst, err := os.Create(localFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer util.Close(dst)

	// 6. 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 8. 设置cookie,将文件路径返回给前端
	cookie := &http.Cookie{
		Name:     CookieFile,
		Value:    base64.StdEncoding.EncodeToString([]byte(localFile)),
		Path:     "/",
		HttpOnly: false,    // 设置为false允许JavaScript访问
		MaxAge:   3600 * 6, // cookie有效期1小时
	}
	http.SetCookie(w, cookie)

	// 9. 重定向到首页
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
