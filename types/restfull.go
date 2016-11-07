package types

import (
	"errors"
	"fmt"

	"encoding/json"

	"github.com/asyou-me/lib.v1/json/lexer"
	"github.com/asyou-me/lib.v1/utils"
	"github.com/asyou-me/postgres"
)

// RestFullConf 资源配置文件
type RestFullConf struct {
	Name   string
	POST   *RestFull `json:"post,omitempty"`
	PATCH  *RestFull `json:"patch,omitempty"`
	GET    *RestFull `json:"get,omitempty"`
	LIST   *RestFull `json:"list,omitempty"`
	DELETE *RestFull `json:"delete,omitempty"`
}

// RestFull 资源定义
type RestFull struct {
	Msg    *string
	Req    *ReqDefine
	Resq   *ReqDefine
	Delete string
	Auth   []string
	Event  Event
}

// Event 事件
type Event struct {
	Type string
}

// REQ json 请求通用结构体
type REQ struct {
	Query         *Query
	Type          string
	Attributes    *map[string]*postgres.V
	Relationships *map[string]*Relationship
}

// Query 查询构造器
type Query struct {
	ID     string
	Fields map[string]uint8
	Filter map[string]string // {"id":"=1","&time":">1"} '&' 为查询条件的 or
	Order  []string          // {"-time","id"}
	Limit  int
	Start  int
}

// RESP json返回通用结构体
type RESP struct {
	Msg      *string                   `json:"msg,omitempty"`
	Type     string                    `json:"type,omitempty"`
	Data     *RESPDATA                 `json:"data,omitempty"`
	Included *map[string]*Relationship `json:"included,omitempty"`
}

// RESPDATA json返回通用结构体
type RESPDATA struct {
	T     uint8
	Map   *map[string]*postgres.V
	Slice *[]*map[string]*postgres.V
}

// MarshalJSON 序列化时调用
func (r *RESPDATA) MarshalJSON() ([]byte, error) {
	if r.T == postgres.Slice {
		return json.Marshal(r.Slice)
	}
	return json.Marshal(r.Map)
}

// Relationship 相关内容
type Relationship struct {
	T     uint8
	Map   *map[string]*postgres.V
	Slice *[]*map[string]*postgres.V
}

// MarshalJSON 序列化时调用
func (r *Relationship) MarshalJSON() ([]byte, error) {
	if r.T == postgres.Slice {
		return json.Marshal(r.Slice)
	}
	return json.Marshal(r.Map)
}

// Array 获取插入到数据库的内容
func (r *Relationship) Array() *[]int64 {
	datas := *r.Slice
	b := make([]int64, len(datas))
	bp := 0
	for _, v := range datas {
		b[bp] = utils.StringToInt64((*v)["id"].V, 0)
		bp = bp + 1
	}
	return &b
}

// ReqDefine 资源内容定义
type ReqDefine struct {
	Attributes    *map[string]*Define
	Relationships *map[string]*RelationshipDefine
}

// RelationshipDefine 关联资源定义
type RelationshipDefine struct {
	T      uint8
	Auth   []string
	Define *map[string]*Define
}

// Define 资源内容字段定义
type Define struct {
	T         uint8
	Require   bool
	Validates []string
}

// Read 读取 json 字节
func (r *REQ) Read(data []byte, define *ReqDefine) error {
	in := &lexer.Lexer{Data: data}
	if in.IsNull() {
		in.Skip()
		return errors.New("请传入正常 json 格式数据")
	}
	if !in.IsDelim('{') {
		return errors.New("请传入正常 json 格式数据")
	}
	in.Delim('{')
	curDef := map[string]struct{}{
		"type":       struct{}{},
		"attributes": struct{}{},
	}
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "type":
			str := in.String()
			if str == "" {
				return errors.New("type 必须为字符串")
			}
			r.Type = str
			delete(curDef, "type")
		case "attributes":
			if define.Attributes == nil {
				in.SkipRecursive()
				break
			}
			if in.IsNull() {
				in.Skip()
				return errors.New("attributes 不能为 nil")
			}
			r.Attributes = &map[string]*postgres.V{}
			err := readMap("attributes", in, r.Attributes, *define.Attributes)
			if err != nil {
				return err
			}
			delete(curDef, "attributes")
		case "relationships":
			if define.Relationships == nil {
				in.SkipRecursive()
				break
			}
			if in.IsNull() {
				in.Skip()
				return errors.New("relationships 不能为 nil")
			}
			r.Relationships = &map[string]*Relationship{}
			err := read("relationships", in, r.Relationships, *define.Relationships)
			if err != nil {
				return err
			}
			delete(curDef, "relationships")
		default:
			in.SkipRecursive()
			return errors.New("未知参数：" + key)
		}
		in.WantComma()
	}
	in.Delim('}')
	if len(curDef) > 0 {
		str := "["
		for k := range curDef {
			if len(str) != 1 {
				str = str + ","
			}
			str = str + k
		}
		str = str + "]"
		return errors.New(str + "为必传参数")
	}
	return nil
}

