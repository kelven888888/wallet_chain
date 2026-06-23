package service

import (
	"mime/multipart"
	"wallet_chain.com/utils/uploads"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Upload
//@description: 创建文件上传记录
//@param: file model.ExaFileUploadAndDownload
//@return: error

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UploadFile
//@description: 根据配置文件判断是文件上传到本地或者七牛云
//@param: header *multipart.FileHeader, noSave string
//@return: file model.ExaFileUploadAndDownload, err error

type FileUploadAndDownloadService struct{}

func (e *FileUploadAndDownloadService) UploadFile(header *multipart.FileHeader, noSave string, types string) (filePath string, err error) {
	oss := upload.NewOss()
	filePath, _, uploadErr := oss.UploadFile(header, types)
	if uploadErr != nil {
		panic(uploadErr)
	}

	return filePath, nil
}
