package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"ssh-client/constkv"
	"ssh-client/ssh_tool"
)

var (
	sshClient = new(ssh_tool.SshClient)
	buffer    = make([]byte, 10240)
)

func handlerCmd(writer http.ResponseWriter, request *http.Request) {
	var resultMap = make(map[string]interface{})
	//处理错误的panic
	defer func() {
		err := recover()
		if err != nil {
			resultMap[constkv.MSG_CODE_KEY] = http.StatusOK
			resultMap[constkv.MSG_CONTENT_KEY] = err.(error).Error()
			bytes, _ := json.Marshal(&resultMap)
			writer.Write(bytes)
		}
	}()

	var bufferSlice = make([]byte, 0)

	reader := request.Body
	bufferReader := bufio.NewReader(reader)

	for {
		length, err := bufferReader.Read(buffer)
		bufferSlice = append(bufferSlice, buffer[:length]...)
		if err == io.EOF {
			break
		}
	}

	log.Println(string(bufferSlice))
	var paramMap = make(map[string]interface{})
	_ = json.Unmarshal(bufferSlice, &paramMap)

	log.Println(paramMap["cmd"])
	cmd := paramMap["cmd"]
	log.Println(cmd.(string) + "--------------------")

	if sshClient.Client != nil {
		runCmd := sshClient.RunCmd(cmd.(string))
		resultMap[constkv.MSG_CODE_KEY] = http.StatusOK
		resultMap[constkv.MSG_CONTENT_KEY] = runCmd
	}

	bytes, _ := json.Marshal(&resultMap)
	_, _ = writer.Write(bytes)
}

func handleUploadFile(writer http.ResponseWriter, request *http.Request) {
	var resultMap = make(map[string]interface{})
	//处理错误的panic
	defer func() {
		err := recover()
		if err != nil {
			resultMap[constkv.MSG_CODE_KEY] = http.StatusOK
			resultMap[constkv.MSG_CONTENT_KEY] = err.(error).Error()
			bytes, _ := json.Marshal(&resultMap)
			writer.Write(bytes)
		}

	}()
	var bufferSlice = make([]byte, 0)
	reader := request.Body
	bufferReader := bufio.NewReader(reader)

	for {
		length, err := bufferReader.Read(buffer)
		bufferSlice = append(bufferSlice, buffer[:length]...)
		if err == io.EOF {
			break
		}
	}

	log.Println(string(bufferSlice))
	var paramMap = make(map[string]interface{})
	_ = json.Unmarshal(bufferSlice, &paramMap)

	localPath := paramMap["localPath"]
	remoteDir := paramMap["remoteDir"]
	remoteFileName := paramMap["remoteFileName"]

	log.Println(localPath)
	log.Println(remoteDir)

	sshClient.UploadFile(localPath.(string), remoteDir.(string), remoteFileName.(string))
}

func handleDownloadFile(writer http.ResponseWriter, request *http.Request) {
	var resultMap = make(map[string]interface{})
	//处理错误的panic
	defer func() {
		err := recover()
		if err != nil {
			resultMap[constkv.MSG_CODE_KEY] = http.StatusOK
			resultMap[constkv.MSG_CONTENT_KEY] = err.(error).Error()
			bytes, _ := json.Marshal(&resultMap)
			writer.Write(bytes)
		}

	}()
	var bufferSlice = make([]byte, 0)
	reader := request.Body
	bufferReader := bufio.NewReader(reader)

	for {
		length, err := bufferReader.Read(buffer)
		bufferSlice = append(bufferSlice, buffer[:length]...)
		if err == io.EOF {
			break
		}
	}

	log.Println(string(bufferSlice))
	var paramMap = make(map[string]interface{})
	_ = json.Unmarshal(bufferSlice, &paramMap)

	remotePath := paramMap["remotePath"]
	localDir := paramMap["localDir"]
	localFileName := paramMap["localFileName"]

	log.Println(remotePath)
	log.Println(localDir)
	log.Println(localFileName)

	sshClient.DownloadFile(remotePath.(string), localDir.(string), localFileName.(string))
}

func main() {
	//初始化ssh客户端
	sshClient.NewSshClient()
	//1.注册一个给定模式的处理器函数到DefaultServeMux
	http.HandleFunc("/", handlerCmd)
	http.HandleFunc("/uploadFile", handleUploadFile)
	http.HandleFunc("/downloadFile", handleDownloadFile)

	//2.设置监听的TCP地址并启动服务
	//参数1：TCP地址(IP+Port)
	//参数2：当设置为nil时表示使用DefaultServeMux
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
