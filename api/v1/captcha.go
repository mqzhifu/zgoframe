package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"zgoframe/core/global"
	httpresponse "zgoframe/http/response"
)

var store = base64Captcha.DefaultMemStore

// @Tags Base
// @Summary 生成图片验证码
// @Description 防止有人恶意攻击，尝试破解密码
// @Security ApiKeyAuth
// @accept application/json
// @Param X-Source-Type header string true "来源" default(1)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY"
// @Produce application/json
// @Success 200 {object} httpresponse.Response
// @Router /base/captcha [get]
func Captcha(c *gin.Context) {
	//字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(global.C.Captcha.ImgHeight, global.C.Captcha.ImgWidth, global.C.Captcha.NumberLength, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	if id, b64s, err := cp.Generate(); err != nil {
		global.V.Zap.Error("验证码获取失败!", zap.Any("err", err))
		httpresponse.FailWithMessage("验证码获取失败", c)
	} else {
		httpresponse.OkWithDetailed(httpresponse.SysCaptchaResponse{
			CaptchaId: id,
			PicPath:   b64s,
		}, "验证码获取成功", c)
	}
}

// @Success 200 {string} string "{"success":true,"data":{},"msg":"验证码获取成功"}"
