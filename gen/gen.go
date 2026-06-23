package gen

import (
	_ "embed"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/global"
	"wallet_chain.com/utils"
)

//go:embed service.tpl
var tplserver string

//go:embed contrl.tpl
var tplctrl string

//go:embed route.tpl
var route string

//go:embed sql.tpl
var sql string

type Gen struct {
}

// go run .\main.go gen AccountTeamActivityLog 奖励日志 155  go run .\main.go gen Optionslog 期权日志 54
// go run .\main.go gen AccountTeamActivityConfig 活动配置 903
// go run .\main.go gen AccountTeamActivityAward 活动奖励 903
// go run .\main.go gen PricePrediction 涨跌预测 118
// go run .\main.go gen NewsMorningPaper 早晚报总结 118
// go run main.go gen WalletChain 充提币渠道 6
// go run main.go gen Goods 产品 54
// go run main.go gen w 玩法配置 54
func (this *Gen) Gener(types string, text string, pid string) {
	var models model.MPlayConfig
	var modelrole, roles, roleinsert model.Role
	global.SHOP_DB.Where("module=?", types).First(&modelrole)
	addroute := true
	if modelrole.Id > 0 {
		addroute = false
	}

	global.SHOP_DB.Unscoped().Model(model.Role{}).Where("module=?", types).Delete(&roles)
	service := strings.Replace(tplserver, "Banner", types, -1)
	file, err := os.Create("admin/service/" + types + ".go")
	if err != nil {

		fmt.Println(err)
	}

	defer file.Close()
	file.WriteString(service)
	fmt.Println("service success")
	ctl := strings.Replace(tplctrl, "Banner", types, -1)
	ctl = strings.Replace(ctl, "banner", strings.ToLower(types), -1)

	file, err = os.Create("admin/controller/" + types + ".go")
	if err != nil {

		fmt.Println(err)
	}
	defer file.Close()
	file.WriteString(ctl)
	fmt.Println("controller success")

	routes := strings.Replace(route, "Banner", types, -1)
	routes = strings.Replace(routes, "banner", strings.ToLower(types), -1)
	file, err = os.Create("router/" + strings.ToLower(types) + ".go")
	if err != nil {

		fmt.Println(err)
	}
	defer file.Close()
	file.WriteString(routes)
	fmt.Println("routerfile success")
	content, err := ioutil.ReadFile("initialize/routes.go")
	contents := string(content)
	contents = strings.Replace(contents, "//routeappend", fmt.Sprintf("router.Init%sRoute(PrivateGroup)  \r\n//routeappend", types), -1)
	if addroute {
		ioutil.WriteFile("initialize/routes.go", []byte(contents), 0777)
	}
	fmt.Println("routerinit success")

	//global.SHOP_DB.Unscoped().Model(model.Role{}).Where("module=?", types).Delete(&roles)
	global.SHOP_DB.Last(&roleinsert)

	sqls := strings.Replace(sql, "Banner", types, -1)
	sqls = strings.Replace(sqls, "{parentid}", strconv.Itoa(int(roleinsert.Id+1)), -1)
	//fmt.Println(text)
	sqls = strings.Replace(sqls, "{text}", text, -1)
	sqls = strings.Replace(sqls, "{pid}", pid, -1)
	sqlexarr := strings.Split(sqls, "\r\n")
	for _, v := range sqlexarr {
		if v != "" {
			err = global.SHOP_DB.Exec(v).Error
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	file, err = os.Create("gen/create" + types + ".tpl")
	if err != nil {

		fmt.Println(err)
	}
	defer file.Close()

	fmt.Println("sql success")

	t := reflect.TypeOf(models) // 获取结构体的reflect.Type

	// 检查类型是否为结构体
	if t.Kind() == reflect.Struct {
		fmt.Println("Type is a struct")
		// 遍历结构体的字段
		header := "<th style=\"width: 30px\"><input type=\"checkbox\" lay-skin=\"primary\" lay-filter=\"all-select\"></th>\n                      "
		content := " {{range .List}}\n                    <tr>\n                        <td><input type=\"checkbox\" lay-skin=\"primary\" value=\"{{.Id}}\" class=\"all-x-select\"></td>\n                        "
		item := ""
		structs := "package model\n\nimport \"time\" \r\ntype " + types + " struct {\r\n"
		search_type_field := ""
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i) // 获取字段的详细信息
			//	fmt.Printf("Field %d: Name=%s Type=%s\n", i, field.Name, field.Type)
			tag1 := field.Tag.Get("comment")
			typess := field.Tag.Get("types")
			ranges := field.Tag.Get("range")
			text := field.Tag.Get("text")
			edit := field.Tag.Get("edit")
			time_format := field.Tag.Get("time_format")
			if time_format != "" {
				structs = structs + field.Name + " " + fmt.Sprintf("%s", field.Type) + " `comment:\"" + tag1 + "\" types:\"" + typess + "\" text:\"" + text + "\" json:\"" + utils.TransformSecondUppercase(field.Name) + "\" form:\"" + utils.TransformSecondUppercase(field.Name) + "\" range:\"" + ranges + "\" edit:\"" + edit + "\" time_format:\"" + time_format + "\" `\r\n"
			} else {
				structs = structs + field.Name + " " + fmt.Sprintf("%s", field.Type) + " `comment:\"" + tag1 + "\" types:\"" + typess + "\" text:\"" + text + "\" json:\"" + utils.TransformSecondUppercase(field.Name) + "\" form:\"" + utils.TransformSecondUppercase(field.Name) + "\" range:\"" + ranges + "\" edit:\"" + edit + "\"  `\r\n"
			}

			if tag1 == "" {

				continue
				//tag1 = "没标签"
			}
			search_type_field = search_type_field + "   <option  value=\"" + utils.TransformSecondUppercase(field.Name) + "\" {{if eq .Search.search_field \"" + utils.TransformSecondUppercase(field.Name) + "\"}}selected{{end}}>" + tag1 + "</option>\r\n"
			header = header + "<th>" + tag1 + "</th>\r\n"

			if fmt.Sprintf("%s", field.Type) == "time.Time" {
				if time_format == "2006-01-02" {
					content = content + "<td>{{TimeDateFormatYMD ." + field.Name + "}}</td>\r\n"
				} else {
					content = content + "<td>{{Formatdatetime ." + field.Name + "}}</td>\r\n"
				}
			} else if typess == "radio" {
				rangesarr := strings.Split(ranges, ",")
				textarr := strings.Split(text, ",")
				fmt.Println(rangesarr, textarr)

				content = content + "<td>"
				for k, v := range rangesarr {
					content = content + "{{if eqpoint ." + field.Name + " " + v + "}}" + textarr[k] + "{{end}}"
				}
				content = content + "</td>\r\n"
				//}
			} else {
				content = content + "<td>{{." + field.Name + "}}</td>\r\n"
			}
			if edit != "0" {
				if typess == "radio" {
					rangesarr := strings.Split(ranges, ",")
					textarr := strings.Split(text, ",")
					item = item + `<div class="layui-form-item">
				<label class="layui-form-label">` + tag1 + `</label>
				<div class="layui-input-inline">
				{{if .IsUpdate}}`
					for k, v := range rangesarr {

						item = item + `<input type="radio" name="` + utils.TransformSecondUppercase(field.Name) + `" value="` + v + `" title="` + textarr[k] + `"  {{if eqpoint .result.` + field.Name + ` ` + v + ` }} checked {{end}}>`
						item = item + "\r\n"
					}
					item = item + `{{else}}`
					for k, v := range rangesarr {
						item = item + `	<input type = "radio" name = "` + utils.TransformSecondUppercase(field.Name) + `" value = "` + v + `" title = "` + textarr[k] + `"  checked>`
						item = item + "\r\n"
					}
					item = item + `
				{{end}}
				</div>
				</div>`
				} else {
					if time_format != "" {
						item = item + "<div class=\"layui-form-item\">\n\t\t\t\t<label for=\"title\" class=\"layui-form-label\"><span class=\"x-red\">*</span>" + tag1 + "</label>\n\t\t\t\t<div class=\"layui-input-block\">\n\t\t\t\t\t<input type=\"text\" id=\"" + utils.TransformSecondUppercase(field.Name) + "\" name=\"" + utils.TransformSecondUppercase(field.Name) + "\" lay-verify=\"required\" autocomplete=\"off\" class=\"layui-input\" value=\"{{TimeDateFormatYMD .result." + field.Name + "}}\">\n\t\t\t\t</div>\n\t\t\t</div>\n"
					} else {
						item = item + "<div class=\"layui-form-item\">\n\t\t\t\t<label for=\"title\" class=\"layui-form-label\"><span class=\"x-red\">*</span>" + tag1 + "</label>\n\t\t\t\t<div class=\"layui-input-block\">\n\t\t\t\t\t<input type=\"text\" id=\"" + utils.TransformSecondUppercase(field.Name) + "\" name=\"" + utils.TransformSecondUppercase(field.Name) + "\" lay-verify=\"required\" autocomplete=\"off\" class=\"layui-input\" value=\"{{.result." + field.Name + "}}\">\n\t\t\t\t</div>\n\t\t\t</div>\n"

					}
				}
			}
		}
		structs = structs + "}"
		file.WriteString(sqls + structs)
		fmt.Println(structs)
		tmpl, err := template.ParseFiles("gen/index.tpl")
		if err != nil {
			panic(err)
		}
		xbox := "<xblock>\n                <button class=\"layui-btn layui-btn-danger\" onclick=\"del_all()\"><i class=\"layui-icon\">&#xe640;</i>批量删除</button>\n                <button class=\"layui-btn\" onclick=\"x_admin_show('添加banner', '{{urlfor \"admin.BannerController.Add\"}}')\"><i class=\"layui-icon\">&#xe608;</i>添加</button>\n                <span class=\"x-right\" style=\"line-height:40px\">共有数据：{{.Count}} 条</span>\n            </xblock>"
		xbox = strings.Replace(xbox, "Banner", types, -1)
		scrip := "\t<script>\n\t\twindow.onload = function() {\n\t\t\tlayui.use(['layer', 'form' ,'laypage', 'form'], function() {\n\t\t\t\t$ = layui.jquery; //jquery\n\t\t\t\tlayer = layui.layer; //弹出层\n                var laypage = layui.laypage; // 分页\n                form=layui.form\n                // 分页\n                laypage.render({\n                    elem: 'page',\n                    count: {{.Count}},\n                limit: {{.Search.limit}},\n                curr:  {{.Search.page}},\n                prev: '<em><</em>',\n                    next: '<em>></em>',\n                    skip: false,\n                    jump: function (obj, first) {\n                    if (first != true) {\n                        var query = $('.x-search-form').serialize();\n                        query += \"&page=\" + obj.curr;\n\n                        //  load_page({{urlfor \"admin.RoleController.Index\"}} + \"?\" + query,\"son\");\n                        window.location.href={{urlfor \"admin.BannerController.Index\"}} + \"?\" + query\n                        //layui.admin.tabsBody(layui.admin.tabsPage.index).find(\".layadmin-iframe\")\n                        // layui.admin.url .tabsBodyChange({{urlfor \"admin.RoleController.Index\"}} + \"?\" + query)\n                    }\n                }\n            });\n\n\t\t\t\tlayer.ready(function() { //为了layer.ext.js加载完毕再执行\n\t\t\t\t\tlayer.photos({\n\t\t\t\t\t\tphotos: '#x-img'\n\t\t\t\t\t\t//,shift: 5 //0-6的选择，指定弹出图片动画类型，默认随机\n\t\t\t\t\t});\n\t\t\t\t});\n\t\t\t});\n            $('#resetBtn').on('click', function(){\n                // 使用form模块的resetField方法来重置表单\n                form.resetField('myForm');\n            });\n\t\t}\n\n\t\t// 批量删除提交\n\t\tfunction del_all() {\n\t\t\tparent.layer.confirm('确认要删除吗？', function(index) {\n                var ids = get_list_ids('all-x-select');\n\t\t\t\t// 发异步删除数据\n\t\t\t\tajax_post({{urlfor \"admin.BannerController.DeleteBatch\"}}, {ids: ids}, false,true, false, true);\n\t\t\t});\n\t\t}\n\n\t\t// 删除\n\t\tfunction del(obj, id, name) {\n            parent.layer.confirm('确认要删除吗？', function(index) {\n\t\t\t\t$(obj).parents(\"tr\").remove();\n\n\t\t\t\t//发异步删除数据\n\t\t\t\tajax_post({{urlfor \"admin.BannerController.Delete\"}}, {id: id, name: name}, reload_page);\n\t\t\t});\n\t\t}\n        $(\"#btn-search\").click(function () {\n            var query = $('.x-search-form').serialize();\n            console.log(query)\n            load_page({{urlfor \"admin.BannerController.Index\"}} + '?' + query);\n        });\n        $(\"#reset\").click(function () {\n            load_page({{urlfor \"admin.BannerController.Index\"}} + '?' );\n        });\n\n    </script>"
		scrip = strings.Replace(scrip, "Banner", types, -1)
		//headercontent := " <xblock>\n                <button class=\"layui-btn layui-btn-danger\" onclick=\"del_all()\"><i class=\"layui-icon\">&#xe640;</i>批量删除</button>\n                <button class=\"layui-btn\" onclick=\"x_admin_show('添加banner', '{{urlfor \"admin.BannerController.Add\"}}')\"><i class=\"layui-icon\">&#xe608;</i>添加</button>\n                <span class=\"x-right\" style=\"line-height:40px\">共有数据：{{.Count}} 条</span>\n            </xblock>\n            <table class=\"layui-table layui-form\">\n                <thead>\n                    <tr>\n                        <th style=\"width: 30px\"><input type=\"checkbox\" lay-skin=\"primary\" lay-filter=\"all-select\"></th>\n                        <th>ID</th>\n                        <th>缩略图</th>\n                        <th>标题</th>\n                        <th>跳转url</th>\n<!--                        <th>广告位</th>-->\n<!--                        <th>链接</th>-->\n<!--                        <th>描述</th>-->\n                        <th>语言</th>\n                        <th>状态</th>\n                        <th>添加时间</th>\n                        <th>操作</th>\n                    </tr>\n                </thead>" // 定义数据
		actionbtn := " <td class=\"td-manage\">\n                            <a href=\"javascript:;\" onclick=\"x_admin_show('编辑群组', {{urlfor \"admin.BannerController.Edit\" \"id\" .Id}})\" class=\"layui-btn layui-btn-xs layui-btn-normal\">\n                            <i class=\"layui-icon\">&#xe642;</i>编辑\n                            </a>\n                            <a class=\"layui-btn layui-btn-xs layui-btn-danger\" href=\"javascript:;\" onclick=\"del(this,'{{.Id}}')\">\n                                <i class=\"layui-icon\">&#xe640;</i>删除\n                            </a>\n                        </td> </tr>\n                {{end}}"
		actionbtn = strings.Replace(actionbtn, "Banner", types, -1)
		content = content + actionbtn
		main := "{{define \"main\"}}\r\n"
		end := "{{end}}\r\n"
		keyword := " <input type=\"text\" name=\"kw\" placeholder=\"请输入关键字\" autocomplete=\"off\" class=\"layui-input\" value=\"{{.Search.kw}}\">"

		data := struct {
			Main            template.HTML
			Header          template.HTML
			Content         template.HTML
			Xbox            template.HTML
			Scrip           template.HTML
			Ends            template.HTML
			Searchtypefield template.HTML
			Keyword         template.HTML
		}{
			Main:            template.HTML(main),
			Header:          template.HTML(header),
			Content:         template.HTML(content),
			Xbox:            template.HTML(xbox),
			Scrip:           template.HTML(scrip),
			Ends:            template.HTML(end),
			Searchtypefield: template.HTML(search_type_field),
			Keyword:         template.HTML(keyword),
		}
		os.Mkdir("views/admin/include/"+strings.ToLower(types)+"/", 0755)
		output := "views/admin/include/" + strings.ToLower(types) + "/" + strings.ToLower(types) + "_index.html"
		outputFile, err := os.Create(output)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("index success ")
		defer outputFile.Close() // 确保在main函数结束时关闭文件
		//err = tmpl.ExecuteTemplate(os.Stdout, "main", data)
		// 创建输出文件
		err = tmpl.ExecuteTemplate(outputFile, "main", data)
		if err != nil {
			fmt.Println(err)
		}

		//// 创建输出文件
		//
		//
		//// 执行模板，并将结果写入到文件中
		//err = tmpl.Execute(outputFile, data)
		//if err != nil {
		//	fmt.Println(err)
		//}

		//form
		tmpl, err = template.ParseFiles("gen/form.tpl")
		if err != nil {
			panic(err)
		}

		scrip = "\t<script>\n\t\twindow.onload = function() {\n\t\t\tlayui.use(['form', 'layer', 'upload'], function() {\n\t\t\t\t$ = layui.jquery;\n\n\t\t\t\tvar form = layui.form\n\t\t\t\t\t,layer = layui.layer;\n  \n\t\t\t\t\n\t\t\t\t// 监听提交\n\t\t\t\tform.on('submit(save)', function(data) {\n\t\t\t\t\tajax_post(\"{{.PostUrl}}\", data.field, load_page, true, true, true);\n\t\t\t\t\treturn false;\n\t\t\t\t});\n\t\t\t}); $(\"#btn-search\").click(function () {\n            var query = $('.x-search-form').serialize();\n            console.log(query)\n            load_page({{urlfor \"admin.BannerController.Index\"}} + '?' + query);\n        });\n\t\t}\n\t</script>"
		btnsub := "<div class=\"layui-form-item\">\n\t\t\t\t<input type=\"hidden\" id=\"id\" name=\"id\" value=\"{{.result.Id}}\">\n\t\t\t\t<button  class=\"layui-btn\" lay-filter=\"save\" lay-submit=\"\">\n\t\t\t\t\t保存\n\t\t\t\t</button>\n\t\t\t</div>"
		datas := struct {
			Main   template.HTML
			Header template.HTML
			Btnsub template.HTML
			Items  template.HTML
			Scrip  template.HTML
			Ends   template.HTML
		}{
			Scrip:  template.HTML(scrip),
			Items:  template.HTML(item),
			Btnsub: template.HTML(btnsub),
			Main:   template.HTML(main),

			Ends: template.HTML(end),
		}
		//err = tmpl.ExecuteTemplate(os.Stdout, "main", datas)
		output = "views/admin/include/" + strings.ToLower(types) + "/" + strings.ToLower(types) + "_form.html"
		outputFile, err = os.Create(output)
		if err != nil {
			fmt.Println(err)
		}
		defer outputFile.Close() // 确保在main函数结束时关闭文件
		// 创建输出文件
		err = tmpl.ExecuteTemplate(outputFile, "main", datas)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("form success ")
		var role []model.Role
		global.SHOP_DB.Find(&role)
		rolestr := ""
		for _, v := range role {
			rolestr = rolestr + fmt.Sprintf("%d", v.Id) + ","
		}
		//var admin model.Group
		global.SHOP_DB.Model(model.Group{}).Where("id=1").Updates(model.Group{
			RoleIds: rolestr,
		})
		//生成路由
	} else {
		fmt.Println("Type is not a struct")
	}

}
