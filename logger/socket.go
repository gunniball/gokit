package logger

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"time"
)

// this is draft...
func Send(d interface{}, conn net.Conn) error {
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
		conn.SetWriteDeadline(now.Add(time.Second * 5)) // TODO: config
		n, err = conn.Write(send_buff[total:])
		total += n
		if err != nil {
			log.Printf("error while sending 4 bytes: %s", err)
			return err
		}
		if total == 4 {
			break
		}
	}
	total = 0
	for {
		conn.SetWriteDeadline(now.Add(time.Second * 5)) // TODO: config
		n, err = conn.Write(data[total:])
		total += n
		if err != nil {
			log.Printf("error while sending request: %s", err)
			return err
		}
		if total == size {
			break
		}
	}
	return nil
}

func Read(conn net.Conn) []byte {
	var MAXSIZE = 99999
	var MAXATTEMPTS = 10
	var err error
	var tmp_buffer [4]byte
	var buffer []byte
	var numbytes, tmp_numbytes, size, attempts int
	for {
		if attempts >= MAXATTEMPTS {
			log.Println("Got maximum attempts count, closing connection")
			return nil
		}
		numbytes = 0
		for {
			tmp_numbytes, err = conn.Read(tmp_buffer[numbytes:])
			if err != nil {
				if err != io.EOF {
					log.Printf("4 bytes read error: %s\n", err.Error())
				}
				return nil
			}
			numbytes += tmp_numbytes
			if numbytes < 4 {
				attempts++
				continue
			}
			break
		}

		size = 0
		for i := 0; i < 4; i++ {
			size = size*256 + int(tmp_buffer[i])
		}
		if size == 0 {
			attempts++
			continue
		}
		if size > MAXSIZE {
			attempts++
			continue
		}

		buffer = make([]byte, size)
		numbytes = 0
		for {
			tmp_numbytes, err = conn.Read(buffer[numbytes:])
			if err != nil {
				if err != io.EOF {
					log.Printf("read in client buffer error: %s\n", err.Error())
				}
				return nil
			}
			numbytes += tmp_numbytes
			if numbytes < size {
				attempts++
				continue
			}
			break
		}
		return buffer
	}
}
