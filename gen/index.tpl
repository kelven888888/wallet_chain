{{define "main"}}
{{.Main}}
<body class="x-iframe-body">


<div class="x-body">
        <div class="layui-card">
         <form class="layui-form x-center x-search-form" style="width:80%">
                        <div class="layui-form-pane" style="margin-top: 15px;">
                            <div class="layui-form-item">

                                <div class="layui-inline">
                                    <label class="layui-form-label">关键字</label>
                                    <div class="layui-input-inline">

                                      {{.Keyword}}
                                    </div>
                                </div>
                                <div class="layui-inline"> <label class="layui-form-label">范围</label>
                                    <div class="layui-input-inline">
                                        <select name="search_field" lay-verify="required" lay-search="">
                                            <option value="">直接选择</option>

                                           {{.Searchtypefield}}

                                        </select>
                                    </div>
                                </div>



                                <div class="layui-inline">
                                    <div class="layui-input-inline" style="width:80px">
                                        <button type="button" class="layui-btn" id="btn-search"><i class="layui-icon">&#xe615;</i></button>
                                    </div>
                                </div>
                                  <div class="layui-inline">
                                                            <div class="layui-input-inline" style="width:80px">
                                                                <button type="reset" class="layui-btn layui-btn-primary" id="reset">重置</button>
                                                            </div>
                                                        </div>
                            </div>
                        </div>

                    </form>
            <xblock>
              {{.Xbox}}
            </xblock>
            <table class="layui-table layui-form">
                <thead>
                    <tr>

                       {{.Header}}
                       <td>操作</td>
                    </tr>
                </thead>
                <tbody id="x-img">
             {{.Content}}

                </tbody>
            </table>
            <div id="page"></div>
        </div>
	</div>
{{.Scrip}}
</body>
{{.Ends}}
{{end}}