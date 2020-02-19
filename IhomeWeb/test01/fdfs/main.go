package main

import (
	"fmt"
	"github.com/tedcy/fdfs_client"
	"path"
)

func UploadByFileName(filename string) (string, error) {
	client, err := fdfs_client.NewClientWithConfig("/Users/mds/Documents/go/sss/IhomeWeb/conf/fdfs.conf")
	if err != nil {
		fmt.Println("FDFS 初始化异常 FileName ：", err)
		return "", err
	}

	fdfs_file_path, err := client.UploadByFilename(filename)
	if err != nil {
		fmt.Println("FDFS 上传失败: ", err)
		return "", err
	}
	fmt.Println(fdfs_file_path)
	return fdfs_file_path, nil
}

func UploadByBuffer(fileBuffer []byte, fileExtName string) (string, error) {
	client, err := fdfs_client.NewClientWithConfig("/Users/mds/Documents/go/sss/IhomeWeb/conf/fdfs.conf")
	if err != nil {
		fmt.Println("FDFS 初始化异常 Buffer：", err)
		return "", nil
	}

	fdfs_file_path, err := client.UploadByBuffer(fileBuffer, fileExtName)
	if err != nil {
		fmt.Println("FDFS 上传失败：", err)
		return "", nil
	}
	fmt.Println(fdfs_file_path)
	return fdfs_file_path, nil
}

func main() {
	filePath := "/Users/mds/Downloads/chrome/fm.jpg"
	//UploadByFileName(filePath)

	d := path.Ext(filePath)
	fmt.Println(d)

}
