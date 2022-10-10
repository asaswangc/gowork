package utils

import (
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	numberBase = "1234567890"                                                               // 用于生成随机数字验证码
	letterBase = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^*()+" // 用于生成随机字符串
)

// https://segmentfault.com/a/1190000017346458
// 加解密算法：对称性加密算法、非对称性加密算法、散列算法，其中散列算法不可逆，无法解密，故而只能用于签名校验、身份验证
// 对称性加密算法：DES、3DES、AES
// 非对称性加密算法：RSA、DSA、ECC
// 散列算法：MD5、SHA1、HMAC

// GenerateRSAKey 生成私钥和公钥, bits参数指定证书大小
// 也可以直接通过openssl命令生成：
// 私钥：openssl genrsa -out rsa_private_key.pem 2048
// 公钥：openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
func GenerateRSAKey(bits int) error {
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(cryptoRand.Reader, bits)
	if err != nil {
		panic(err)
	}

	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//使用pem格式对x509输出的内容进行编码
	//创建文件保存私钥
	privateFile, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	defer func() {
		_ = privateFile.Close()
	}()

	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	//将数据保存到文件
	err = pem.Encode(privateFile, &privateBlock)
	if err != nil {
		return err
	}

	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}
	//pem格式编码
	//创建用于保存公钥的文件
	publicFile, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer func() {
		_ = publicFile.Close()
	}()

	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	//保存到文件
	err = pem.Encode(publicFile, &publicBlock)
	if err != nil {
		return err
	}
	return nil
}

// EncryptWithRSA rsa加密
func EncryptWithRSA(plainText string, publicKeyPath string) (string, error) {
	//打开公钥文件
	keyFile, err := os.Open(publicKeyPath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = keyFile.Close()
	}()

	//读取公钥内容
	info, _ := keyFile.Stat()
	buf := make([]byte, info.Size())
	_, err = keyFile.Read(buf)
	if err != nil {
		return "", err
	}

	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	cipherBytes, err := rsa.EncryptPKCS1v15(cryptoRand.Reader, publicKey, []byte(plainText))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(cipherBytes), nil
}

// DecryptWithRSA rsa解密
func DecryptWithRSA(cipherText string, privateKeyPath string) (string, error) {
	//打开私钥文件
	keyFile, err := os.Open(privateKeyPath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = keyFile.Close()
	}()

	//获取私钥内容
	info, _ := keyFile.Stat()
	buf := make([]byte, info.Size())
	_, err = keyFile.Read(buf)
	if err != nil {
		return "", err
	}

	//pem解码
	block, _ := pem.Decode(buf)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	//对密文进行解密
	cipherBytes, err := base64.RawURLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	plainBytes, err := rsa.DecryptPKCS1v15(cryptoRand.Reader, privateKey, cipherBytes)
	if err != nil {
		return "", err
	}
	return string(plainBytes), nil
}

// EncryptWithSha256 sha256加密
func EncryptWithSha256(data string) string {
	// 先base64对原始数据进行编码
	tmp := base64.StdEncoding.EncodeToString([]byte(data))

	// 再使用sha256进行两层加密
	h := sha256.New()
	h.Write([]byte(tmp))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

// GenRandomString 生成随机字符串或数字验证码，n表示要生成的字符串长度，code表示类型，1 数字验证码，2 随机字符串
// https://colobu.com/2018/09/02/generate-random-string-in-Go/
func GenRandomString(n int, code int) string {
	var letterBytes string
	if code == 1 {
		letterBytes = numberBase
	} else {
		letterBytes = letterBase
	}

	letterIdxBits := 6                    // 6 bits to represent a letter index
	letterIdxMask := 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax := 63 / letterIdxBits    // # of letter indices fitting in 63 bits
	src := rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)

	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache) & letterIdxMask; idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// CheckPassword 验证密码
func CheckPassword(password string) (b bool) {
	if ok, _ := regexp.MatchString("^[a-zA-Z\\d]{4,16}$", password); !ok {
		return false
	}
	return true
}

// CheckUsername 验证名称  必须是、下划线、@、.和6-10位之间的字母
func CheckUsername(username string) (b bool) {
	if ok, _ := regexp.MatchString("^[a-zA-Z\\d\u4e00-\u9fa5]{1,32}$", username); !ok {
		return false
	}
	return true
}

// Intersect 求交集
func Intersect(slice1, slice2 []int) []int {
	m := make(map[int]int)
	nn := make([]int, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// Difference 求差集 slice1-并集
func Difference(slice1, slice2 []int) []int {
	m := make(map[int]int)
	nn := make([]int, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

func CheckMenu(pageType int, path string) string {
	if pageType == 0 {
		return "m-" + GenerateFeatureLabel(path)
	}
	return GenerateFeatureLabel(path)
}

func GenerateFeatureLabel(url string) string {
	pt := `^/.*`
	buf := []byte(url)
	ok, _ := regexp.Match(pt, buf)
	println(ok)
	if ok {
		buf = buf[1:]
	}
	return strings.ReplaceAll(string(buf), "/", "-")
}

func StructToJson(Struct interface{}) (string, error) {
	bytes, err := json.Marshal(Struct)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
