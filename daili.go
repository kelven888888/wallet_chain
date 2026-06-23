package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func mains() {
	// 创建一个HTTP服务器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 目标服务器的URL（这里以Google为例，实际使用时需要替换为有效的URL）
		targetURL, err := url.Parse("http://api.tushare.pro/")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("请求转发")
		fmt.Println(targetURL.Host)
		now := time.Now()

		timenow := now.Format("2006-01-02 15:04:05")
		fmt.Println(timenow)

		// 创建一个反向代理
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// 修改请求头（可选）
		// 例如，如果需要传递原始请求的Host头，可以这样做：
		// r.Host = targetURL.Host
		// 但是在这个例子中，我们使用默认的代理行为，所以不需要修改请求头

		// 将请求转发到目标服务器，并将响应返回给客户端
		proxy.ServeHTTP(w, r)
	})

	// 启动服务器
	fmt.Println("Starting proxy server on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
