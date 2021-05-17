package handlers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/nju-iot/security-certificate/logs"
)

var prvKey2 []byte
var pubKey1 []byte
var pubKey2 []byte

func readFromFile() {
	data, err := ioutil.ReadFile("keys/prvKey2.txt")
	if err != nil {
		logs.Error("%v", err)
	}
	prvKey2 = data
	data, err = ioutil.ReadFile("keys/pubKey1.txt")
	if err != nil {
		logs.Error("%v", err)
	}
	pubKey1 = data
	// pubKey1 = ([]byte)("MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCulH3yfmzpG8N5QJQttaNMW9Jw\nnQovz4xUYuXNEf8XeYR+l07nLSBJBHY2qFFBOp/0zzGC7Yz1ZxMJOx3ojw+Nnflb\nRR+oRIGed23WigKJrja5uSWg4CGic7OvkV6DKcrkWVhDnVlBnoblrg20hmpNHCO+\nGWWOolPyXyjeXZj26wIDAQAB")
	data, err = ioutil.ReadFile("keys/pubKey2.txt")
	if err != nil {
		logs.Error("%v", err)
	}
	pubKey2 = data
}

func EncryptData(data []byte) []byte {
	if prvKey2 == nil || pubKey2 == nil || pubKey1 == nil {
		readFromFile()
	}
	return rsaEncrypt(data, pubKey1)
}

func testEncryptData(data []byte) []byte {
	if prvKey2 == nil || pubKey2 == nil || pubKey1 == nil {
		readFromFile()
	}
	return rsaEncrypt(data, pubKey2)
}

func DecryptData(data []byte) []byte {
	if prvKey2 == nil || pubKey2 == nil || pubKey1 == nil {
		readFromFile()
	}
	return rsaDecrypt(data, prvKey2)
}

func CheckSign(data, sign, key []byte) bool {
	return rsaVerySignWithSha256(data, sign, pubKey1)
}

//签名
func rsaSignWithSha256(data []byte, keyBytes []byte) []byte {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error"))
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		panic(err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		panic(err)
	}

	return signature
}

//验证
func rsaVerySignWithSha256(data, signData, keyBytes []byte) bool {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signData)
	if err != nil {
		panic(err)
	}
	return true
}

// 公钥加密
func rsaEncrypt(data, keyBytes []byte) []byte {
	//解密pem格式的公钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		panic(err)
	}
	return ciphertext
}

// 私钥解密
func rsaDecrypt(ciphertext, keyBytes []byte) []byte {
	//获取私钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error"))
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// 解密
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		panic(err)
	}
	return data
}
