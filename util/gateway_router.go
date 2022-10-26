package util

import "zgoframe/protobuf/pb"

func (netWay *NetWay) Router(msg pb.Msg, conn *Conn) (data interface{}, err error) {
	return netWay.Option.RouterBack(msg, conn)
}
