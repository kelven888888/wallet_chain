package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"wallet_chain.com/admin/model"
	api_v1 "wallet_chain.com/api/v1/route"
	"wallet_chain.com/global"
	"wallet_chain.com/middleware"
	statics "wallet_chain.com/public"
	"wallet_chain.com/router"
	"wallet_chain.com/utils/gconv"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

type Option func(engine *gin.Engine)

var Options = []Option{}

// 注册app的路由配置
func Include(opts ...Option) {
	Options = append(Options, opts...)
}

func IninRoute() *gin.Engine {
	Include(api_v1.Routers)
	gin.SetMode(gin.ReleaseMode)
	Router := gin.Default()
	Router.StaticFS(global.SHOP_CONFIG.Local.Path, http.Dir(global.SHOP_CONFIG.Local.Path)) // 为用户头像和文件提供静态地址
	// Router.Use(middleware.LoadTls())  // 打开就能玩https了
	//global.SHOP_LOG.Info("use middleware logger")
	// 跨域
	//Router.Use(middleware.Cors()) // 如需跨域可以打开
	//	global.SHOP_LOG.Info("use middleware cors")
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.SHOP_LOG.Info("register swagger handler")
	//store := cookie.NewStore([]byte("SDF324324324"))
	store, _ := redis.NewStoreWithDB(10, "tcp", fmt.Sprintf("%s", global.SHOP_CONFIG.Redis.Addr), global.SHOP_CONFIG.Redis.Password, "3", []byte("SDF324324324"))

	store.Options(
		sessions.Options{
			MaxAge: 60 * 60 * 24 * 7,
			Path:   "/",
		})

	Router.Use(sessions.Sessions("mysession", store))
	Router.Use(gin.Recovery())
	Router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})

	GinLoadHtml(Router)
	//Router.Use(middleware.TestMiddleware()).Use(middleware.DefaultLimit())

	PublicGroup := Router.Group("")
	{
		router.InitPublicRouter(PublicGroup)
	}
	PrivateGroup := Router.Group("admin")
	PrivateGroup.Use(middleware.AdminMiddleware()).Use(middleware.OperationRecord())
	//PrivateGroup.Use(middleware.OperationRecord())
	{
		router.InitAdminRoute(PrivateGroup)    // 注册功能api路由
		router.InitGroupRoute(PrivateGroup)    // 注册功能api路由
		router.InitIndexRoute(PrivateGroup)    // 注册功能api路由
		router.InitRoleRoute(PrivateGroup)     // 注册功能api路由
		router.InitAdminlogRoute(PrivateGroup) // 注册功能api路由

		router.InitConfigRoute(PrivateGroup) // 注册功能api路由

		router.InitSettingRoute(PrivateGroup) // 设置模块

		router.InitAccesslogRoute(PrivateGroup) // 访问日志模块
		//routeappend

	}
	for _, opt := range Options {
		fmt.Println("6666666666666")
		opt(Router)
	}

	return Router

}

type KV struct {
	Code  string
	Value string
}

