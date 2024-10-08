package im

import uuid "github.com/satori/go.uuid"

type Msg struct {
	SessionId string `json:"session_id"`
	Content   string `json:"content"`
	SendUid   int    `json:"send_uid"`
	RecvUid   int    `json:"recv_uid"`
	GroupId   int    `json:"group_id"`
}

type Im struct {
}

func New() *Im {
	im := new(Im)
	return im
}

func (im *Im) Init() {

}

func (im *Im) Send(msg Msg) {
	if msg.SendUid == 0 || msg.RecvUid == 0 {

	}

	if msg.SessionId == "" {
		msg.SessionId = im.MakeSessionId()
	}

	if msg.GroupId == 0 {

	}
}

func (im *Im) MakeSessionId() string {
	return uuid.NewV4().String()
}

func (im *Im) Consumer() {

}
