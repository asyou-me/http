package errors

import (
	base_errors "github.com/asyou-me/lib.v1/errors"
)

// 错误类型
const (
	MISSING = iota
	INVALID
	SYSTEM_ERR
	INVALID_PARAME
	INVALID_DATA
	PARAME_REQUIRE
	NULL
)

// 定义错误容器集合
var codes = &base_errors.ErrCodes{
	SYSTEM_ERR: &base_errors.ErrStruct{
		Code: 500,
		Format: map[string][]string{
			"zh": []string{
				"系统出现 %s 错误: %s", "未知", "请联系系统维护者",
			},
			"en": []string{
				"system %s error: %s", "unknown", "Please contact the system administrator",
			},
		},
		Level:    "ERROR",
		ValueLen: 2,
		Type:     "system_error",
	}, MISSING: &base_errors.ErrStruct{
		Code: 404,
		Format: map[string][]string{
			"zh": []string{
				"无法找到 %s: %s", "", "",
			},
			"en": []string{
				"%s not found: %s", "", "",
			},
		},
		Level:    "INFO",
		ValueLen: 2,
		Type:     "miss_error",
	}, INVALID: &base_errors.ErrStruct{
		Code: 404,
		Format: map[string][]string{
			"zh": []string{
				"无效 %s: %s", "", "",
			},
			"en": []string{
				"invalid %s: %s", "", "",
			},
		},
		Level:    "INFO",
		ValueLen: 2,
		Type:     "invalid_request_error",
	}, INVALID_PARAME: &base_errors.ErrStruct{
		Code: 404,
		Format: map[string][]string{
			"zh": []string{
				"无效请求参数 %s: %s", "", "",
			},
			"en": []string{
				"invalid parameter %s: %s", "", "",
			},
		},
		Level:    "INFO",
		ValueLen: 2,
		Type:     "invalid_request_error",
	}, INVALID_DATA: &base_errors.ErrStruct{
		Code: 404,
		Format: map[string][]string{
			"zh": []string{
				"无法找到符合 %s 的 %s 记录", "", "",
			},
			"en": []string{
				"Unable to find %s with %s", "", "",
			},
		},
		Level:    "INFO",
		ValueLen: 2,
		Type:     "invalid_request_error",
	}, PARAME_REQUIRE: &base_errors.ErrStruct{
		Code: 404,
		Format: map[string][]string{
			"zh": []string{
				"缺少参数 %s", "",
			},
			"en": []string{
				"invalid parameter %s", "",
			},
		},
		Level:    "INFO",
		ValueLen: 1,
		Type:     "invalid_request_error",
	}, NULL: &base_errors.ErrStruct{
		Code: 404,
		Format: map[string][]string{
			"zh": []string{
				"%s", "",
			},
			"en": []string{
				"%s", "",
			},
		},
		Level:    "INFO",
		ValueLen: 1,
		Type:     "invalid_request_error",
	},
}
