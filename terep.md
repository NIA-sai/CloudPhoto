部分接口说明(所有图片返回暂时为base64编码)：


|接口| 请求方法 | 请求体 | 请求头 | 返回内容说明                                              |
|-----| -------- | ------ | ------ |-----------------------------------------------------|
| /captcha/image | GET | 无 | 无 | JSON:<br />图片验证码(captchaImage)<br />及其ID(captchaId) |
| /api/cutOutFigure | POST | formdata:<br />figures[File]: 需要抠图的文件(可一次性传多个) | "captcha-id":图片验证码id<br />"captcha-code":图片验证码答案 | String:直接就是任务id(taskId) (逗号分隔，一张图一个task)            |
| /task/ask/cutOut/{taskId} | GET | 无 | 无 | JSON:<br />status(状态)<br />result(结果, 字符串数组)        |
|          |        |        |              |                                                     |

