
layui.config({
    base: '/resource/layuiadmin/' //静态资源所在路径
}).extend({
    index: 'lib/index' //主入口模块
}).use(['index', 'table'], function(){
    var $ = layui.$
        ,form = layui.form
        ,table = layui.table;

    //监听搜索
    form.on('submit(LAY-user-back-search)', function(data){
        var field = data.field;

        //执行重载
        table.reload('LAY-user-back-manage', {
            where: field
        });
    });

    //事件
    var active = {
        batchdel: function(){
            var checkStatus = table.checkStatus('LAY-user-back-manage')

                ,checkData = checkStatus.data; //得到选中的数据

            if(checkData.length === 0){
                return layer.msg('请选择数据');
            }
            var ids = new Array();
            for (var i = 0; i < checkStatus.data.length; i++) {
                ids.push(checkStatus.data[i].id);//循环获取选中的行ID
            }
            layer.prompt({
                formType: 1
                ,title: '敏感操作，请验证口令'
                ,value:"123456"
            }, function(value, index){
                layer.close(index);

                layer.confirm('确定删除吗？', function(index) {
                    var datas =new Array();
                    datas["ids"]=ids;
                    //执行 Ajax 后重载
                    console.log(ids)
                    $.ajax({
                        url: '../admin/deletebatch',
                        type: 'Post',
                        async: true,
                        // data :JSON.stringify(datas),
                        data: JSON.stringify({ 'ids': ids ,"password":value}),
                        dataType: 'json',
                        contentType: 'application/json',
                        timeout: 30000,
                        success: successCallback,
                        error: errorCallback,
                        complete: completeCallback,
                        statusCode: {
                            404: 404,
                            500: 500
                        }

                    });
                    table.reload('LAY-user-back-manage');
                    layer.msg('已删除');
                });
            });
        }
        ,add: function(){
            layer.open({
                type: 2
                ,title: '添加管理员'
                ,content: '../admin/add'
                ,area: ['720px', '520px']
                ,btn: ['确定', '取消']
                ,yes: function(index, t){
                    var iframeWindow = window['layui-layer-iframe'+ index]
                        ,submitID = 'LAY-user-back-submit'
                        ,submit = t.find('iframe').contents().find('#LAY-user-back-submit');
                    // ,submit = layero.find('iframe').contents().find("#LAY-user-role-submit");

                    //监听提交
                    iframeWindow.layui.form.on('submit(LAY-user-back-submit)', function(t){
                        // var field = data.field; //获取提交的字段

                        //提交 Ajax 成功后，静态更新表格中的数据

                        t.field.id = parseInt(t.field.id)
                        t.field.roleid = parseInt(t.field.roleid)
                        t.field.status= parseInt(t.field.status)
                        t.field.mobile= parseInt(t.field.mobile)
                        $.ajax({
                            async: true,
                            url: '../admin/add',
                            type: 'Post',
                            data :JSON.stringify(t.field),
                            dataType: 'json',
                            timeout: 30000,
                            success: successCallback,
                            error: errorCallback,
                            complete: completeCallback,
                            statusCode: {
                                404: 404,
                                500: 500
                            }
                        })


                        table.reload('LAY-user-back-manage'); //数据刷新
                        layer.close(index); //关闭弹层
                    });

                    submit.trigger('click');
                }
            });
        }
    }
    $('.layui-btn.layuiadmin-btn-admin').on('click', function(){
        var type = $(this).data('type');
        active[type] ? active[type].call(this) : '';
    });
});

function successCallback(json) {
    layui.table.reload('LAY-user-back-manage');
}

function errorCallback(xhr, status){
    console.log('出问题了！');
}

