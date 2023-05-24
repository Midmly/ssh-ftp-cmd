package utils

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"path"
	"strconv"
)

type SshConfig struct {
	Label      string      `json:"label" form:"label" query:"label"`
	Username   string      `json:"username" form:"username" query:"username"`
	Password   string      `json:"password" form:"password" query:"password"`
	Address    string      `json:"address" form:"address" query:"address"`
	Port       int         `json:"port" form:"port" query:"port"`
	PrivateKey string      `json:"private_key" form:"private_key" query:"private_key"`
	Width      int         `json:"width"`
	Height     int         `json:"height"`
	Client     *ssh.Client `json:"-" form:"-" query:"-"`
}

func (c *SshConfig) NewSshClient() error {
	var auth ssh.AuthMethod
	signer, err := ssh.ParsePrivateKey([]byte(c.PrivateKey))
	if err != nil {
		log.Printf("无法解析私钥：%v", err)
		auth = ssh.Password(c.Password)
	} else {
		auth = ssh.PublicKeys(signer)
	}
	// 连接SSH服务器
	c.Client, err = ssh.Dial("tcp", c.Address+":"+strconv.Itoa(c.Port), &ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	return err
}

func (c *SshConfig) RunCmd(cmd string) string {
	session, err := c.Client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	runResult, err := session.CombinedOutput(cmd)
	if err != nil {
		return err.Error()
	}
	return string(runResult)
}

func (c *SshConfig) UploadFile(localPath string, remoteDir string, remoteFileName string) {
	ftpClient, err := sftp.NewClient(c.Client)
	if err != nil {
		fmt.Println("创建ftp客户端失败", err)
		panic(err)
	}

	defer ftpClient.Close()

	fmt.Println(localPath, remoteFileName)
	srcFile, err := os.Open(localPath)
	if err != nil {
		fmt.Println("打开文件失败", err)
		panic(err)
	}
	defer srcFile.Close()

	dstFile, e := ftpClient.Create(path.Join(remoteDir, remoteFileName))
	if e != nil {
		fmt.Println("创建文件失败", e)
		panic(e)
	}
	defer dstFile.Close()

	buffer := make([]byte, 1024000)
	for {
		n, err := srcFile.Read(buffer)
		dstFile.Write(buffer[:n])
		//注意，由于文件大小不定，不可直接使用buffer，否则会在文件末尾重复写入，以填充1024的整数倍
		if err != nil {
			if err == io.EOF {
				fmt.Println("已读取到文件末尾")
				break
			} else {
				fmt.Println("读取文件出错", err)
				panic(err)
			}
		}
	}
}

func (c *SshConfig) DownloadFile(remotePath string, localDir string, localFilename string) {
	ftpClient, err := sftp.NewClient(c.Client)
	if err != nil {
		fmt.Println("创建ftp客户端失败", err)
		panic(err)
	}

	defer ftpClient.Close()

	srcFile, err := ftpClient.Open(remotePath)
	if err != nil {
		fmt.Println("文件读取失败", err)
		panic(err)
	}
	defer srcFile.Close()

	dstFile, e := os.Create(path.Join(localDir, localFilename))
	if e != nil {
		fmt.Println("文件创建失败", e)
		panic(e)
	}
	defer dstFile.Close()
	if _, err1 := srcFile.WriteTo(dstFile); err1 != nil {
		fmt.Println("文件写入失败", err1)
		panic(err1)
	}
	fmt.Println("文件下载成功")
}
