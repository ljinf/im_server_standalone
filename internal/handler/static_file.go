package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "github.com/ljinf/im_server_standalone/api/v1"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type FileHandler interface {
	UploadImage(ctx *gin.Context)

	//Download(ctx *gin.Context)
}

type fileHandler struct {
	*Handler
}

func (h *fileHandler) UploadImage(ctx *gin.Context) {

	// 获取上传的文件（表单字段名是 "file"）
	file, err := ctx.FormFile("file")
	if err != nil {
		h.logger.Error("FormFileErr" + err.Error())
		v1.HandleError(ctx, http.StatusOK, v1.ErrFileUploadFiled, nil)
		return
	}

	// 检查文件类型
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	if !allowedExts[strings.ToLower(ext)] {
		h.logger.Error("NotAllowedExt", zap.Any("ext", ext))
		v1.HandleError(ctx, http.StatusOK, v1.ErrNotAllowedFileExt, nil)
		return
	}

	//图片长宽  width x height:300x400
	size, exist := ctx.GetPostForm("size")
	if !exist {
		h.logger.Error("GetFileWidthAndHeightFiled")
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)

		return
	}

	now := time.Now()
	subDir := now.Format("20060102")

	// 指定存储路径
	uploadDir := h.conf.GetString("assets.dir")
	fileName := fmt.Sprintf("%v_%v", strings.ToUpper(size), generateRandomFilename(file.Filename))
	dst := filepath.Join(uploadDir, subDir, fileName)

	// 保存文件
	if err = ctx.SaveUploadedFile(file, dst); err != nil {
		h.logger.Error("SaveUploadedFileErr" + err.Error())
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrSaveFileFiled, nil)
		return
	}

	v1.HandleSuccess(ctx, fmt.Sprintf("/%v/%v", subDir, fileName))
}

func (h *fileHandler) MultiUpload(ctx *gin.Context) {
	// 获取multipart表单
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "获取表单失败",
		})
		return
	}

	// 获取所有文件
	files := form.File["files"] // "files"对应前端表单的name属性

	// 遍历保存每个文件
	for _, file := range files {
		filename := filepath.Base(file.Filename)
		// 保存文件到uploads目录
		if err := ctx.SaveUploadedFile(file, "uploads/"+filename); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "文件保存失败",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "所有文件上传成功",
		"count":   len(files),
	})
}

func (h *fileHandler) Download(ctx *gin.Context) {
	v1.HandleSuccess(ctx, nil)
}

func generateRandomFilename(original string) string {
	ext := filepath.Ext(original)
	buf := make([]byte, 8)
	rand.Read(buf)
	return hex.EncodeToString(buf) + ext
}

func NewFileHandler(h *Handler) FileHandler {
	return &fileHandler{
		Handler: h,
	}
}
