package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
	"wallet_chain.com/global"
)

type Parms struct {
	Keys  string
	Value string
}

func Response(ctx *gin.Context, httpStatus int, code int, msg string, data any, args ...Parms) {

	language, _ := ctx.Get("Language")
	if language == nil {
		language = global.SHOP_CONFIG.System.Language
	}
	msg = Languageresponse(msg, language.(string))
	for _, v := range args {
		msg = strings.Replace(msg, v.Keys, v.Value, 1)
	}

	//arg := &args[0]
	//fmt.Println(fmt.Sprintf("%+v", arg))
	//
	//val := reflect.ValueOf(arg)
	//
	//// 检查反射值的类型是否为map
	//fmt.Println(val.Kind(), val, arg)
	//if val.Kind() != reflect.Map {
	//
	//	fmt.Println("Error: Not a map type.")
	//
	//}
	//
	//fmt.Println("--- Iterating Map ---")
	//
	//// 获取map的所有键
	//
	//keys := val.MapKeys()
	//
	//// 遍历键
	//
	//for _, key := range keys {
	//
	//	// 通过键获取对应的值
	//
	//	value := val.MapIndex(key)
	//
	//	// 将reflect.Value转换回原始接口类型，以便打印或进一步处理
	//
	//	fmt.Printf("Key: %v, Value: %v\n", key.Interface(), value.Interface())
	//
	//}
	ctx.JSON(httpStatus, gin.H{
		"code":     code,
		"msg":      msg,
		"data":     data,
		"language": language,
		"website":  global.SHOP_CONFIG.System.WebApiURL,
	})
}

func Success(ctx *gin.Context, msg string, data any, args ...Parms) {

	Response(ctx, http.StatusOK, 200, msg, data, args...)
}

func Fail(ctx *gin.Context, msg string, data any, args ...Parms) {
	fmt.Println(args)
	Response(ctx, http.StatusOK, 400, msg, data, args...)
}

// 获取全部请求参数
func DataMapByRequest(ctx *gin.Context) (dataMap map[string]any, err error) {
	//必须先解析Form
	err = ctx.Request.ParseForm()
	dataMap = make(map[string]any)
	//说明:须post方法,加: 'Content-Type': 'application/x-www-form-urlencoded'
	for key, _ := range ctx.Request.PostForm {
		dataMap[key] = ctx.PostForm(key)
	}
	// 获取Url上的请求参数
	for key, _ := range ctx.Request.URL.Query() {
		dataMap[key] = ctx.Query(key)
	}
	return
}

// 生成指定长度的随机字符
func RandomString(n int) string {
	var letters = []byte("qwertyuioplkjhgfdsazxcvbnmMNBVCXZASDFGHJKLPOIUYTREWQ")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// 获取外网IP
func ExternalIp() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

// 格式化 IP
func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil
	}
	return ip
}

// 过滤指定数组中的key
func ParamsFilter(isFilterStr string, params map[string]any) map[string]any {
	var data = make(map[string]any)
	for key, value := range params {
		if value != "" {
			find := strings.Contains(isFilterStr, key)
			if !find {
				data[key] = value
			}
		}
	}
	return data
}

func UUID() string {
	return uuid.Must(uuid.NewV4()).String()
}

// AnyToMap interface 转 map
func AnyToMap(v any) (map[string]any, error) {
	dataJson, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var MapData map[string]any
	err = json.Unmarshal(dataJson, &MapData)
	if err != nil {
		return nil, err
	}
	return MapData, nil
}
