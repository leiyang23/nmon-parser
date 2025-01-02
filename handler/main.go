package handler

import (
	"encoding/base64"
	"fmt"
	"github.com/leiyang23/nmon-parser/core"
	"github.com/leiyang23/nmon-parser/resource"
	"github.com/leiyang23/nmon-parser/util"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
)

type TplData struct {
	Categories map[string]core.Category
	AAA        []core.AAA
	BBBP       []core.BBBP
	Chart      template.HTML
	CookieFile string
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	// handler 异常，说明无法处理该文件，直接清空掉 Cookie，退回到上传页面
	defer func() {
		if err := recover(); err != nil {
			cookie := &http.Cookie{
				Name:   CookieFile,
				Value:  "",
				MaxAge: -1,
			}
			http.SetCookie(w, cookie)

			// 返回 500
			tpl, _ := template.ParseFS(resource.Resource, "assets/500.html")
			_ = tpl.Execute(w, struct {
				Message string
			}{fmt.Sprintf("can not deal this file: %v", err)})
		}
	}()

	var err error

	tpl, _ := template.ParseFS(resource.Resource, "assets/index.html")
	tplData := TplData{CookieFile: CookieFile}

	category := r.FormValue("category")
	fp, err := r.Cookie(CookieFile)
	if err != nil {
		//slog.Error("read cookie error", "error", err)
		err = tpl.Execute(w, tplData)
		if err != nil {
			slog.Error("generate tpl error", "error", err)
		}
		return
	}

	// 上传的文件通过 Cookie 来传递保存路径
	decoded, err := base64.StdEncoding.DecodeString(fp.Value)
	if err != nil {
		http.Error(w, "file path err"+err.Error(), http.StatusBadRequest)
		return
	}
	srcFile := string(decoded)

	//	首页默认为 AAA
	if category == "" {
		category = "AAA"
	}

	// 首先解析出所有的表头指标
	var categories map[string]core.Category
	categories, err = core.GetAllCategory(srcFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tplData.Categories = categories

	// 再根据点击的指标，解析出对应的数据
	if strings.ToUpper(category) == "AAA" {
		tplData.AAA, err = core.GetAAA(srcFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if strings.ToUpper(category) == "BBBP" {
		tplData.BBBP, err = core.GetBBBP(srcFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		chart, err1 := core.LineChart(srcFile, categories[category])
		if err1 != nil {
			http.Error(w, util.WrapErr(err, "generate chart error").Error(), http.StatusInternalServerError)
			return
		}
		tplData.Chart = template.HTML(chart.Element + chart.Script)
	}

	// 模版渲染
	err = tpl.Execute(w, tplData)
	if err != nil {
		http.Error(w, util.WrapErr(err, "parse tpl error").Error(), http.StatusInternalServerError)
		return
	}
}
