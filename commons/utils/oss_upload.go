package utils

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	slog "github.com/m2c/kiplestar/commons/log"
	"io"
	"io/ioutil"
)

const ossEndPoint =  "oss-ap-southeast-3.aliyuncs.com"
func UploadByReader(fileReader io.Reader, fileName string,ossKeyID,ossKeySecret,ossBucket string) (err error) {

	client, err := oss.New(ossEndPoint, ossKeyID, ossKeySecret)
	if err != nil {
		slog.Errorf("Error:%s", err)
		return
	}

	bucket, err := client.Bucket(ossBucket)
	if err != nil {
		slog.Errorf("Error:%s", err)
		return
	} else {
		fmt.Println("bukect ok")
	}
	objectName := "bank_file/" + fileName

	err = bucket.PutObject(objectName, fileReader)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	return
}

func DownloadFile(file_name string,ossKeyID,ossKeySecret,ossBucket string) (data []byte, err error) {

	client, err := oss.New(ossEndPoint, ossKeyID, ossKeySecret)
	if err != nil {
		slog.Errorf("Error:%s", err)
		return
	}

	bucket, err := client.Bucket(ossBucket)
	if err != nil {
		slog.Errorf("Error:%s", err)
		return
	} else {
		fmt.Println("bukect ok")
	}

	// 下载文件到流。
	objectName := "bank_file/" + file_name

	body, err := bucket.GetObject(objectName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 数据读取完成后，获取的流必须关闭，否则会造成连接泄漏，导致请求无连接可用，程序无法正常工作。
	defer body.Close()

	data, err = ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	return
}
