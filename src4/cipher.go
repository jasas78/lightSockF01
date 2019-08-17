//package lightsocks
package main

type _STcipher struct {
	// 编码用的密码
	encodePassword *_STpassword
	// 解码用的密码
	decodePassword *_STpassword
}

// 加密原数据
func (__Vcipher *_STcipher) _Fencode(__Vbs1 []byte) {
	for __Vi1, __Vv1 := range __Vbs1 {
		__Vbs1[__Vi1] = __Vcipher.encodePassword[__Vv1]
	}
}

// 解码加密后的数据到原数据
func (__Vcipher *_STcipher) _Fdecode(__Vbs2 []byte) {
	for __Vi2, __Vv2 := range __Vbs2 {
		__Vbs2[__Vi2] = __Vcipher.decodePassword[__Vv2]
	}
}

// 新建一个编码解码器
func _FnewCipher(___VencodePassword *_STpassword) *_STcipher {
	__VdecodePassword := &_STpassword{}
	for __Vi3, __Vv3 := range ___VencodePassword {
		//___VencodePassword[__Vi3] = __Vv3
		__VdecodePassword[__Vv3] = byte(__Vi3)
	}
	return &_STcipher{
		encodePassword: ___VencodePassword,
		decodePassword: __VdecodePassword,
	}
}
