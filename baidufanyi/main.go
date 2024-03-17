package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"net/http"
	"net/url"
)

// 官方文档：https://api.fanyi.baidu.com/doc/21

const apiURL = "https://fanyi-api.baidu.com/api/trans/vip/translate"

type result struct {
	Source      string `json:"src"` // 原文
	Destination string `json:"dst"` // 译文
}

type response struct {
	From      string   `json:"from"`         // 原语言
	To        string   `json:"to"`           // 目标语言
	Results   []result `json:"trans_result"` // 翻译结果
	ErrorCode int      `json:"error_code"`   // 错误码
}

// 生成签名
func signature(appID, secret, content, salt string) string {
	data := md5.Sum([]byte(appID + content + salt + secret))

	return hex.EncodeToString(data[:])
}

// 从响应体获取翻译结果
func getTranslation(body []byte) string {
	resp := response{}
	err := json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatal(err)
	}
	if resp.ErrorCode != 0 {
		log.Fatalf("translation error code: %d", resp.ErrorCode)
	}
	if len(resp.Results) == 0 {
		log.Fatal("translation is empty")
	}

	return resp.Results[0].Destination
}

// 请求翻译
func requestTranslation(appID, secret, from, to, content string) string {
	// 生成随机数作为 salt
	saltNum, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		log.Fatalf("random number generating error: %v", err)
	}
	salt := saltNum.String()

	sign := signature(appID, secret, content, salt)

	form := url.Values{}
	form.Set("q", content)
	form.Set("from", from)
	form.Set("to", to)
	form.Set("appid", appID)
	form.Set("salt", salt)
	form.Set("sign", sign)
	resp, err := http.PostForm(apiURL, form)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return getTranslation(body)
}

func main() {
	appID := flag.String("appid", "", "百度翻译的 APP ID（必需）")
	secret := flag.String("secret", "", "百度翻译的密钥（必需）")
	from := flag.String("from", "auto", "翻译的原语言（具体见<https://api.fanyi.baidu.com/doc/21>）")
	to := flag.String("to", "zh", "翻译的目标语言（具体见<https://api.fanyi.baidu.com/doc/21>）")
	helpShort := flag.Bool("h", false, "打印帮助信息")
	helpLong := flag.Bool("help", false, "打印帮助信息")

	flag.Parse()

	if flag.NFlag() == 0 || *helpShort || *helpLong {
		fmt.Println("baidufanyi [参数] 翻译内容")
		flag.PrintDefaults()

		return
	}

	if *appID == "" || *secret == "" {
		log.Fatal("需要设置 appid 和 secret")
	}
	if flag.NArg() == 0 {
		log.Fatal("翻译内容为空")
	}

	content := flag.Arg(0)

	translation := requestTranslation(*appID, *secret, *from, *to, content)
	if translation != "" {
		fmt.Println(translation)
	}
}
