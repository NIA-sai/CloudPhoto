package tool

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	_ "fmt"
	"net/http"
)

var httpClient = &http.Client{}

func SendHttpReq(req *http.Request, f func(*http.Response)) int {
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	f(resp)
	return resp.StatusCode
}

func PanicIfErr(err ...error) {
	for _, e := range err {
		if e != nil {
			print(e.Error())
			panic(e)
		}
	}
}
func HandleErr(err error, f func()) {
	if err != nil {
		print(err.Error())
		f()
	}
}
func HexOfHash256(data []byte) string {
	sha := sha256.Sum256(data)
	return hex.EncodeToString(sha[:])
}
func HmacSHA256(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}
