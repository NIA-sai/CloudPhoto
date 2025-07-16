package aitool

import (
	"CloudPhoto/config"
	"CloudPhoto/internal/model"
	"CloudPhoto/internal/storage"
	"CloudPhoto/internal/tool"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	cutOutUrl    string
	rawCutOutUrl *url.URL
)

func changeFace(c *gin.Context) {

	face, err := c.FormFile("face")
	if err != nil {
		print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing face image"})
		return
	}
	bodyId := c.Query("body_id")
	if bodyId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing base_id"})
		return
	}
	taskId := uuid.New().String()
	storage.SetChangeFaceTask(
		taskId,
		model.TaskResult{
			Status: "processing",
		})
	go func() {
		bodyFile, err := os.Open(storage.GetBodyFilePath(bodyId))
		tool.PanicIfErr(err)
		defer tool.PanicIfErr(bodyFile.Close())
		bodyByte, err := io.ReadAll(bodyFile)
		tool.PanicIfErr(err)
		faceFile, err := face.Open()
		tool.PanicIfErr(err)
		defer tool.PanicIfErr(faceFile.Close())
		faceByte, err := io.ReadAll(faceFile)
		tool.PanicIfErr(err)
		result, err := callPhotosAiApi(storage.ChangeFace, &bodyByte, &faceByte)
		if err != nil {
			print(err.Error())
			storage.SetChangeFaceTask(taskId, model.TaskResult{
				Status: "failed",
			})
			return
		} else {
			storage.SetChangeFaceTask(taskId, model.TaskResult{
				Status: "success",
				Result: result,
			})
		}
	}()
	c.String(202, "%s", taskId)
}

func cutOutFigure(c *gin.Context) {
	files, err := c.MultipartForm()
	if err != nil {
		print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing files"})
		return
	}
	figures, ok := files.File["figures"]
	if !ok || len(figures) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing figures"})
		return
	} else if len(figures) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "too many figures"})
		return
	}
	taskIds := ""
	//taskIdss := make([]string, len(figures))
	for _, f := range figures {
		taskId := uuid.New().String()
		taskIds += taskId + ","
		//taskIdss[i] = taskId
		go func() {
			figureFile, err := f.Open()
			if err != nil {
				print(err.Error())
				return
			}
			defer tool.PanicIfErr(figureFile.Close())
			data, err := io.ReadAll(figureFile)
			if err != nil {
				print(err.Error())
				return
			}
			result, err := callPhotosAiApi(storage.CutOut, &data)
			//可能不必
			if err != nil {
				print(err.Error())
				storage.SetCutOutTask(taskId, model.TaskResult{
					Status: "failed",
					Result: result,
				})
				return
			} else {
				storage.SetCutOutTask(taskId, model.TaskResult{
					Status: "complete",
					Result: result,
				})
			}
		}()
	}
	c.String(202, "%s", taskIds)
	//c.JSON(202, taskIdss)
}

func callPhotosAiApi(api string, photo *[]byte, photos ...*[]byte) (any, error) {
	return func(f *[]byte, fs ...*[]byte) (any, error) {
		switch api {
		case storage.CutOut:
			{
				all := make([]*[]byte, 0, len(photos)+1)
				all = append(all, photo)
				all = append(all, photos...)
				body := getCutOutReqBaseBody()
				results := make([]string, 0, len(all))
				var tmp int

				for _, p := range all {
					body.BinaryDataBase64[0] = base64.StdEncoding.EncodeToString(*p)
					tmp = tool.SendHttpReq(getCutOutBaseReq(body),
						func(resp *http.Response) {
							if resp.StatusCode != 200 {
								fmt.Println(resp.Body)
								return
							}
							tool.PanicIfErr(json.NewDecoder(resp.Body).Decode(&result))
						},
					)
					if tmp != 200 {
						//存在部分错误，信息
						//在里面处理？
						results = append(results, "failed")
					} else {

						results = append(results, result.Data.BinaryDataBase64[0])

					}
				}
				return results, nil
			}
		default:
			return nil, errors.New("unknown api")
		}
	}(photo, photos...)
}

func getCutOutReqBaseBody() *cutOutReqBody {
	return &cutOutReqBody{
		ReqKey:           "saliency_seg",
		OnlyMask:         3,
		RefineMask:       0,
		RGB:              [3]int{-1, -1, -1},
		BinaryDataBase64: make([]string, 1),
		ReturnUrl:        true,
		LogoInfo: cutOutLogoInfo{
			AddLogo:         true,
			Position:        rand.Int() % 4,
			LogoTextContent: "山东大学学生在线",
		},
	}
}

type cutOutReqBody struct {
	ReqKey           string         `json:"req_key"`
	BinaryDataBase64 []string       `json:"binary_data_base64"`
	ReturnUrl        bool           `json:"return_url"`
	OnlyMask         int            `json:"only_mask"`
	RGB              [3]int         `json:"rgb"`
	RefineMask       int            `json:"refine_mask"`
	LogoInfo         cutOutLogoInfo `json:"logo_info"`
}
type cutOutLogoInfo struct {
	AddLogo         bool    `json:"add_logo"`
	Position        int     `json:"position"`
	Language        int     `json:"language"`
	Opacity         float64 `json:"opacity"`
	LogoTextContent string  `json:"logo_text_content"`
}

var result cutOutRespBody

type cutOutRespBody struct {
	Code int `json:"code"`
	Data struct {
		BinaryDataBase64 []string `json:"binary_data_base64"`
		ImageUrls        []string `json:"image_urls"`
	} `json:"data"`
}

// 并非Base
func getCutOutBaseReq(data *cutOutReqBody) *http.Request {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	//os.WriteFile("output.txt", []byte(jsonBytes), 0644)
	req, _ := http.NewRequest(config.Get().Ai.CutOut.Method, cutOutUrl, bytes.NewReader(jsonBytes))
	now := time.Now().UTC()
	xDate := now.Format("20060102T150405Z")
	xxDate := now.Format("20060102")
	//沟槽签名
	rawCredential := xxDate + "/" + config.Get().Ai.CutOut.Region + "/" + config.Get().Ai.CutOut.Service + "/request"
	credential := "Credential=" + config.Get().Ai.CutOut.AccessId + "/" + rawCredential
	signedHeaders := "SignedHeaders=x-date"

	requestPayload := tool.HexOfHash256(jsonBytes)
	canonicalRequest :=
		config.Get().Ai.CutOut.Method + "\n" +
			"/\n" +
			rawCutOutUrl.RawQuery + "\n" +
			"x-date:" + xDate + "\n\n" +
			"x-date" + "\n" +
			requestPayload
	stringToSign :=
		"HMAC-SHA256" + "\n" +
			xDate + "\n" +
			rawCredential + "\n" +
			tool.HexOfHash256([]byte(canonicalRequest))
	kSecret := config.Get().Ai.CutOut.SecretKey

	kDate := tool.HmacSHA256([]byte(kSecret), xxDate)
	kRegion := tool.HmacSHA256(kDate, config.Get().Ai.CutOut.Region)
	kService := tool.HmacSHA256(kRegion, config.Get().Ai.CutOut.Service)
	kSigning := tool.HmacSHA256(kService, "request")
	signature := hex.EncodeToString(tool.HmacSHA256(kSigning, stringToSign))
	authorization := "HMAC-SHA256 " + credential + ", " + signedHeaders + ", Signature=" + signature
	req.Header.Set("X-Date", xDate)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json")
	fmt.Println(authorization)
	return req
}
