package utils

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	slog "github.com/m2c/kiplestar/commons/log"
	"github.com/m2c/kiplestar/config"
	"io"
	"io/ioutil"
)

<<<<<<< Updated upstream
type OSSClient interface {
	UploadAndSignUrl(fileReader io.Reader, objectName string, expiredInSec int64) (string, error)
	DeleteByObjectName(objectName string)
}

type ossClientImp struct {
	ossBucket       string
	accessKeyID     string
	accessKeySecret string
	ossEndPoint     string
}
=======
const ossEndPoint =  "oss-ap-southeast-3.aliyuncs.com"

func UploadByReader(fileReader io.Reader, fileName string,ossKeyID,ossKeySecret,ossBucket string) (err error) {
>>>>>>> Stashed changes

func OSSClientInstance() OSSClient {
	return  &ossClientImp{
		ossBucket: config.Configs.Oss.OssBucket,
		accessKeyID: config.Configs.Oss.AccessKeyID,
		accessKeySecret: config.Configs.Oss.AccessKeySecret,
		ossEndPoint: config.Configs.Oss.OssEndPoint,
	}
}
func(slf *ossClientImp) UploadAndSignUrl(fileReader io.Reader, objectName string, expiredInSec int64) (string, error){
	// 创建OSSClient实例。
	client, err := oss.New(slf.ossEndPoint, slf.accessKeyID, slf.accessKeySecret)
	if err != nil {
		slog.Errorf("Error:%s", err)
		return "", err
	}
	bucket, err := client.Bucket(slf.ossBucket)
	if err != nil {
		slog.Errorf("Error:%s", err)
		return "", err
	}
	err = bucket.PutObject(objectName, fileReader)
	if err != nil {
		slog.Errorf("Error:%s", err)
		return "", err
	}
	signedURL, err := bucket.SignURL(objectName, oss.HTTPGet, expiredInSec, oss.Process("image/format,png"))
	if err != nil {
		bucket.DeleteObject(objectName)
		return "", err
	}
	return signedURL, nil
}

func(slf *ossClientImp)  DeleteByObjectName(objectName string) {
	client, err := oss.New(slf.ossEndPoint, slf.accessKeyID, slf.accessKeySecret)
	if err != nil {
		slog.Errorf("Error:%s", err)
		return
	}
	bucket, err := client.Bucket(slf.ossBucket)
	if err != nil {
		slog.Errorf("Error:%s", err)
		return
	}
	err = bucket.DeleteObject(objectName)
	if err != nil {
		slog.Errorf("Error:%s", err)
	}
}

func(slf *ossClientImp) UploadByReader(fileReader io.Reader, fileName string) (err error) {

	client, err := oss.New(slf.ossEndPoint, slf.accessKeyID, slf.accessKeySecret)
	if err != nil {
		slog.Errorf("Error:%s", err)
		return
	}

	bucket, err := client.Bucket(slf.ossBucket)
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

func(slf *ossClientImp)  DownloadFile(file_name string) (data []byte, err error) {

	client, err := oss.New(slf.ossEndPoint, slf.accessKeyID, slf.accessKeySecret)
	if err != nil {
		slog.Errorf("Error:%s", err)
		return
	}

	bucket, err := client.Bucket(slf.ossBucket)
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


