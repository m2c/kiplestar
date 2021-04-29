package utils

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	slog "github.com/mark-jiang-gp/kiplestar/commons/log"
	"io"
	"io/ioutil"
	"time"
)

type OSSClient interface {
	UploadAndSignUrl(fileReader io.Reader, objectName string, expiredInSec int64) (string, error)
	DeleteByObjectName(objectName string)
	UploadByReader(fileReader io.Reader, fileName string) (err error)
	DownloadFile(fileName string) (data []byte, err error)
	IsFileExist(fileName string) (isExist bool, err error)
	GetFileURL(fileName string, expireTime time.Duration) (url string, err error)
}

type ossClientImp struct {
	ossBucket       string
	accessKeyID     string
	accessKeySecret string
	ossEndPoint     string
}

func (slf *ossClientImp) GetFileURL(fileName string, expireTime time.Duration) (url string, err error) {
	client, err := oss.New(slf.ossEndPoint, slf.accessKeyID, slf.accessKeySecret)
	if err != nil {
		slog.Errorf("ossClientImp IsFileExist Error:%s", err)
		return "", err
	}
	bucket, err := client.Bucket(slf.ossBucket)
	if err != nil {
		slog.Errorf("ossClientImp IsFileExist  Error:%s", err)
		return "", err
	}
	//oss.Process("image/format,png")
	url, err = bucket.SignURL(fileName, oss.HTTPGet, int64(expireTime))
	if err != nil {
		return "", err
	}
	return

}

func OSSClientInstance(ossBucket, accessKeyID, accessKeySecret, ossEndPoint string) OSSClient {
	return &ossClientImp{
		ossBucket:       ossBucket,
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		ossEndPoint:     ossEndPoint,
	}
}

func (slf *ossClientImp) IsFileExist(fileName string) (isExist bool, err error) {
	client, err := oss.New(slf.ossEndPoint, slf.accessKeyID, slf.accessKeySecret)
	if err != nil {
		slog.Errorf("ossClientImp IsFileExist Error:%s", err)
		return false, err
	}
	bucket, err := client.Bucket(slf.ossBucket)
	if err != nil {
		slog.Errorf("ossClientImp IsFileExist  Error:%s", err)
		return false, err
	}
	return bucket.IsObjectExist(fileName)
}

func (slf *ossClientImp) UploadAndSignUrl(fileReader io.Reader, objectName string, expiredInSec int64) (string, error) {
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
	//oss.Process("image/format,png")
	signedURL, err := bucket.SignURL(objectName, oss.HTTPGet, expiredInSec)
	if err != nil {
		bucket.DeleteObject(objectName)
		return "", err
	}
	return signedURL, nil
}

func (slf *ossClientImp) DeleteByObjectName(objectName string) {
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

func (slf *ossClientImp) UploadByReader(fileReader io.Reader, fileName string) (err error) {

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

	err = bucket.PutObject(fileName, fileReader)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	return
}

func (slf *ossClientImp) DownloadFile(file_name string) (data []byte, err error) {

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

	body, err := bucket.GetObject(file_name)
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
