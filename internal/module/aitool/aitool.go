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
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	facefusion "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/facefusion/v20220927"
)

var (
	cutOutUrl    string
	rawCutOutUrl *url.URL
)

//func changeFace(c *gin.Context) {
//
//	face, err := c.FormFile("face")
//	if err != nil {
//		print(err.Error())
//		c.JSON(http.StatusBadRequest, gin.H{"error": "missing face image"})
//		return
//	}
//	bodyId := c.Query("body_id")
//	if bodyId == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "missing base_id"})
//		return
//	}
//	taskId := uuid.New().String()
//	storage.SetChangeFaceTask(
//		taskId,
//		model.TaskResult{
//			Status: "processing",
//		})
//	go func() {
//		bodyFile, err := os.Open(storage.GetBodyFilePath(bodyId))
//		tool.PanicIfErr(err)
//		defer tool.PanicIfErr(bodyFile.Close())
//		bodyByte, err := io.ReadAll(bodyFile)
//		tool.PanicIfErr(err)
//		faceFile, err := face.Open()
//		tool.PanicIfErr(err)
//		defer tool.PanicIfErr(faceFile.Close())
//		faceByte, err := io.ReadAll(faceFile)
//		tool.PanicIfErr(err)
//		result, err := callPhotosAiApi(storage.ChangeFace, &bodyByte, &faceByte)
//		if err != nil {
//			print(err.Error())
//			storage.SetChangeFaceTask(taskId, model.TaskResult{
//				Status: "failed",
//			})
//			return
//		} else {
//			storage.SetChangeFaceTask(taskId, model.TaskResult{
//				Status: "success",
//				Result: result,
//			})
//		}
//	}()
//	c.String(202, "%s", taskId)
//}

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
						results = append(results, "cutOut failed")
					} else {
						data, err := base64.StdEncoding.DecodeString(result.Data.BinaryDataBase64[0])
						if err != nil {
							results = append(results, "cutOut failed")
						} else {
							tool.PanicIfErr(os.MkdirAll(config.Get().App.StaticRoot+storage.CutOutFilePath, os.ModePerm))
							fileName := "cutOuted_" + result.RequestId + ".png"
							fullPath := filepath.Join(config.Get().App.StaticRoot+storage.CutOutFilePath, fileName)

							if err := os.WriteFile(fullPath, data, 0644); err != nil {
								results = append(results, "cutOut failed")
							} else {
								results = append(results, config.Get().App.Domain+"/"+filepath.Join(config.Get().App.StaticRelativePath+storage.CutOutFilePath, fileName))
							}
						}

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
	}
}

type cutOutReqBody struct {
	ReqKey           string   `json:"req_key"`
	BinaryDataBase64 []string `json:"binary_data_base64"`
	ReturnUrl        bool     `json:"return_url"`
	OnlyMask         int      `json:"only_mask"`
	RGB              [3]int   `json:"rgb"`
	RefineMask       int      `json:"refine_mask"`
}

var result cutOutRespBody

type cutOutRespBody struct {
	Code int `json:"code"`
	Data struct {
		BinaryDataBase64 []string `json:"binary_data_base64"`
		ImageUrls        []string `json:"image_urls"`
	} `json:"data"`
	RequestId string `json:"request_id"`
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
	return req
}

//人脸融合:

// 定义融合请求结构体，用于向后端、向腾讯云发请求
type FusionRequest struct {
	ModelId    string      `json:"model_id"`     // 模板素材ID
	RspImgType string      `json:"rsp_img_type"` // 返回图像格式(url/base64)
	MergeInfos []MergeInfo `json:"merge_infos"`  // 融合信息数组
}

// 定义MergeInfo结构体
type MergeInfo struct {
	//选择url/base64一种提交就行
	Image string `json:"image,omitempty"` // 输入图片base64 (可选)
	Url   string `json:"url,omitempty"`   // 输入图片URL (可选)
}

// 定义响应结构体,用于返回前端
type FusionResponse struct {
	Success   bool   `json:"success"`
	FusedURL  string `json:"fused_url,omitempty"`
	Error     string `json:"error,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
}

func fuseFace(c *gin.Context) { //前端调用这个函数，这里面调用callFaceFusionApi函数

	// 声明请求结构体实例
	var req FusionRequest

	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式: " + err.Error()})
		return
	}

	// 验证必要参数
	if req.ModelId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必填参数: model_id"})
		return
	}

	if req.RspImgType != "url" && req.RspImgType != "base64" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rsp_img_type必须是'url'或'base64'"})
		return
	}

	if len(req.MergeInfos) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "至少需要提供一个MergeInfo"})
		return
	}

	// 调用人脸融合API
	result, err := callFaceFusionApi(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, FusionResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// 返回成功结果
	c.JSON(http.StatusOK, FusionResponse{
		Success:  true,
		FusedURL: result,
	})
}

// 腾讯云人脸融合API调用函数
func callFaceFusionApi(req FusionRequest) (string, error) {
	// 1获取配置
	cfg := config.Get().Ai.FaceFusion

	// 2创建凭证对象
	credential := common.NewCredential(
		cfg.SecretId,
		cfg.SecretKey,
	)

	// 3创建客户端配置
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = cfg.Url

	// 4创建客户端
	client, err := facefusion.NewClient(credential, cfg.Region, cpf)
	if err != nil {
		return "", fmt.Errorf("failed to create client: %v", err)
	}

	// 5创建请求对象
	request := facefusion.NewFuseFaceRequest()

	// 6设置请求参数
	request.ProjectId = common.StringPtr(cfg.ProjectId)
	request.ModelId = common.StringPtr(req.ModelId)       // 模板素材ID
	request.RspImgType = common.StringPtr(req.RspImgType) // 返回图像格式(url或base64)

	// 转换 MergeInfos 从 []MergeInfo 到 []*MergeInfo
	mergeInfos := make([]*facefusion.MergeInfo, len(req.MergeInfos))
	for i, info := range req.MergeInfos {
		mergeInfos[i] = &facefusion.MergeInfo{
			Image: common.StringPtr(info.Image),
			Url:   common.StringPtr(info.Url),
		}
	}
	request.MergeInfos = mergeInfos

	// 7发送请求
	response, err := client.FuseFace(request)
	if err != nil {
		if sdkErr, ok := err.(*tcerr.TencentCloudSDKError); ok {
			return "", fmt.Errorf("API error: %s, %s", sdkErr.Code, sdkErr.Message)
		}
		return "", fmt.Errorf("failed to fuse face: %v", err)
	}

	// 8处理响应
	if response.Response != nil && response.Response.FusedImage != nil {
		return *response.Response.FusedImage, nil
	}

	return "", errors.New("no fused image in response")
}
