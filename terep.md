部分接口说明(所有图片返回暂时为base64编码)：


|接口| 请求方法 | 请求体                                            | 请求头                                              | 返回内容说明                                              |
|-----|------|------------------------------------------------|--------------------------------------------------|-----------------------------------------------------|
| /captcha/image | GET  | 无                                              | 无                                                | JSON:<br />图片验证码(captchaImage)<br />及其ID(captchaId) |
| /api/cutOutFigure | POST | formdata:<br />figures[File]: 需要抠图的文件(可一次性传多个) | "captcha-id":图片验证码id<br />"captcha-code":图片验证码答案 | String:直接就是任务id(taskId) (逗号分隔，一张图一个task)            |
| /task/ask/cutOut/{taskId} | GET  | 无                                              | 无                                                | JSON:<br />status(状态)<br />result(结果, 字符串数组)        |
| /api/fuseFace        | POST    | 见下方示例                                          | 像扣图一样的验证码请求                                      | 见下方示例                                               |
| /api/getModelList            | POST     |    见下方示例                                            |       像扣图一样的验证码请求                                           |                   见下方示例                                  |


/api/fuseFace人脸融合接口示例：

    {
        "model_id": "mt_1947953263779856384", // string，必填，模板素材ID，也就是换脸底图的id
        "rsp_img_type": "url", // string，必填，返回图像格式，可选值："url" 或 "base64"。如果要返回url，有效期是7天，及时下载下来就好了
        "merge_infos": [ // 必填，融合信息数组，这里是可以把多张图片人脸融合在一起，但是咱要换脸，所以只需要传一张用户的图片。
            {
                "image": "/9j/4AAQSkZJRgABAQEBLAEsAAD/4QCORXhpZgAATU0AKgAAAAgAAgESAAMAAAABAAEA" //string，输入用户图片的base64编码，无前缀
            }
        ]
    }
成功响应(200)：

    {
        "success": true,
        "fused_url": "string"     // 返回换脸后的图片URL或base64数据(url或base64是在发请求的时候就选好了的)，建议及时把图片下载下来
    }
失败响应(400/500)：

    {
        "success": false,
        "error": "string",        // 错误描述
        "error_code": "string"    // 错误代码（可选）
    }

/api/getModelList获取人脸融合素材模板列表接口示例：
    
    //这里是分页返回，防止模板太多一次性返回不方便
    {
        "limit": 20,              // int,单次返回的模板素材数量限制，默认是20，最大能填20
        "offset": null               // int,偏移量，用于分页，默认是0
    }
成功响应(200)：

    {
        "data": {
            "MaterialInfos": [
                {
                    "MaterialId": "mt_1947668117646295040",
                    "MaterialName": "测试1.jpg",
                },
                {
                    "MaterialId": "mt_1947953263779856384",
                    "MaterialName": "测试1模板.jpg",
                }
            ],
            "Count": 2, //素材模板条数
            "RequestId": "cd29dd64-57f2-434d-82a5-4d780300d490" //唯一请求 ID，由服务端生成，每次请求都会返回(若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId)。定位问题时需要提供该次请求的 RequestId。
        },
        "success": true
    }

失败响应(400/500)：

    {
        "error": "string"         // 错误描述
    }