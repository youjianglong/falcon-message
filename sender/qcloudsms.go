package sender

import (
	"github.com/youjianglong/falcon-message/config"
	"github.com/qichengzx/qcloudsms_go"
	"fmt"
	"strings"
)

type QCloudSMS struct {
	cfg    config.QCloudSms
	client *qcloudsms.QcloudSMS
}

//使用模板发送短信
func (this *QCloudSMS) SendSmsWithTpl(phoneNumber string, params []string, tpl ...int) error {
	tplID := this.cfg.TplID
	if len(tpl) > 0 {
		tplID = tpl[0]
	}
	ss := qcloudsms.SMSSingleReq{
		Tel: qcloudsms.SMSTel{
			Nationcode: "86",
			Mobile:     phoneNumber,
		},
		Type:   0,
		TplID:  tplID,
		Params: params,
	}
	ok, err := this.client.SendSMSSingle(ss)
	if ok {
		return nil
	} else {
		return err
	}
}

//发送短信
func (this *QCloudSMS) SendSms(phoneNumber string, content string, tpl ...int) error {
	TplID := this.cfg.TplID
	if len(tpl) > 0 {
		TplID = tpl[0]
	}
	ss := qcloudsms.SMSSingleReq{
		Tel: qcloudsms.SMSTel{
			Nationcode: "86",
			Mobile:     phoneNumber,
		},
		Type:  0,
		Msg:   content,
		TplID: TplID,
	}
	ok, err := this.client.SendSMSSingle(ss)
	if ok {
		return nil
	} else {
		return err
	}
}

//发送语音通知
func (this *QCloudSMS) SendVoiceWithTpl(phoneNumber string, params []string, tpl ...int) error {
	tplID := this.cfg.TplID
	if len(tpl) > 0 {
		tplID = tpl[0]
	}
	tp, err := this.GetTemplate(uint(tplID))
	if err != nil {
		return err
	}
	content := tp.Text
	for i := 0; i < len(params); i++ {
		content = strings.Replace(content, fmt.Sprintf("{%d}", i+1), params[i], 1)
	}
	v := qcloudsms.VoiceReq{
		Tel: qcloudsms.SMSTel{
			Nationcode: "86",
			Mobile:     phoneNumber,
		},
		Promptfile: content,
		Prompttype: tp.Type,
	}
	ok, err := this.client.SendVoice(v)
	if ok {
		return nil
	} else {
		return err
	}
}

//发送语音通知
func (this *QCloudSMS) SendVoice(phoneNumber string, content string) error {
	v := qcloudsms.VoiceReq{
		Tel: qcloudsms.SMSTel{
			Nationcode: "86",
			Mobile:     phoneNumber,
		},
		Promptfile: content,
		Prompttype: 2,
	}
	ok, err := this.client.SendVoice(v)
	if ok {
		return nil
	} else {
		return err
	}
}

//发送通知
func (this *QCloudSMS) Send(phoneNumber string, content string, way ...string) error {
	w := this.cfg.Way
	if len(way) > 0 {
		w = way[0]
	}
	if w == "voice" {
		if this.cfg.Receivers != nil && len(this.cfg.Receivers) > 0 {
			for _, r := range this.cfg.Receivers {
				this.SendVoice(r, content)
			}
		}
		return this.SendVoice(phoneNumber, content)
	} else {
		if this.cfg.Receivers != nil && len(this.cfg.Receivers) > 0 {
			for _, r := range this.cfg.Receivers {
				this.SendSms(r, content)
			}
		}
		return this.SendSms(phoneNumber, content)
	}
}

//发送通知
func (this *QCloudSMS) SendWithTpl(phoneNumber string, params []string, args ...interface{}) error {
	way := this.cfg.Way
	tpl := this.cfg.TplID
	for _, v := range args {
		switch v.(type) {
		case int:
			tpl = v.(int)
			break
		case string:
			way = v.(string)
			break
		}
	}
	if way == "voice" {
		if this.cfg.Receivers != nil && len(this.cfg.Receivers) > 0 {
			for _, r := range this.cfg.Receivers {
				this.SendVoiceWithTpl(r, params,tpl)
			}
		}
		return this.SendVoiceWithTpl(phoneNumber, params, tpl)
	} else {
		if this.cfg.Receivers != nil && len(this.cfg.Receivers) > 0 {
			for _, r := range this.cfg.Receivers {
				this.SendSmsWithTpl(r, params,tpl)
			}
		}
		return this.SendSmsWithTpl(phoneNumber, params, tpl)
	}
}

func (this *QCloudSMS) GetTemplate(id uint) (*qcloudsms.Template, error) {
	rs, err := this.client.GetTemplateByID([]uint{id})
	if err != nil {
		return nil, err
	}
	if rs.Count > 0 {
		return &(rs.Data[0]), nil
	} else {
		return nil, fmt.Errorf("template %d not found", id)
	}
}

func NewQCloudSMS(cfg config.QCloudSms) *QCloudSMS {
	qc := &QCloudSMS{cfg: cfg}
	opt := qcloudsms.NewOptions()
	opt.UserAgent = "GameWinner Inc"
	opt.Debug = cfg.Debug
	opt.APPID = cfg.APPID
	opt.APPKEY = cfg.APPKEY
	opt.SIGN = cfg.SIGN
	qc.client = qcloudsms.NewClient(opt)
	return qc
}