func LoadTemplate(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	funcMap := template.FuncMap{
		"StringToLower": func(str string) string {
			return strings.ToLower(str)
		},
		"Str2html": func(str string) template.HTML {
			return template.HTML(str)

		},
		"Formatdatetime": func(str time.Time) string {
			return str.Format("2006-01-02 15:04:05")
		},
		"Formatdatetimepoint": func(str model.LocalTime) any {
			sts := str.String()
			return sts
		},
		"TimeDateFormatYMD": func(str time.Time) string {
			if (str.Format("2006-01-02")) == "0001-01-01" {
				return ""
			}
			return str.Format("2006-01-02")
		},
		"Str2FormatYMD": func(str string) string {
			t, _ := time.Parse(time.RFC3339, str)
			return t.Format("2006-01-02")
		},
		"TradeAcountType": func(str string) string {
			TradeType := map[string]string{
				"2":                    "购买量化或基金",
				"AFUND-TO-TDAIFUND":    "钱包划转到基金",
				"AFUND-TO-TDFUND":      "钱包划转到现货",
				"AFUND-TO-TDQUANTFUND": "钱包划转到量化",

				"B-CMS":                      "做多股票手续费",
				"B-FEE":                      "做多股票费用",
				"S-FEE":                      "做空股票费用",
				"S-CMS":                      "做空股票手续费",
				"QUANTDFUND-TO-TFUND":        "量化赎回",
				"AI_QUANTDFUND-TO-TFUND":     "基金赎回到基金账号",
				"TDFUND-TO-AFUND":            "现货划转到钱包",
				"TDAIFUND-TO-AFUND":          "基金划转到钱包",
				"TDQUANTFUND-TO-AFUND":       "量化划转到钱包",
				"QUANTDFUND-TO-AFUND":        "购买量化",
				"AI_QUANTDFUND-BUYSSS":       "购买基金",
				"Cancel_Order":               "取消订单",
				"Submit_Order":               "下单",
				"Submit_Order_Lock":          "下单锁定",
				"Submit_Order_UnLock":        "撤单释放",
				"Submit_Orderoptions_Unlock": "撤单释放",
				"Options_Premium":            "期权权利金",
				"Options_Loss":               "期权交割亏损",
				"Options_Profit":             "期权交割收益",
				"B-FEE-DIFF":                 "做多补返差价",
				"S-FEE-DIFF":                 "做空补扣差价",
			}
			value, exists := TradeType[str]
			if !exists {
				return str
			} else {
				return value
			}

		},
		"Calprofit": func(amount int64, price float64, ag_price float64, cms float64, dire int) string {
			profit := 0.0
			if dire == 0 {
				profit = -cms
			} else {
				profit = (price-ag_price)*float64(-amount) - cms
			}

			return fmt.Sprintf("%f", profit)
		},
		"StrP2FormatYMD": func(str *string) string {
			if str == nil {
				return "-"
			}
			t, _ := time.Parse(time.RFC3339, *str)
			return t.Format("2006-01-02")
		},
		"TimeDateFormatYM": func(str time.Time) string {
			return str.Format("2006-01")
		},
		"datetime": func() string {
			return time.Now().Format("2006-01-02 15:04:05.00000")
		},
		"eqpoint": func(v int, values int) bool {

			return v == values
		},
		"langdefault": func(values string) string {
			if global.SHOP_CONFIG.System.Version == "NQ" {
				return values
			}
			content := make(map[string]string)
			json.Unmarshal([]byte(values), &content)
			return content["en"]
		},
		"language": func(code string, values string) string {
			if global.SHOP_CONFIG.System.Version == "NQ" {
				return values
			}
			content := make(map[string]string)
			json.Unmarshal([]byte(values), &content)
			return content[code]
		},
		"qiuyu10": func(v int, values int) bool {
			return v%values == 24
		},
		"urlfor": func(endpoint string, values ...interface{}) string {
			str := strings.ReplaceAll(endpoint, ".", "/")
			str = strings.ReplaceAll(str, "Controller", "")
			str = strings.ToLower(str)
			params := make(map[string]string)
			//parm := "?"
			if len(values) > 0 {
				key := ""
				for k, v := range values {
					if k%2 == 0 {
						key = fmt.Sprint(v)
						//parm = parm + key + "?"
					} else {
						params[key] = fmt.Sprint(v)
						//parm = parm + key + "=" + params[key]
					}
				}
			}
			if len(params) == 0 {
				return "../../" + str
			}
			u := "?"
			for k, v := range params {
				u += k + "=" + v + "&"
			}
			urls := strings.TrimRight(u, "&")

			/*parm := "?"
			if len(values) > 0 {

				for k, v := range values {
					parm = parm + fmt.Sprintf("?%sk=%sv", fmt.Sprint(k), fmt.Sprint(v))
				}
			}*/
			return "../../" + str + urls
			/*html = strings.ReplaceAll(html, "&#34;", "")
			return template.HTML(html)*/
		},
		"itoa": func(in interface{}) string {
			if in == nil || in == 0 {
				return ""
			}

			s := fmt.Sprint(in)

			if s == "0" {
				return ""
			}

			return s
		},
		"compare": func(a, b interface{}) (equal bool) {
			equal = false
			if strings.TrimSpace(fmt.Sprintf("%v", a)) == strings.TrimSpace(fmt.Sprintf("%v", b)) {
				equal = true
			}
			return equal
		},
		"str2int": func(s uint) int {
			return gconv.Int(s)
		},
		"map_get": func() string {
			return time.Now().Format("2006-01-02 15:04:05.00000")
		},
		"substr_no_html": func() string {
			return time.Now().Format("2006-01-02 15:04:05.00000")
		},
		"ts2TimeDate": func(value int64) string {
			return time.Unix(value, 0).Format("2006-01-02 15:04:05")
		},
		"Cmp": func(i *int, j int) bool { return *i == j },
		"Map2Json": func(v map[string]any) string {
			if v == nil {
				return "-"
			}
			b, _ := json.Marshal(v)
			return string(b)
		},
		"StringMap2ValueEmpty": func(data map[string]string, key string) bool {
			if v, ok := data[key]; ok {
				if v == "" {
					return false
				}
				return true
			} else {
				return false
			}
		},
		"StringMap2Value": func(data map[string]string, key string) string {
			if v, ok := data[key]; ok {
				return v
			} else {
				return ""
			}
		},
		"ParsteFloat": func(nmu float64) string {
			return fmt.Sprintf("%.4f", nmu)
		},
		"Mul": func(num1 decimal.Decimal, num2 decimal.Decimal) decimal.Decimal {
			return num1.Mul(num2)
		},
	}
	//// 非模板嵌套
	htmls, err := filepath.Glob(templatesDir + "/*.html")
	if err != nil {
		panic(err.Error())
	}
	for _, html := range htmls {
		//r.AddFromGlob(filepath.Base(html), html)
		baseName := filepath.Base(html)

		r.AddFromFilesFuncs(baseName, funcMap, html)

	}
	// 布局模板
	layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	if err != nil {
		panic(err.Error())
	}

	// 嵌套的内容模板
	includes, err := filepath.Glob(templatesDir + "/include/**/*.html")
	if err != nil {
		panic(err.Error())
	}
	// 将主模板，include页面，layout子模板组合成一个完整的html页面
	for _, include := range includes {
		// 文件名称
		//兼容win下目录处理
		include = strings.ReplaceAll(include, "\\", "/")
		baseName := filepath.Base(include)

		files := []string{}

		files = append(files, templatesDir+"/layouts/layout.html", include)

		files = append(files, layouts...)
		//fmt.Println(files)

		r.AddFromFilesFuncs(baseName, funcMap, files...)

	}

	return r
}
func GinLoadHtml(v *gin.Engine) {

	//加载后台模板文件
	//v.LoadHTMLGlob("views/admin/**/*")
	//v.Delims("{[{", "}]}")

	//funcMap := template.FuncMap{
	//	"StringToLower": func(str string) string {
	//		return strings.ToLower(str)
	//	},
	//	"Str2html": func(str string) template.HTML {
	//		return template.HTML(str)
	//
	//	},
	//	"Formatdatetime": func(str time.Time) string {
	//		return str.Format("2006-01-02 15:04:05")
	//	},
	//	"TimeDateFormatYMD": func(str time.Time) string {
	//		return str.Format("2006-01-02")
	//	},
	//	"Str2FormatYMD": func(str string) string {
	//		t, _ := time.Parse(time.RFC3339, str)
	//		return t.Format("2006-01-02")
	//	},
	//	"TradeAcountType": func(str string) string {
	//		TradeType := map[string]string{
	//			"2":                          "购买量化或基金",
	//			"AFUND-TO-TDAIFUND":          "钱包划转到基金",
	//			"AFUND-TO-TDFUND":            "钱包划转到现货",
	//			"B-CMS":                      "做多股票手续费",
	//			"B-FEE":                      "做多股票费用",
	//			"S-FEE":                      "做空股票费用",
	//			"S-CMS":                      "做空股票手续费",
	//			"QUANTDFUND-TO-TFUND":        "量化赎回到钱包",
	//			"AI_QUANTDFUND-TO-TFUND":     "基金赎回到基金账号",
	//			"TDFUND-TO-AFUND":            "现货划转到钱包",
	//			"TDAIFUND-TO-AFUND":          "基金划转到钱包",
	//			"QUANTDFUND-TO-AFUND":        "购买量化",
	//			"AI_QUANTDFUND-BUYSSS":       "购买基金",
	//			"Cancel_Order":               "取消订单",
	//			"Submit_Order":               "下单",
	//			"Submit_Order_Lock":          "下单锁定",
	//			"Submit_Order_UnLock":        "撤单释放",
	//			"Submit_Orderoptions_Unlock": "撤单释放",
	//			"Options_Premium":            "期权权利金",
	//			"Options_Loss":               "期权交割亏损",
	//			"Options_Profit":             "期权交割收益",
	//			"B-FEE-DIFF":                 "做多补返差价",
	//			"S-FEE-DIFF":                 "做空补扣差价",
	//		}
	//		value, exists := TradeType[str]
	//		if !exists {
	//			return str
	//		} else {
	//			return value
	//		}
	//
	//	},
	//	"Calprofit": func(amount int64, price float64, ag_price float64, cms float64, dire int) string {
	//		profit := 0.0
	//		if dire == 0 {
	//			profit = -cms
	//		} else {
	//			profit = (price-ag_price)*float64(-amount) - cms
	//		}
	//
	//		return fmt.Sprintf("%f", profit)
	//	},
	//	"StrP2FormatYMD": func(str *string) string {
	//		if str == nil {
	//			return "-"
	//		}
	//		t, _ := time.Parse(time.RFC3339, *str)
	//		return t.Format("2006-01-02")
	//	},
	//	"TimeDateFormatYM": func(str time.Time) string {
	//		return str.Format("2006-01")
	//	},
	//	"datetime": func() string {
	//		return time.Now().Format("2006-01-02 15:04:05.00000")
	//	},
	//	"eqpoint": func(v int, values int) bool {
	//		return v == values
	//	},
	//	"langdefault": func(values string) string {
	//		if global.SHOP_CONFIG.System.Version == "NQ" {
	//			return values
	//		}
	//		content := make(map[string]string)
	//		json.Unmarshal([]byte(values), &content)
	//		return content["en"]
	//	},
	//	"language": func(code string, values string) string {
	//		if global.SHOP_CONFIG.System.Version == "NQ" {
	//			return values
	//		}
	//		content := make(map[string]string)
	//		json.Unmarshal([]byte(values), &content)
	//		return content[code]
	//	},
	//	"qiuyu10": func(v int, values int) bool {
	//		return v%values == 24
	//	},
	//	"urlfor": func(endpoint string, values ...interface{}) string {
	//		str := strings.ReplaceAll(endpoint, ".", "/")
	//		str = strings.ReplaceAll(str, "Controller", "")
	//		str = strings.ToLower(str)
	//		params := make(map[string]string)
	//		//parm := "?"
	//		if len(values) > 0 {
	//			key := ""
	//			for k, v := range values {
	//				if k%2 == 0 {
	//					key = fmt.Sprint(v)
	//					//parm = parm + key + "?"
	//				} else {
	//					params[key] = fmt.Sprint(v)
	//					//parm = parm + key + "=" + params[key]
	//				}
	//			}
	//		}
	//		if len(params) == 0 {
	//			return "../../" + str
	//		}
	//		u := "?"
	//		for k, v := range params {
	//			u += k + "=" + v + "&"
	//		}
	//		urls := strings.TrimRight(u, "&")
	//
	//		/*parm := "?"
	//		if len(values) > 0 {
	//
	//			for k, v := range values {
	//				parm = parm + fmt.Sprintf("?%sk=%sv", fmt.Sprint(k), fmt.Sprint(v))
	//			}
	//		}*/
	//		return "../../" + str + urls
	//		/*html = strings.ReplaceAll(html, "&#34;", "")
	//		return template.HTML(html)*/
	//	},
	//	"itoa": func(in interface{}) string {
	//		if in == nil || in == 0 {
	//			return ""
	//		}
	//
	//		s := fmt.Sprint(in)
	//
	//		if s == "0" {
	//			return ""
	//		}
	//
	//		return s
	//	},
	//	"compare": func(a, b interface{}) (equal bool) {
	//		equal = false
	//		if strings.TrimSpace(fmt.Sprintf("%v", a)) == strings.TrimSpace(fmt.Sprintf("%v", b)) {
	//			equal = true
	//		}
	//		return equal
	//	},
	//	"str2int": func(s uint) int {
	//		return gconv.Int(s)
	//	},
	//	"map_get": func() string {
	//		return time.Now().Format("2006-01-02 15:04:05.00000")
	//	},
	//	"substr_no_html": func() string {
	//		return time.Now().Format("2006-01-02 15:04:05.00000")
	//	},
	//	"ts2TimeDate": func(value int64) string {
	//		return time.Unix(value, 0).Format("2006-01-02 15:04:05")
	//	},
	//	"Cmp": func(i *int, j int) bool { return *i == j },
	//	"Map2Json": func(v map[string]any) string {
	//		if v == nil {
	//			return "-"
	//		}
	//		b, _ := json.Marshal(v)
	//		return string(b)
	//	},
	//	"StringMap2ValueEmpty": func(data map[string]string, key string) bool {
	//		if v, ok := data[key]; ok {
	//			if v == "" {
	//				return false
	//			}
	//			return true
	//		} else {
	//			return false
	//		}
	//	},
	//	"StringMap2Value": func(data map[string]string, key string) string {
	//		if v, ok := data[key]; ok {
	//			return v
	//		} else {
	//			return ""
	//		}
	//	},
	//	"ParsteFloat": func(nmu float64) string {
	//		return fmt.Sprintf("%.4f", nmu)
	//	},
	//}
	//const (
	//	layoutsDir   = "templates/layouts"
	//	templatesDir = "/views/include"
	//	extension    = "/*.html"
	//)
	//
	//tmplFiles, err := fs.ReadDir(tpl.Templates, templatesDir)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//for _, tmpl := range tmplFiles {
	//	fmt.Println(tmpl, "999999999999999999999999999999999999999999999999999999999")
	//	if tmpl.IsDir() {
	//		continue
	//	}
	//
	//	pt, err := template.ParseFS(tpl.Templates, templatesDir+"/"+tmpl.Name(), layoutsDir+extension)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	//	fmt.Println(pt)
	//}
	////v.SetHTMLTemplate(template.Must(template.New("").Funcs(funcMap).ParseFS(tpl.Templates, "admin/**/**/*.html")))
	//tmplFiles, err = fs.ReadDir(tpl.Templates, "admin/include")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//tmpls := template.New("").Funcs(funcMap)
	//r := multitemplate.NewRenderer()
	//for _, tmpl := range tmplFiles {
	//	if tmpl.IsDir() {
	//		soutpm, _ := fs.ReadDir(tpl.Templates, fmt.Sprintf("admin/include/%s", tmpl.Name()))
	//		//fmt.Println(fmt.Sprintf("admin/include/%s", tmpl))
	//		//fmt.Println(soutpm)
	//		for _, tpls := range soutpm {
	//			if tpls.IsDir() {
	//				continue
	//			}
	//			//fmt.Println(tpls.Name())
	//			//fmt.Println(fmt.Sprintf("admin/include/%s/%s", tmpl.Name(), tpls.Name()))
	//			//tplarr := strings.Split(tpls.Name(), ".")
	//			//fmt.Println(tplarr[1])
	//			//tmpls.Must(template.New(tpls.Name()).Funcs(funcMap).ParseFS(tpl.Templates, fmt.Sprintf("admin/include/%s/%s", tmpl.Name(), tpls.Name()), "admin/layouts/*.html")))
	//			//r.Add(tmpl.Name(), template.Must(template.New("").Funcs(funcMap).ParseFS(tpl.Templates, "admin/include/banner/*.html")))
	//			//r.AddFromFilesFuncs(tmpl.Name(), funcMap, files...)
	//		}
	//
	//		//v.SetHTMLTemplate(template.Must(template.New("").Funcs(funcMap).ParseFS(tpl.Templates, fmt.Sprintf("admin/include/%s", soutpm), "admin/**/*.html")))
	//	}
	//	//fmt.Println(tmpl.Name())
	//}
	//
	//v.HTMLRender = r
	//r.seth
	//tmpls := template.New("").Funcs(funcMap)

	v.HTMLRender = LoadTemplate("./views/admin")
	//v.SetHTMLTemplate(template.Must(template.New("").Funcs(funcMap).ParseFS(tpl.Templates, "admin/include/**/*.html", "admin/layouts/*.html")))
	//v.SetHTMLTemplate(template.Must(template.New("").Funcs(funcMap).ParseFS(tpl.Templates, "admin/**/**/*.html", "admin/**/*.html")))
	//tmpl := template.New("").Funcs(funcMap)
	//tmpl = template.Must(tmpl.ParseFS(tpl.Templates, "/**/**/*.html", "/**/*.html"))
	//v.SetHTMLTemplate(tmpl)
	//

	//加载静态文件
	//v.Static("/resource", "./public/resource")
	//v.StaticFS("/resource/", http.FS(statics.Static))
	st, _ := fs.Sub(statics.Static, "resource")
	v.StaticFS("resource", http.FS(st))

	//v.Static("/static", "./public/resource/static")
	//box := packr.NewBox("/public/resource")
	//v.StaticFS("/resource", box)
}
