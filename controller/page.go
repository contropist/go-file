package controller

import (
	"github.com/gin-gonic/gin"
	"go-file/common"
	"go-file/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func GetIndexPage(c *gin.Context) {
	query := c.Query("query")

	files, _ := model.Query(query)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"message": "",
		"files":   files,
	})
}

func GetExplorerPage(c *gin.Context) {
	path := c.DefaultQuery("path", "/")
	path, _ = url.PathUnescape(path)

	rootPath := filepath.Join(common.LocalFileRoot, path)
	if !strings.HasPrefix(rootPath, common.LocalFileRoot) {
		// We may being attacked!
		rootPath = common.LocalFileRoot
	}
	root, err := os.Stat(rootPath)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"message": err.Error(),
		})
	}
	if root.IsDir() {
		var localFiles []model.LocalFile
		var tempFiles []model.LocalFile
		files, err := ioutil.ReadDir(rootPath)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"message": err.Error(),
			})
		}
		if path != "/" {
			parts := strings.Split(path, "/")
			if len(parts) > 0 {
				parts = parts[:len(parts)-1]
			}
			parentPath := strings.Join(parts, "/")
			parentFile := model.LocalFile{
				Name:         "..",
				Link:         "explorer?path=" + url.QueryEscape(parentPath),
				Size:         "",
				IsFolder:     true,
				ModifiedTime: "",
			}
			localFiles = append(localFiles, parentFile)
			path = strings.Trim(path, "/") + "/"
		} else {
			path = ""
		}
		for _, f := range files {
			file := model.LocalFile{
				Name:         f.Name(),
				Link:         "explorer?path=" + url.QueryEscape(path+f.Name()),
				Size:         common.Bytes2Size(f.Size()),
				IsFolder:     f.Mode().IsDir(),
				ModifiedTime: f.ModTime().String()[:19],
			}
			if file.IsFolder {
				localFiles = append(localFiles, file)
			} else {
				tempFiles = append(tempFiles, file)
			}

		}
		localFiles = append(localFiles, tempFiles...)

		c.HTML(http.StatusOK, "explorer.html", gin.H{
			"message": "",
			"files":   localFiles,
		})
	} else {
		c.File(filepath.Join(common.LocalFileRoot, path))
	}
}

func GetManagePage(c *gin.Context) {
	c.HTML(http.StatusOK, "manage.html", gin.H{
		"message": "",
	})
}
