
    layui.config({
    base: '/resource/layuiadmin/' //静态资源所在路径
}).extend({
    index: 'lib/index' //主入口模块
}).use(['index', 'user'], function(){
    var $ = layui.$
    ,setter = layui.setter
    ,admin = layui.admin
    ,form = layui.form
    ,router = layui.router()
    ,search = router.search;
        initCode();
    form.render();

    //提交
    form.on('submit(LAY-user-login-submit)', function(obj){

    //请求登入接口
    /*   admin.req({
         url:  '../public/loginsubmit' //实际使用请改成服务端真实接口
         ,data: obj.field

         ,method: 'post'
         ,done: function(res){

           //请求成功后，写入 access_token
           layui.data(setter.tableName, {
             key: setter.request.tokenName
             ,value: res.data.access_token
           });
           console.log(res)

           //登入成功的提示与跳转
           layer.msg('登入成功', {
             offset: '15px'
             ,icon: 1
             ,time: 1000
           }, function(){
             location.href = '../admin/index'; //后台主页
           });
         }
       });*/
    $.ajax({
    type: "POST",
    url:  '../public/loginsubmit' ,//实际使用请改成服务端真实接口
    data: obj.field,
    success: function (res) {
    console.log(res);
    if (res.code == 200) {
    layer.msg(res.msg, { icon: 1, time: 500 }, function () {
    location.href = '../admin/main/index'; //后台主页
})
} else {
    layer.msg(res.msg, { icon: 2, time: 2000 }, function () {

    initCode();
})
}

    layer.closeAll('loading');
}
});

});


    //实际使用时记得删除该代码
    /*   layer.msg('为了方便演示，用户名密码可随意输入', {
         offset: '15px'
         ,icon: 1
       });*/

});
function initCode(){
        layui.admin.req({
            async: true,
            url: '../public/captcha',
            type: 'Get',
            // dataType: 'json',
            contentType: 'application/json',
            timeout: 30000,
            done:function(res){

                layui.$("#LAY-user-get-vercode").attr("src",res.data.PicPath)
                layui.$("#CaptchaId").val(res.data.CaptchaId)
            }
        })
    }
