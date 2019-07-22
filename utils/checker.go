package utils

// 校验证书
import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"regexp"
	"time"
)

func ParseCert(crt string) map[string]interface{} {
	var cert tls.Certificate
	//加载PEM格式证书到字节数组
	certPEMBlock := []byte(crt)
	//获取下一个pem格式证书数据 -----BEGIN CERTIFICATE-----   -----END CERTIFICATE-----
	certDERBlock, restPEMBlock := pem.Decode(certPEMBlock)
	if certDERBlock == nil {
		return nil
	}
	//附加数字证书到返回
	cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
	//继续解析Certifacate Chan,这里要明白证书链的概念
	certDERBlockChain, _ := pem.Decode(restPEMBlock)
	if certDERBlockChain != nil {
		//追加证书链证书到返回
		cert.Certificate = append(cert.Certificate, certDERBlockChain.Bytes)
		fmt.Println("存在证书链")
	}
	//第一个叶子证书就是我们https中使用的证书
	c, err := x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		fmt.Println("x509证书解析失败")
		return nil
	} else {
		m := make(map[string]interface{})
		// CRT证书文件信息：
		// 域名:	beta.webvpn.net.cn
		// SAN域名:	*.beta.webvpn.net.cn, beta.webvpn.net.cn
		// 公司/组织:	null
		// 部门:	null
		// 城市:	null
		// 省份:	null
		// 国家:	null
		// 有效期从:	2019-02-27 09:10:02
		// 有效期到:	2019-05-28 09:10:02
		// 过期剩余天数:	82
		// 序列号:	416839859003663373424025564381137678463467
		// 颁发者:	Let's Encrypt
		// 哈希签名算法:	SHA256
		// 加密算法:sha256WithRSAEncryption
		m["Organization"] = c.Issuer.Organization
		m["Subject"] = c.Subject.CommonName
		m["DNSNames"] = c.DNSNames
		m["NotBefore"] = c.NotBefore.Unix()
		m["NotAfter"] = c.NotAfter.Unix()
		m["Expire"] = (c.NotAfter.Unix() - strTimeYMDtoTimeStamp10(TimeStampToTimeStr(time.Now().Unix())[0:10])) //证书实际过期的秒
		return m
	}
}

// 年月日字符串转日期  例如:2019-05-23
func strTimeYMDtoTimeStamp10(str string) int64 {
	if len(str) != 10 {
		return 0
	}
	return TimeStrToTimeStamp10(str + " 15:04:05")
}
func TimeStrToTimeStamp10(t string) int64 {
	if len(t) < 18 {
		return 0
	}
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, t, loc)
	sr := theTime.Unix()
	return sr
}
func TimeStrToTime(t string) (ti time.Time, err error) {
	if len(t) < 18 {
		return ti, errors.New("wrong time str")
	}
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	ti, err = time.ParseInLocation(timeLayout, t, loc)
	return
}
func TimeStampToTimeStr(t int64) string {
	//  yyyy-MM-dd'T'HH:mm:ss'.000Z'
	timeLayout := "2006-01-02 15:04:05"
	dataTimeStr := time.Unix(t, 0).Format(timeLayout)
	// fmt.Println(dataTimeStr)
	return dataTimeStr
}

// 使用openssl解
// openssl x509 -in  /etc/letsencrypt/archive/newyingyong.cn/cert1.pem -noout -text

func CheckName(name string) bool {
	// reg := `^[a-zA-Z]{1}([a-zA-Z0-9]|[._]){4,19}$`
	// r, err := regexp.Compile(reg)
	// if err != nil {
	// 	return false
	// }
	// if r.MatchString(name) {
	// 	return true
	// }
	if len(name) >= 6 && len(name) <= 30 {
		return true
	}
	return false
}

func CheckTelephone(phone string) bool {
	if len(phone) != 11 {
		return false
	}
	reg := `^1([38]\d|5[0359]\d|7[3678]\d)\d{8}$`
	r, _ := regexp.Compile(reg)
	return r.MatchString(phone)
}

func CheckEmail(email string) bool {
	if len(email) == 0 || len(email) >= 255 {
		return false
	}
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegexp.MatchString(email) {
		return false
	}
	return true
}
