//package lightsocks
package main

import (
    "fmt"
)

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
func _Fdump256byte(___VpreStr, ___VpostStr string, ___VstPass *_STpassword) {
	fmt.Printf(___VpreStr)
	for __Vi1 := 0; __Vi1 < 256; __Vi1++ {
		fmt.Printf("%2x", ___VstPass[__Vi1])
		if __Vi1 != 255 && (__Vi1&7 == 7) {
			fmt.Printf(" ")
		}
	}
	fmt.Printf(___VpostStr)
} // _Fdump256byte

// 新建一个编码解码器
func _FnewCipher(___VencodePassword *_STpassword) *_STcipher {
	__VdecodePassword := &_STpassword{}
	for __Vi3, __Vv3 := range ___VencodePassword {
		//___VencodePassword[__Vi3] = __Vv3
		__VdecodePassword[__Vv3] = byte(__Vi3)
	}
	_Fdump256byte("enPass : ", "\n", ___VencodePassword)
	_Fdump256byte("dePass : ", "\n", __VdecodePassword)

	return &_STcipher{
		encodePassword: ___VencodePassword,
		decodePassword: __VdecodePassword,
	}
}
