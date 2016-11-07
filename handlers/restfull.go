package handlers

import (
	"fmt"

	"github.com/asyou-me/postgres"

	"git.asyou.me/old/potential-mentor/pkg/errors"
	"git.asyou.me/old/potential-mentor/pkg/models"
	"git.asyou.me/old/potential-mentor/pkg/types"
)

// RestFulls 资源定义
var RestFulls = map[string]*types.RestFullConf{}

// POST 创建
func POST(req *types.REQ, resq *types.RESP) error {
	// 获取当前资源的定义文件
	post := RestFulls[req.Type].POST
	// 当前资源的关联数据定义文件
	relationships := *(post.Req.Relationships)
	if req.Relationships != nil {
		for k, v := range *req.Relationships {
			if v.T == postgres.Map { // 单条关联数据
				id := (*v.Map)["id"].V
				m := *relationships[k].Define
				for k := range m {
					(*v.Map)[k] = &postgres.V{}
				}
				err := models.GET(k, "id="+id, *v.Map)
				if err != nil {
					er := errors.New(
						errors.INVALID_PARAME, "relationships."+k+".id",
						"没有找到 id 为 "+id+" 的 "+k,
					)
					er.LogMessage(err.Error())
					return er
				}
				(*req.Attributes)[k] = (*v.Map)["id"]
			} else { // 多条关联数据
				for _, value := range *v.Slice {
					id := (*value)["id"].V
					m := *relationships[k].Define
					for k := range m {
						(*value)[k] = &postgres.V{}
					}
					err := models.GET(k, "id="+id, *value)
					if err != nil {
						er := errors.New(
							errors.INVALID_PARAME, "relationships."+k+".id",
							"没有找到 id 为 "+id+" 的 "+k,
						)
						er.LogMessage(err.Error())
						return er
					}
				}
				(*req.Attributes)[k] = &postgres.V{
					T:        postgres.IntArray,
					IntArray: v.Array(),
				}
			}
		}
	}
	_, err := models.INSERT(req.Type, *req.Attributes)
	if err != nil {
		msg := "无法创建数据，请检查参数"
		if post.Msg != nil {
			msg = *post.Msg
		}
		er := errors.New(errors.NULL, msg)
		er.LogMessage(err.Error())
		return er
	}
	resq.Data = &types.RESPDATA{
		T:   postgres.Map,
		Map: req.Attributes,
	}
	resq.Type = req.Type
	resq.Included = req.Relationships
	return nil
}

// PACTH 更新
func PACTH(req *types.REQ, resq *types.RESP) error {
	relationships := *(RestFulls[req.Type].POST.Req.Relationships)
	if req.Relationships != nil {
		for k, v := range *req.Relationships {
			if v.Map == nil {
				v.Map = &map[string]*postgres.V{}
			}

			if v.T == postgres.Map {
				id := (*v.Map)["id"].V
				m := *relationships[k].Define
				for k, d := range m {
					(*v.Map)[k] = &postgres.V{
						T: d.T,
					}
				}
				// TODO 此处为检查数据是否存在 不需要获取全部的内容
				err := models.GET(k, "id="+id, *v.Map)
				if err != nil {
					er := errors.New(
						errors.INVALID_PARAME, "relationships."+k+".id",
						"没有找到 id 为 "+id+" 的 "+k,
					)
					er.LogMessage(err.Error())
					return er
				}
				(*req.Attributes)[k] = (*v.Map)["id"]
			} else {
				if len(*v.Slice) == 0 {
					continue
				}
				for _, value := range *v.Slice {
					id := (*value)["id"].V
					m := *relationships[k].Define
					for k, d := range m {
						(*v.Map)[k] = &postgres.V{
							T: d.T,
						}
					}
					err := models.GET(k, "id="+id, *value)
					if err != nil {
						er := errors.New(
							errors.INVALID_PARAME, "relationships."+k+".id",
							"没有找到 id 为 "+id+" 的 "+k,
						)
						er.LogMessage(err.Error())
						return er
					}
				}
				(*req.Attributes)[k] = &postgres.V{
					T:        postgres.IntArray,
					IntArray: v.Array(),
				}
			}
		}
	}

	_, err := models.UPDATE(req.Type, "id="+req.Query.ID, *req.Attributes)
	if err != nil {
		return err
	}
	resq.Type = req.Type
	resq.Data = &types.RESPDATA{
		T: postgres.Map,
		Map: &map[string]*postgres.V{
			"id": &postgres.V{
				T: postgres.Int8,
				V: req.Query.ID,
			},
		},
	}
	return nil
}

