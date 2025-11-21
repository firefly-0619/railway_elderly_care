package factories

import (
	"context"
	"elderly-care-backend/config"
	. "elderly-care-backend/global"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
)

const (
	MUSIC_SOURCE_BUCKET   = "music-source"
	ACCOUNT_AVATAR_BUCKET = "account-avatar"
	FILE_BUCKET           = "file"
)

type OssType int

const (
	ALIYUN OssType = iota
	MINIO
)

type OssFactory struct {
	ossMap map[OssType]OssClient
}

var OssClientFactory *OssFactory

func InitOssFactory() {
	ossFactory := &OssFactory{
		ossMap: make(map[OssType]OssClient),
	}
	ossFactory.initMinioClient()
	OssClientFactory = ossFactory
}

func (this *OssFactory) initMinioClient() {

	minioConfig := config.Config.Oss.Minio
	enable := minioConfig.Enable
	if !enable {
		Logger.Warn("minio not enable")
		return
	}
	endpoint := minioConfig.Endpoint
	accessKeyID := minioConfig.AccessKeyId
	secretAccessKey := minioConfig.SecretAccessKey
	useSSL := minioConfig.UseSsl

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Fprintf(log.Writer(), "Cannot create minio client: %v\n", err)
	}
	Logger.Info("minio client init success")
	client := (*MiniOssClient)(minioClient)
	this.ossMap[MINIO] = client
}

func (this *OssFactory) GetOssClient(ossType OssType) OssClient {
	switch ossType {
	case MINIO:
		return this.ossMap[MINIO]
	default:
		return nil
	}
}

type OssClient interface {
	Upload(bucketName string, objectName string, reader io.Reader, size int64) (string, error)
	Download(bucketName string, objectName string) ([]byte, error)
	GetServiceUrl() string
	GetRawClient() interface{}
}

type MiniOssClient minio.Client

func (client *MiniOssClient) Upload(bucketName string, objectName string, reader io.Reader, size int64) (string, error) {
	ctx := context.Background()
	contentType := "application/octet-stream"
	minioClient := (*minio.Client)(client)
	_, err := minioClient.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	endPoint := minioClient.EndpointURL()
	// 构造公开访问 URL
	return getPolicyUrl(endPoint.Scheme, endPoint.Host, bucketName, objectName), err

}

func (client *MiniOssClient) Download(bucketName string, objectName string) ([]byte, error) {
	ctx := context.Background()
	minioClient := (*minio.Client)(client)
	object, err := minioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (client *MiniOssClient) GetServiceUrl() string {
	minioClient := (*minio.Client)(client)
	endPoint := minioClient.EndpointURL()
	return fmt.Sprintf("%s://%s", endPoint.Scheme, endPoint.Host)
}

// 获取原生客户端
func (client *MiniOssClient) GetRawClient() interface{} {
	return (*minio.Client)(client)
}

func getPolicyUrl(scheme string, host string, bucketName string, objectName string) string {
	return fmt.Sprintf("%s://%s/%s/%s", scheme, host, bucketName, objectName)
}
