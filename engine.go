package http

import "github.com/labstack/echo"

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/asyou-me/lib.v1/utils"

	"git.asyou.me/old/potential-mentor/pkg/handlers"
	"github.com/asyou-me/http/types"
)

// Engine http engine
type Engine struct {
	*echo.Echo
}

// LoadRestfull 加载所有 restfull 文件
func (*Engine) LoadRestfull(path string) error {
	// 完善路径
	utils.CompletePath(&path)
	// 列出路径下所有配置文件
	files, err := utils.ListDir(path, ".json")
	if err != nil {
		return err
	}
	for _, file := range files {
		// 读取文件
		fi, err := os.Open(file)
		if err != nil {
			return errors.New("传入的配置文件路径: " + file + " 不存在")
		}
		defer fi.Close()
		fd, err := ioutil.ReadAll(fi)
		if err != nil {
			return err
		}

		// 将文件内容填入 RestFullConf
		v := &types.RestFullConf{}
		if err = json.Unmarshal(fd, v); err != nil {
			if err != nil {
				return err
			}
		}
		handlers.RestFulls[v.Name] = v
	}
	return nil
}
