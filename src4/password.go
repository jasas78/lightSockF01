//package lightsocks
package main

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"strings"
	"time"
)

const _VpasswordLength = 256

type _STpassword [_VpasswordLength]byte

func _Finit1() {
	// 更新随机种子，防止生成一样的随机密码
	rand.Seed(time.Now().Unix())
}

// 采用base64编码把密码转换为字符串
func (___Vpassword *_STpassword) String() string {
	return base64.StdEncoding.EncodeToString(___Vpassword[:])
}

// 解析采用base64编码的字符串获取密码
func _FparsePassword(passwordString string) (*_STpassword, error) {
	__Vbs, __Verr := base64.StdEncoding.DecodeString(strings.TrimSpace(passwordString))
	if __Verr != nil || len(__Vbs) != _VpasswordLength {
		return nil, errors.New("不合法的密码")
	}
	__Vpassword := _STpassword{}
	copy(__Vpassword[:], __Vbs)
	__Vbs = nil
	return &__Vpassword, nil
}

// 产生 256个byte随机组合的 密码，最后会使用base64编码为字符串存储在配置文件中
// 不能出现任何一个重复的byte位，必须又 0-255 组成，并且都需要包含
func _FrandPassword() string {
	// 随机生成一个由  0~255 组成的 byte 数组
	__VintArr := rand.Perm(_VpasswordLength)
	__Vpassword := &_STpassword{}
	for __Vi, __Vv := range __VintArr {
		__Vpassword[__Vi] = byte(__Vv)
		if __Vi == __Vv {
			// 确保不会出现如何一个byte位出现重复
			return _FrandPassword()
		}
	}
	return __Vpassword.String()
}
