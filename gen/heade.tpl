 <tbody id="x-img">
                {{range .List}}
                    <tr>
                        <td><input type="checkbox" lay-skin="primary" value="{{.Id}}" class="all-x-select"></td>
                        <td>{{.Id}}</td>
                        <td><img src="../../{{.Image}}" width="200" /></td>
                        <td>{{.Title}}</td>
                        <td>{{.PointUrl}}</td>
                        <td>{{.Language}}</td>

                        <td data-field="Status" >
                               {{if eqpoint .Status 0}}
                            <div class="layui-unselect layui-form-switch layui-checkbox-disbaled layui-disabled" lay-skin="_switch"><em>禁用</em><i></i></div>
                             {{else}}
                                <div class="layui-unselect layui-form-switch layui-form-onswitch layui-checkbox-disbaled layui-disabled" lay-skin="_switch"><em>启用</em><i></i></div>
 {{end}}

                        </td>

                        <td>{{Formatdatetime .CreateTime}}</td>
                        <td class="td-manage">
                            <a href="javascript:;" onclick="x_admin_show('编辑群组', {{urlfor "admin.BannerController.Edit" "id" .Id}})" class="layui-btn layui-btn-xs layui-btn-normal">
                            <i class="layui-icon">&#xe642;</i>编辑
                            </a>
                            <a class="layui-btn layui-btn-xs layui-btn-danger" href="javascript:;" onclick="del(this,'{{.Id}}', '{{.Title}}')">
                                <i class="layui-icon">&#xe640;</i>删除
                            </a>
                        </td>
                    </tr>
                {{end}}
                </tbody>