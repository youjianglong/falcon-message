package sender

import (
	"testing"
	"github.com/youjianglong/falcon-message/config"
)

func getQCloudSMS() *QCloudSMS {
	return NewQCloudSMS(config.QCloudSms{APPID:"",APPKEY:"",SIGN:""})
}
func TestQCloudSMS_SendSms(t *testing.T) {
	t.Log(getQCloudSMS().SendSms("18718809947",`服务器127.0.0.1发生告警，问题：测试。
请及时检查并处理。`))
}
func TestQCloudSMS_SendSmsWithTplID(t *testing.T) {
	t.Log(getQCloudSMS().SendSmsWithTpl("18718809947",[]string{"127.0.0.1","测试"}),512)
}
func TestQCloudSMS_SendVoice(t *testing.T) {
	t.Log(getQCloudSMS().SendVoice("18718809947","服务器127.0.0.1发生告警，问题：测试。请及时检查并处理。"))
}