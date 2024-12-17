// -------------------------------------------------
// Package kb_apis
// Author: hanzhi
// Date: 2024/12/16
// -------------------------------------------------

package kb_apis

import (
	"gcnote/server/config"
	"gcnote/server/dto"
	"gcnote/server/model"
	"gcnote/server/router/wrench"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
)

// RecycleKBFile
// [rewrite]需要重写 -
// todo: 还有个事情，就是如果在共享表中存在这个，那么就要询问用户，是否确定删除，确定，那么先删除所有共享记录，然后才能彻底删除。
//
// todo: 但是这个就不要在一个接口里耦合了，分成查询共享记录、删除文件、删除共享记录几个部分，然后前端整合在一起显示。
//
// @Summary 删除知识库里的文件（fixme:新写一个吧，写错了，应该是收到回收站的，这个留给回收站了）
// @Description 删除知识库内部的文件(MYSQL + 文件系统已完成 + redis不管了)
// @ID			recycle-kb-file
// @Tags		index
// @Accept		json
// @Produce		json
// @Param		request		body		dto.KBFileUDRequest true "文档请求体"
// @Success 	200			{object} 	dto.BaseResponse	"成功响应，返回success"
// @Failure		400			{object}	dto.BaseResponse	"参数错误(code:40000)"
// @Failure		401			{object}	dto.BaseResponse	"Token错误(code:40101)"
// @Failure		404			{object}	dto.BaseResponse	"该知识库不存在(code:40201)"
// @Failure		500			{object}	dto.BaseResponse	"服务器内部错误(code:50000)"
// @Router 		/index/recycle_file [post]
func RecycleKBFile(ctx *gin.Context) {
	var req dto.KBFileUDRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Fail(dto.ParamsErrCode))
		return
	}
	// 获取用户信息
	claims, exists := ctx.Get("claims")
	if !exists {
		zap.S().Infof("Unable to get the claims")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.UserTokenErrCode))
		return
	}
	currentUser := claims.(jwt.MapClaims)
	currentUserId := currentUser["sub"].(string)
	if currentUserId == "" {
		zap.S().Debugf("currentUserId is empty")
		ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.ParamsErrCode))
		return
	}

	// 回收文件
	src := filepath.Join(config.PathCfg.KnowledgeBasePath, req.IndexId, req.KBFileId)
	dst := filepath.Join(config.PathCfg.RecycleBinPath, req.IndexId, req.KBFileId)
	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		zap.S().Errorf("Failed to create recycle dir, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 复制文件到要求的路径
	err = wrench.CopyDir(src, dst)
	if err != nil {
		zap.S().Errorf("Failed to create recycle dir, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}

	// 删除原来的文件夹
	err = wrench.RemoveContents(src)
	if err != nil {
		zap.S().Errorf("Delete kb_file id :%v err:%v", req.KBFileId, err)
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 然后sql操作
	recycleModel := model.Recycle{
		UserId:        currentUserId,
		SourceIndexId: req.IndexId,
		KBFileId:      req.KBFileId,
		KBFileName:    req.KBFileName,
	}
	kbFile := model.KBFile{
		IndexId:    req.IndexId,
		KBFileId:   req.KBFileId,
		KBFileName: req.KBFileName,
	}

	// 首先删除旧表
	tx := config.DB.Begin()
	tx.Model(&kbFile).Where(
		"kb_file_id = ?",
		kbFile.KBFileId,
	).Delete(&kbFile)
	if tx.Error != nil {
		zap.S().Errorf("Delete kb_file name :%v err:%v", req.KBFileName, tx.Error)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 然后放到回收站里
	tx.Create(&recycleModel)
	if tx.Error != nil {
		zap.S().Errorf("recycle kb_file name :%v err:%v", req.KBFileName, tx.Error)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		zap.S().Errorf("Failed to commit transaction, err: %v", err)
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, dto.Fail(dto.InternalErrCode))
		return
	}
	zap.S().Infof("Delete kb_file_name %v , user id : %v done.", req.KBFileName, currentUserId)
	ctx.JSON(http.StatusOK, dto.Success())
}
