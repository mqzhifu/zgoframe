package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"math/rand"
	"sort"
	"time"
)

type AccessToken struct {
	AppID          string
	AppCertificate string
	ChannelName    string
	UidStr         string
	Ts             uint32
	Salt           uint32
	Message        map[uint16]uint32
	Signature      string
	CrcChannelName uint32
	CrcUid         uint32
	MsgRawContent  string
}

func random(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func panichandler() {
	if r := recover(); r != nil {
		fmt.Println("error: ", r)
	}
}

func getVersion() string {
	return "006"
}

func CreateAccessToken(appID, appCertificate, channelName string, uid uint32) AccessToken {
	var uidStr string
	if uid == 0 {
		uidStr = ""
	} else {
		uidStr = fmt.Sprintf("%d", uid)
	}
	ts := uint32(time.Now().Unix()) + 24*3600
	salt := uint32(random(1, 99999999))
	message := make(map[uint16]uint32)
	return AccessToken{appID, appCertificate, channelName, uidStr, ts, salt, message, "", 0, 0, ""}
}

func CreateAccessToken2(appID, appCertificate, channelName string, uid string) AccessToken {
	ts := uint32(time.Now().Unix()) + 24*3600
	salt := uint32(random(1, 99999999))
	message := make(map[uint16]uint32)
	return AccessToken{appID, appCertificate, channelName, uid, ts, salt, message, "", 0, 0, ""}
}

func (token *AccessToken) FromString(originToken string) bool {
	defer panichandler()

	dk6version := getVersion()
	originVersion := originToken[:VERSION_LENGTH]
	if originVersion != dk6version {
		return false
	}

	//originAppID := originToken[VERSION_LENGTH:(VERSION_LENGTH + APP_ID_LENGTH)]
	originContent := originToken[(VERSION_LENGTH + APP_ID_LENGTH):]
	originContentDecoded, err := base64.StdEncoding.DecodeString(originContent)
	if err != nil {
		return false
	}

	signature_, crc_channel_name_, crc_uid_, msg_raw_content_, err := unPackContent(originContentDecoded)
	if err != nil {
		return false
	}
	token.Signature = signature_
	token.CrcChannelName = crc_channel_name_
	token.CrcUid = crc_uid_
	token.MsgRawContent = msg_raw_content_

	salt_, ts_, messages_, err := unPackMessages(token.MsgRawContent)
	if err != nil {
		return false
	}
	token.Salt = salt_
	token.Ts = ts_
	token.Message = messages_

	return true
}

func (token *AccessToken) AddPrivilege(privilege Privileges, expireTimestamp uint32) {
	pri := uint16(privilege)
	token.Message[pri] = expireTimestamp
}

func (token *AccessToken) Build() (string, error) {
	ret := ""
	version := getVersion()

	buf_m := new(bytes.Buffer)
	if err := packUint32(buf_m, token.Salt); err != nil {
		return ret, err
	}
	if err := packUint32(buf_m, token.Ts); err != nil {
		return ret, err
	}
	if err := packMapUint32(buf_m, token.Message); err != nil {
		return ret, err
	}
	bytes_m := buf_m.Bytes()

	buf_val := new(bytes.Buffer)
	val := token.AppID + token.ChannelName + token.UidStr
	buf_val.Write([]byte(val))
	buf_val.Write(bytes_m)
	bytes_val := buf_val.Bytes()

	buf_sig := hmac.New(sha256.New, []byte(token.AppCertificate))
	buf_sig.Write(bytes_val)
	bytes_sig := buf_sig.Sum(nil)

	crc32q := crc32.MakeTable(0xedb88320)
	crc_channel_name := crc32.Checksum([]byte(token.ChannelName), crc32q)
	crc_uid := crc32.Checksum([]byte(token.UidStr), crc32q)

	buf_content := new(bytes.Buffer)
	if err := packString(buf_content, string(bytes_sig[:])); err != nil {
		return ret, err
	}
	if err := packUint32(buf_content, crc_channel_name); err != nil {
		return ret, err
	}
	if err := packUint32(buf_content, crc_uid); err != nil {
		return ret, err
	}
	if err := packString(buf_content, string(bytes_m[:])); err != nil {
		return ret, err
	}
	bytes_content := buf_content.Bytes()

	ret = version + token.AppID + base64.StdEncoding.EncodeToString(bytes_content)
	return ret, nil
}

func packUint16(w io.Writer, n uint16) error {
	return binary.Write(w, binary.LittleEndian, n)
}

func packUint32(w io.Writer, n uint32) error {
	return binary.Write(w, binary.LittleEndian, n)
}

func packString(w io.Writer, s string) error {
	err := packUint16(w, uint16(len(s)))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(s))
	return err
}

