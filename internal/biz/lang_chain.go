package biz

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/d2jvkpn/x-ai/pkg/lang_chain"
	// "github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

type LangChainAgent struct {
	*lang_chain.LangChain
}

type ChainQuery struct {
	FaissIndex string `json:"faissIndex"`
	Query      string `json:"query"`
}

func NewLangChainAgent(key, path string) (lca *LangChainAgent, err error) {
	lca = new(LangChainAgent)

	if lca.LangChain, err = lang_chain.NewLangChain(key, path); err != nil {
		return nil, err
	}

	return lca, nil
}

func (lca *LangChainAgent) Filename(name string) (ext string, err error) {
	switch {
	case strings.HasSuffix(name, ".pdf"):
		return "pdf", nil
	case strings.HasSuffix(name, ".txt"):
		return "txt", nil
	case strings.HasSuffix(name, ".doc"):
		return "doc", nil
	case strings.HasSuffix(name, ".docx"):
		return "docx", nil
	default:
		return "", fmt.Errorf("unknow file type(ext)")
	}
}

func (lca *LangChainAgent) PyIndex(prefix string) {
	var err error

	if err = lca.LangChain.PyIndex(context.TODO(), prefix+".tmp.yaml", prefix); err != nil {
		log.Printf("!!! PyIndex: %v\n", err)
		return
	}

	if err = os.Rename(prefix+".tmp.yaml", prefix+".yaml"); err != nil {
		log.Printf("!!! PyIndex: %v\n", err)
		return
	}
}

func (lca *LangChainAgent) HandleIndex(ctx *gin.Context) (indexName string, err error) {
	var (
		idx         int
		msg         string
		prefix, ext string
		list        []string
		fh          *multipart.FileHeader
		form        *multipart.Form
		files       []*multipart.FileHeader
		index       *lang_chain.FaissIndex
	)

	msg = "ok"
	if form, err = ctx.MultipartForm(); err != nil {
		return "", err
	}

	files = form.File["sources"]
	if len(files) == 0 {
		msg = "has no sources"
		ctx.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": msg})
		return "", fmt.Errorf(msg)
	}

	sources := make([]lang_chain.Source, 0, len(files))
	if index, err = lang_chain.NewFaissIndex(sources); err != nil {
		return "", err
	}

	indexName = "faiss-index_" + index.UtcTime.Format(time.DateOnly) + "/" + index.Uuid()
	prefix = filepath.Join(
		lca.GetPath(),
		"faiss-index_"+index.UtcTime.Format(time.DateOnly),
		index.Uuid(),
	)

	list = make([]string, 0, len(files))
	for idx, fh = range files {
		if ext, err = lca.Filename(fh.Filename); err != nil {
			msg = fmt.Sprintf("invalid file: %s", fh.Filename)
			ctx.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": msg})
			return "", fmt.Errorf(msg)
		}
		fp := fmt.Sprintf("%s_%03d.%s", prefix, idx, ext)

		index.Sources = append(index.Sources, lang_chain.Source{
			Name:   fh.Filename,
			Type:   ext,
			Source: filepath.Base(fp),
		})

		list = append(list, fp)
	}

	if err = os.MkdirAll(filepath.Dir(prefix), 0755); err != nil {
		msg = "internal error"
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": msg})
		return "", err
	}

	for idx, fh = range files {
		if err = ctx.SaveUploadedFile(fh, list[idx]); err != nil {
			msg = "internal error"
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 2, "msg": msg})
			return "", err
		}
	}

	if err = index.SaveYaml(prefix + ".tmp.yaml"); err != nil {
		msg = "internal error"
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 3, "msg": msg})
		return "", err
	}

	// TODO: ?? async
	/*
		if err = lca.LangChain.PyIndex(ctx, prefix+".tmp.yaml", prefix); err != nil {
			msg = "internal error"
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 4, "msg": msg})
			return "", err
		}
	*/
	go lca.PyIndex(prefix)

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": msg, "data": gin.H{"faissIndex": indexName}})
	return indexName, nil
}

func (lca *LangChainAgent) HandleQuery(ctx *gin.Context) (err error) {
	var (
		prefix   string
		msg, ans string
		query    ChainQuery
	)

	msg = "ok"

	if err = ctx.BindJSON(&query); err != nil {
		msg = "unmarshal failed"
		ctx.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": msg})
		return err
	}
	if query.FaissIndex == "" || query.Query == "" {
		msg = "faissIndex or query is empty"
		ctx.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": msg})
		return err
	}

	prefix = filepath.Join(lca.GetPath(), query.FaissIndex)
	if _, err = os.Stat(prefix + ".yaml"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			msg = "faissIndex not exists"
			ctx.JSON(http.StatusNotAcceptable, gin.H{"code": -3, msg: msg})
			return
		}

		msg = "internal error"
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 1, msg: msg})
		return
	}

	if ans, err = lca.PyQuery(ctx, prefix, query.Query); err != nil {
		msg = "internal error"
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 2, msg: msg})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, msg: msg, "data": gin.H{"ans": ans}})

	return nil
}