function completeCallback(xhr, status) {
    console.log('ok！');
}
/** layuiAdmin.std-v1.4.0 LPPL License By https://www.layui.com/admin/ */
;
layui.define(["table", "form"],
    function(e) {
        var t = layui.$,
            i = layui.table;
        layui.form;
        i.render({
            elem: "#LAY-user-back-manage",
            url: "../admin/getlist",
            cols: [[{
                type: "checkbox",
                fixed: "left"
            },
                {
                    field: "id",
                    width: 80,
                    title: "Id",
                    sort: !0
                },
                {
                    field: "account",
                    sort: !0,
                    title: "登录名"
                },
                {
                    field: "mobile",
                    title: "手机"
                },
                {
                    field: "mail",
                    title: "邮箱"
                },
                {
                    field: "groupname",
                    title: "角色"
                },
                {
                    field: "CreatedAt",
                    title: "加入时间",
                    sort: !0,
                    templet: function(d){
                        return layui.util.toDateString(d.CreatedAt, 'yyyy-MM-dd HH:mm:ss');
                    }
                },
                {
                    field: "Status",
                    title: "审核状态",
                    templet: function (d) {
                        return '<input type="checkbox" name="status" disabled value="' + d.Status + '" lay-skin="switch" lay-text="启用|禁用" lay-filter="status" ' + (d.status == 1 ? 'checked' : '') + '>';

                    },
                    minWidth: 80,
                    align: "center"
                },
                {
                    title: "操作",
                    width: 150,
                    align: "center",
                    fixed: "right",
                    templet: function (d) {
                        return ' <a class="layui-btn layui-btn-normal layui-btn-xs" lay-event="edit"><i class="layui-icon layui-icon-edit"></i>编辑</a>' +
                            '  <a class="layui-btn layui-btn-danger layui-btn-xs" lay-event="del"><i class="layui-icon layui-icon-delete"></i>删除</a>';

                    },
                }]],
            text: "对不起，加载出现异常！",
            page: !0,
            limit: 15,
            limits: [10, 15, 20, 25, 30],
            done : function(res, curr, count)
            {

                if (res.count < curr*limit)
                {
                    $(".layui-table-main").html('<div class="layui-none">暂无数据</div>');
                }
            }
        }),
            i.on("tool(LAY-user-back-manage)",
                function(e) {
                    e.data;
                    if ("del" === e.event) layer.prompt({
                            formType: 1,
                            title: "敏感操作，请验证口令",
                            value:'123456' +
                                '',
                            btn:['提交','取消'],
                            /* yes: function (index, layero) {
                                 var value1 = $('layui-layer-input').val();//获取多行文本框的值
                                 alert('您刚才输入了:' + value1);}*/

                        },
                        function(t, i) {
                            var value1 = layui.$('.layui-layer-input').val();//获取多行文本框的值

                            layer.close(i),
                                layer.confirm("确定删除此管理员？",
                                    function(t) {
                                        //  console.log('您刚才输入了:' + value1);
                                        var datas= {id:e.data.id, password: value1};
                                        layui.$.ajax({
                                            async: true,
                                            url: '../admin/delete',
                                            type: 'Post',
                                            data :JSON.stringify(datas),
                                            // dataType: 'json',
                                            contentType: 'application/json',
                                            timeout: 30000,
                                            success: successCallback,
                                            error: errorCallback,
                                            complete: completeCallback,
                                            statusCode: {
                                                404: 404,
                                                500: 500
                                            }
                                        })
                                        console.log(e),
                                            e.del(),
                                            layer.close(t)
                                    })
                        });
                    else if ("edit" === e.event) {
                        t(e.tr);
                        layer.open({
                            type: 2,
                            title: "编辑管理员",
                            content: "../admin/edit?id="+e.data.id,
                            area: ["720px", "520px"],
                            btn: ["确定", "取消"],
                            yes: function(e, t) {


                                var l = window["layui-layer-iframe" + e],
                                    // r = "LAY-user-back-submit",
                                    n = t.find("iframe").contents().find("#LAY-user-back-submit");
                                l.layui.form.on("submit(LAY-user-back-submit)",
                                    function(t) {
                                        console.log(JSON.stringify(t.field))

                                        t.field.id = parseInt(t.field.id)
                                        t.field.roleid = parseInt(t.field.roleid)
                                        t.field.status= parseInt(t.field.status)
                                        t.field.mobile= parseInt(t.field.mobile)
                                        $.ajax({
                                            async: true,
                                            url: '../admin/edit',
                                            type: 'Post',
                                            data :JSON.stringify(t.field),
                                            // dataType: 'json',
                                            contentType: 'application/json',
                                            timeout: 30000,
                                            success: successCallback,
                                            error: errorCallback,
                                            complete: completeCallback,
                                            statusCode: {
                                                404: 404,
                                                500: 500
                                            }
                                        })

                                        // layui.table.reload('LAY-user-back-manage');

                                        layer.close(e)
                                    }),

                                    n.click()
                            },
                            success: function(e, t) {



                                }

                        })
                    }
                }),

            e("useradmin", {})
    });