func packHexString(w io.Writer, s string) error {
	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	return packString(w, string(b))
}

func packExtra(w io.Writer, extra map[uint16]string) error {
	keys := []int{}
	if err := packUint16(w, uint16(len(extra))); err != nil {
		return err
	}
	for k := range extra {
		keys = append(keys, int(k))
	}
	//should sorted keys
	sort.Ints(keys)

	for _, k := range keys {
		v := extra[uint16(k)]
		if err := packUint16(w, uint16(k)); err != nil {
			return err
		}
		if err := packString(w, v); err != nil {
			return err
		}
	}
	return nil
}

func packMapUint32(w io.Writer, extra map[uint16]uint32) error {
	keys := []int{}
	if err := packUint16(w, uint16(len(extra))); err != nil {
		return err
	}
	for k := range extra {
		keys = append(keys, int(k))
	}
	//should sorted keys
	sort.Ints(keys)

	for _, k := range keys {
		v := extra[uint16(k)]
		if err := packUint16(w, uint16(k)); err != nil {
			return err
		}
		if err := packUint32(w, v); err != nil {
			return err
		}
	}
	return nil
}

func unPackUint16(r io.Reader) (uint16, error) {
	var n uint16
	err := binary.Read(r, binary.LittleEndian, &n)
	return n, err
}

func unPackUint32(r io.Reader) (uint32, error) {
	var n uint32
	err := binary.Read(r, binary.LittleEndian, &n)
	return n, err
}

func unPackString(r io.Reader) (string, error) {
	n, err := unPackUint16(r)
	if err != nil {
		return "", err
	}

	buf := make([]byte, n)
	r.Read(buf)
	s := string(buf[:])
	return s, err
}

func unPackContent(buff []byte) (string, uint32, uint32, string, error) {
	in := bytes.NewReader(buff)
	sig, err := unPackString(in)
	if err != nil {
		return "", 0, 0, "", err
	}

	crc_channel_name, err := unPackUint32(in)
	if err != nil {
		return "", 0, 0, "", err
	}
	crc_uid, err := unPackUint32(in)
	if err != nil {
		return "", 0, 0, "", err
	}
	m, err := unPackString(in)
	if err != nil {
		return "", 0, 0, "", err
	}

	return sig, crc_channel_name, crc_uid, m, nil
}

func unPackMessages(msg_str string) (uint32, uint32, map[uint16]uint32, error) {
	msg_map := make(map[uint16]uint32)

	msg_byte := []byte(msg_str)
	in := bytes.NewReader(msg_byte)

	salt, err := unPackUint32(in)
	if err != nil {
		return 0, 0, msg_map, err
	}
	ts, err := unPackUint32(in)
	if err != nil {
		return 0, 0, msg_map, err
	}

	len, err := unPackUint16(in)
	if err != nil {
		return 0, 0, msg_map, err
	}
	for i := uint16(0); i < len; i++ {
		key, err := unPackUint16(in)
		if err != nil {
			return 0, 0, msg_map, err
		}
		value, err := unPackUint32(in)
		if err != nil {
			return 0, 0, msg_map, err
		}
		msg_map[key] = value
	}

	return salt, ts, msg_map, nil
}

