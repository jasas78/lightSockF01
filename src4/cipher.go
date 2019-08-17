//package lightsocks
package main

type _STcipher struct {
	// 编码用的密码
	encodePassword *_STpassword
	// 解码用的密码
	decodePassword *_STpassword
}

// 加密原数据
func (__Vcipher *_STcipher) _Fencode(bs []byte) {
	for i, v := range bs {
		bs[i] = __Vcipher.encodePassword[v]
	}
}

// 解码加密后的数据到原数据
func (__Vcipher *_STcipher) _Fdecode(bs []byte) {
	for i, v := range bs {
		bs[i] = __Vcipher.decodePassword[v]
	}
}

// 新建一个编码解码器
func newCipher(encodePassword *_STpassword) *_STcipher {
	decodePassword := &_STpassword{}
	for i, v := range encodePassword {
		encodePassword[i] = v
		decodePassword[v] = byte(i)
	}
	return &_STcipher{
		encodePassword: encodePassword,
		decodePassword: decodePassword,
	}
}
