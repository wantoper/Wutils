package main

import (
	"WUtils/WHttp"
	"fmt"
	"io"
)

func main() {
	server := WHttp.NewServer(":8080")

	server.GET("/", func(w WHttp.ResponseWriter, r *WHttp.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("<h1>欢迎使用WHttp框架</h1>"))
	})

	server.GET("/hello", func(w WHttp.ResponseWriter, r *WHttp.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("Hello, World!"))
	})

	server.POST("/api/data", func(w WHttp.ResponseWriter, r *WHttp.Request) {
		data, _ := io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(fmt.Sprintf(`{"message":"收到数据","data":"%s"}`, string(data))))
	})

	// 启动服务器
	fmt.Println("服务器启动在 :8080")
	server.ListenAndServe()
}