// LIST 获取列表
// TODO include 查询字段
func LIST(req *types.REQ, resq *types.RESP) error {
	attrs := *(RestFulls[req.Type].LIST.Req.Attributes)
	files := map[string]uint8{}
	files["id"] = postgres.Int8
	for k, v := range attrs {
		files[k] = v.T
	}
	req.Query.Fields = files
	// 长度不超过 req.Query.Limit 时可直接传递 slice 本身
	datas, err := models.List(req.Type, req.Query)
	if err != nil {
		er := errors.New(errors.NULL, err.Error())
		er.LogMessage(err.Error())
		return er
	}
	resq.Type = req.Type
	resq.Data = &types.RESPDATA{
		T:     postgres.Slice,
		Slice: &datas,
	}
	return nil
}

// GET 获取
func GET(req *types.REQ, resq *types.RESP) error {
	get := RestFulls[req.Type].GET.Resq
	relationships := *(RestFulls[req.Type].GET.Resq.Relationships)
	attrs := *(RestFulls[req.Type].GET.Resq.Attributes)
	vMap := map[string]*postgres.V{}
	vMap["id"] = &postgres.V{
		T: postgres.Int8,
	}

	for k, v := range *get.Attributes {
		vMap[k] = &postgres.V{
			T: v.T,
		}
	}

	// 获取参数中需要包含的信息
	Includes := map[string]*postgres.V{}
	if req.Query.Include != nil {
		for _, k := range *req.Query.Include {
			r, ok := relationships[k]
			if ok {
				v := &postgres.V{
					T: r.T,
				}
				vMap[k] = v
				Includes[k] = v
			}
		}
	}

	err := models.GET(req.Type, "id="+req.Query.ID, vMap)
	if err != nil {
		er := errors.New(
			errors.INVALID_PARAME, "id",
			"没有找到 id 为 "+req.Query.ID+" 的 "+req.Type,
		)
		er.LogMessage(err.Error())
		return er
	}

	if len(Includes) > 0 {
		resq.Included = &map[string]*types.Relationship{}
		for k, v := range Includes {
			ins := &types.Relationship{}
			m := *relationships[k].Define
			if v.T == postgres.IntArray {
				ins.T = postgres.Slice
				// TODO 检查数据是否为真
				if v.IntArray == nil {
					continue
				}
				s := make([]*map[string]*postgres.V, len(*v.IntArray))
				ins.Slice = &s
				for index, value := range *v.IntArray {
					id := fmt.Sprint(value)
					item := &map[string]*postgres.V{}
					for k, d := range m {
						(*item)[k] = &postgres.V{
							T: d.T,
						}
					}
					err := models.GET(k, "id="+id, *item)
					if err == nil {
						s[index] = item
					}
				}
				(*resq.Included)[k] = ins
			} else {
				ins.T = postgres.Map
				if v.V == "" {
					continue
				}
				item := &map[string]*postgres.V{}
				for k, d := range m {
					(*item)[k] = &postgres.V{
						T: d.T,
					}
				}
				err := models.GET(k, "id="+v.V, *item)
				if err == nil {
					ins.Map = item
				}
			}
		}
	}

	resq.Type = req.Type
	resq.Data = &types.RESPDATA{
		T:   postgres.Map,
		Map: &vMap,
	}
	return nil
}

// DELETE 删除
func DELETE(req *types.REQ, resq *types.RESP) error {
	err := models.DB.Del(req.Type, "id='"+req.Query.ID+"'")
	if err != nil {
		er := errors.New(errors.INVALID, "资源删除出错")
		er.LogMessage(err.Error())
		return er
	}
	resq.Type = req.Type
	resq.Data = &types.RESPDATA{
		T: postgres.Map,
		Map: &map[string]*postgres.V{
			"id": &postgres.V{
				T: postgres.Int8,
				V: req.Query.ID,
			},
		},
	}
	return nil
}
