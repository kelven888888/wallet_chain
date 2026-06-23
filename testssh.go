package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

func RestartTrade() {
	pid := Exce("lsof -t -i:2666 -sTCP:LISTEN")
	if pid == "" {
		fmt.Println("no process listen port 2666")
		go Exce("cd /www/python/market/&&source py3market/bin/activate &&nohup python manage.py runserver>>/www/python/market/nohup.out&\n")
		time.Sleep(5 * time.Second)
		pid := Exce("lsof -t -i:2666 -sTCP:LISTEN")
		if pid != "" {
			fmt.Println("启动成功.进程id", pid)
		} else {
			fmt.Println("kill port 2666 fail")
		}
		return
	}
	pids := Exce(fmt.Sprintf(" kill -9 %s 2>/dev/null", pid))
	if pids == "" {
		go Exce("cd /www/python/market/&&source py3market/bin/activate &&nohup python manage.py runserver>>/www/python/market/nohup.out&\n")
		time.Sleep(5 * time.Second)
		pid := Exce("lsof -t -i:2666 -sTCP:LISTEN")
		if pid != "" {
			fmt.Println("启动成功.进程id", pid)
		} else {
			fmt.Println("kill port 2666 fail")
		}
	} else {
		fmt.Println("kill port 2666 fail")
	}
}
func Exce(comand string) string {
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("UdkH@fuVZMfs"), // 或者使用 ssh.PublicKeys(你的公钥)
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 注意：生产环境下不要使用 InsecureIgnoreHostKey()
	}

	// 建立 SSH 连接
	conn, err := ssh.Dial("tcp", "18.163.2.192:22", config)
	if err != nil {
		log.Fatalf("无法建立连接: %v", err)
	}
	defer conn.Close()

	// 创建 SSH 会话
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("无法创建会话: %v", err)
	}
	defer session.Close()

	// 设置会话选项，例如保持远程命令的输出打开，直到它完成。

	// 如果需要的话，可以启动一个终端会话。例如，对于需要交互式输入的命令。
	// session.RequestPty("xterm", 80, 40, modes)

	// 执行命令并获取输出
	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	err = session.Run(comand)

	//err = session.Run("cd /www/python/market/ &&source py3market/bin/activate &&nohup python manage.py runserver>>/www/python/market/nohup.out&") // 这里替换为你想执行的命令
	if err != nil && stderrBuf.String() != "" {
		log.Fatalf("远程命令执行失败: %v\n错误输出: %s,命令%s", err, stderrBuf.String(), comand)
	}

	// 输出命令的标准输出结果
	fmt.Printf("命令输出: %s\n", stdoutBuf.String())
	return stdoutBuf.String()
}
