package logger

import (
	"encoding/json"
	"time"
	"errors"
)

var notOpenError = errors.New("no")

func (s *SocketLogger) Send(d interface{}) error {
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	size := len(data)
	send_buff := make([]byte, 0, 4)
	send_buff = append(send_buff, byte(size>>24))
	send_buff = append(send_buff, byte(size>>16))
	send_buff = append(send_buff, byte(size>>8))
	send_buff = append(send_buff, byte(size))
	var n int
	var total int
	var now = time.Now()
	for {
		if s.Conn == nil {
			return notOpenError
		}
		s.Conn.SetWriteDeadline(now.Add(time.Second * 5)) // TODO: config
		n, err = s.Conn.Write(send_buff[total:])
		total += n
		if err != nil {
			return err
		}
		if total == 4 {
			break
		}
	}
	total = 0
	for {
		s.Conn.SetWriteDeadline(now.Add(time.Second * 5)) // TODO: config
		n, err = s.Conn.Write(data[total:])
		total += n
		if err != nil {
			return err
		}
		if total == size {
			break
		}
	}
	return nil
}

