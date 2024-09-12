package tool

import (
	"crypto/sha512"
	"fmt"

	"github.com/anaskhan96/go-password-encoder"
)

var options = &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New} //生成加密方法

func Salt_Encryption(p string) string {
	salt, encodedPwd := password.Encode(p, options)
	back := fmt.Sprintf("sha512$%s%s", salt, encodedPwd) //将盐值和生成的密钥保存//一个是$之前的剩下的一个是16个一个是64个
	return back
}

// 前边是密码  后边是加密后的密码
func CHECK_Password(p1 string, p2 string) bool {
	//先取值
	salt := p2[7:23]
	encodedPwd := p2[23:]
	check := password.Verify(p1, salt, encodedPwd, options)
	return check
}