func read(path string, in *lexer.Lexer, data *map[string]*Relationship, define map[string]*RelationshipDefine) error {
	if in.IsNull() {
		in.Skip()
		return errors.New(path + " 必须是 json 格式")
	}
	if !in.IsDelim('{') {
		return errors.New(path + " 必须是 json 格式")
	}
	in.Delim('{')

	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		if key == "" {
			in.Skip()
			in.WantComma()
			continue
		}
		def, ok := define[key]
		if !ok {
			return errors.New(path + "." + key + " 为多余参数")
		}
		if def.T == postgres.Map {
			if !in.IsDelim('{') {
				return errors.New(path + "." + key + " 必须为对象")
			}
			d := &map[string]*postgres.V{}
			err := readMap(path+"."+key, in, d, *def.Define)
			if err != nil {
				return err
			}
			(*data)[key] = &Relationship{
				T:   postgres.Map,
				Map: d,
			}
		} else if def.T == postgres.Slice {
			if !in.IsDelim('[') {
				return errors.New(path + "." + key + " 必须为数组")
			}
			d := &[]*map[string]*postgres.V{}
			err := readSlice(path+"."+key, in, d, *def.Define)
			if err != nil {
				return err
			}
			(*data)[key] = &Relationship{
				T:     postgres.Slice,
				Slice: d,
			}
		} else {
			return errors.New(path + "." + key + " 必须为数组或对象")
		}
		in.WantComma()
	}

	in.Delim('}')
	return nil
}

func readMap(path string, in *lexer.Lexer, data *map[string]*postgres.V, define map[string]*Define) error {
	if in.IsNull() {
		in.Skip()
		return errors.New(path + " 必须是 json 格式")
	}
	if !in.IsDelim('{') {
		return errors.New(path + " 必须是 json 格式")
	}
	in.Delim('{')

	// 这里必须这样做 因为 map 的引用传值特性
	curDef := map[string]*Define{}
	for k := range define {
		curDef[k] = define[k]
	}
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		def, ok := define[key]
		if !ok {
			return errors.New("未知参数：" + path + "." + key)
		}

		v, err := getValue(in, def.T)
		if err != nil {
			return errors.New("参数：" + path + "." + key + err.Error())
		}

		(*data)[key] = &postgres.V{
			T: def.T,
			V: v,
		}
		delete(curDef, key)
		in.WantComma()
	}
	if len(curDef) > 0 {
		num := 0
		str := "["
		for k := range curDef {
			if curDef[k].Require == false {
				continue
			}
			if len(str) != 1 {
				str = str + ","
			}
			str = str + k
			num = num + 1
		}
		str = str + "]"
		if num > 0 {
			return errors.New(path + "." + str + "为必传参数")
		}
	}
	in.Delim('}')
	return nil
}

func readSlice(path string, in *lexer.Lexer, data *[]*map[string]*postgres.V, define map[string]*Define) error {
	in.Delim('[')
	var i = 0
	for !in.IsDelim(']') {
		d := &map[string]*postgres.V{}
		err := readMap(path+"["+fmt.Sprint(i)+"]", in, d, define)
		if err != nil {
			return err
		}
		(*data) = append(*data, d)
		i++
		in.WantComma()
	}
	in.Delim(']')
	return nil
}

func getValue(in *lexer.Lexer, t uint8) (string, error) {
	switch t {
	case postgres.Int, postgres.Int8, postgres.Int16, postgres.Int32, postgres.Int64:
		str, ok := in.Mumber()
		if !ok {
			return "", errors.New("类型应为数字")
		}
		return str, nil
	case postgres.String:
		str, ok := in.Str()
		if !ok {
			return "", errors.New("类型应为 string")
		}
		return str, nil
	case postgres.Bool:
		str, ok := in.Boolean()
		if !ok {
			return "", errors.New("类型应为 bool")
		}
		return str, nil
	}
	return "", errors.New("类型没有定义")
}