//=================================================================================================================

//RtmTokenBuilder class
type RtmTokenBuilder struct {
}

//BuildToken method
// appID: The App ID issued to you by Agora. Apply for a new App ID from
//        Agora Dashboard if it is missing from your kit. See Get an App ID.
// appCertificate:	Certificate of the application that you registered in
//                  the Agora Dashboard. See Get an App Certificate.
// userAccount: The user account.
// role: Role_Rtm_User = 1
// privilegeExpireTs: represented by the number of seconds elapsed since
//                    1/1/1970. If, for example, you want to access the
//                    Agora Service within 10 minutes after the token is
//                    generated, set expireTimestamp as the current
//                    timestamp + 600 (seconds)./
func RTMBuildToken(appID string, appCertificate string, userAccount string, role RTMRole, privilegeExpiredTs uint32) (string, error) {
	token := CreateAccessToken2(appID, appCertificate, userAccount, "")
	token.AddPrivilege(KLoginRtm, privilegeExpiredTs)
	return token.Build()
}

//=====================================

//RtcTokenBuilder class
type RtcTokenBuilder struct {
}

//BuildTokenWithUserAccount method
// appID: The App ID issued to you by Agora. Apply for a new App ID from
//        Agora Dashboard if it is missing from your kit. See Get an App ID.
// appCertificate:	Certificate of the application that you registered in
//                  the Agora Dashboard. See Get an App Certificate.
// channelName:Unique channel name for the AgoraRTC session in the string format
// uid: User ID. A 32-bit unsigned integer with a value ranging from
//      1 to (232-1). optionalUid must be unique.
// role: Role_Publisher = 1: A broadcaster (host) in a live-broadcast profile.
//       Role_Subscriber = 2: (Default) A audience in a live-broadcast profile.
// privilegeExpireTs: represented by the number of seconds elapsed since
//                    1/1/1970. If, for example, you want to access the
//                    Agora Service within 10 minutes after the token is
//                    generated, set expireTimestamp as the current
//                    timestamp + 600 (seconds)./
func RTCBuildTokenWithUserAccount(appID string, appCertificate string, channelName string, userAccount string, role RTCRole, privilegeExpiredTs uint32) (string, error) {
	token := CreateAccessToken2(appID, appCertificate, channelName, userAccount)
	token.AddPrivilege(KJoinChannel, privilegeExpiredTs)

	if (role == RoleAttendee) || (role == RolePublisher) || (role == RoleAdmin) {
		token.AddPrivilege(KPublishVideoStream, privilegeExpiredTs)
		token.AddPrivilege(KPublishAudioStream, privilegeExpiredTs)
		token.AddPrivilege(KPublishDataStream, privilegeExpiredTs)
	}
	return token.Build()
}

//BuildTokenWithUID method
// appID: The App ID issued to you by Agora. Apply for a new App ID from
//        Agora Dashboard if it is missing from your kit. See Get an App ID.
// appCertificate:	Certificate of the application that you registered in
//                  the Agora Dashboard. See Get an App Certificate.
// channelName:Unique channel name for the AgoraRTC session in the string format
// userAccount: The user account.
// role: Role_Publisher = 1: A broadcaster (host) in a live-broadcast profile.
//       Role_Subscriber = 2: (Default) A audience in a live-broadcast profile.
// privilegeExpireTs: represented by the number of seconds elapsed since
//                    1/1/1970. If, for example, you want to access the
//                    Agora Service within 10 minutes after the token is
//                    generated, set expireTimestamp as the current
func BuildTokenWithUID(appID string, appCertificate string, channelName string, uid uint32, role RTCRole, privilegeExpiredTs uint32) (string, error) {
	uidStr := fmt.Sprint(uid)
	if uid == 0 {
		uidStr = ""
	}
	return RTCBuildTokenWithUserAccount(appID, appCertificate, channelName, uidStr, role, privilegeExpiredTs)
}

//==========================================================
