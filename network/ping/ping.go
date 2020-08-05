package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	count = flag.Int64("c", 0, "stop after sending (and receiving) count ECHO_RESPONSE packets.")
	dest  = flag.String("d", "", "Specify the sending destination")
)

type PingOptions struct {
	Count       int64
	Destination string
}

var (
	options *PingOptions
)

func init() {
	flag.Parse()
	if *dest == "" {
		fmt.Println("invalid destination")
		os.Exit(1)
	}

	options = &PingOptions{
		Count:       *count,
		Destination: *dest,
	}
}

func ping(seq int16) {
	msg := generateICMPMsg(seq)

	start := time.Now()
	conn, err := net.Dial("ip4:icmp", options.Destination)
	checkErr(err)

	_, err = conn.Write(msg)
	checkErr(err)

	receive := make([]byte, 76)
	_, err = conn.Read(receive)
	checkErr(err)

	// 判断请求和应答的ID标识符, sequence序列码是否一致, 以及ICMP是否超时
	if receive[24] != msg[4] || receive[25] != msg[5] || receive[26] != msg[6] || receive[27] != msg[7] || receive[20] == 11 {
		fmt.Printf("Request timeout for icmp_seq %d\n", seq)
		return
	}

	fmt.Printf("%d bytes from %s: icmp_seq=%d ttl=%d time=%.3f ms\n", len(receive), conn.RemoteAddr().String(), seq, receive[8], float64(time.Since(start).Microseconds())/1000)
}

func main() {
	cName, err := net.LookupCNAME(options.Destination)
	checkErr(err)

	// 这里建立链接只是为了打印链接信息, 不会使用该链接发送ping包
	conn, err := net.Dial("ip4:icmp", options.Destination)
	checkErr(err)

	fmt.Printf("PING %s (%s): 56 data bytes\n", cName, conn.RemoteAddr().String())

	var seq int16
	for {
		if options.Count != 0 && seq >= int16(options.Count) {
			break
		}

		ping(seq)
		seq++

		time.Sleep(time.Second)
	}
}

func generateICMPMsg(seq int16) []byte {
	// ICMP头8字节, 数据部分48字节
	msg := make([]byte, 56)
	// 8 表示回显请求
	// ping的请求和应答, 该code都为0
	msg[0], msg[1] = 8, 0
	// 校验码占2个字节
	msg[2], msg[3] = 0, 0
	// ID标示符占2个字节
	msg[4], msg[5] = genIdentifier()
	// 序列号占2个字段
	msg[6], msg[7] = genSequence(seq)

	// 计算校验和
	checkResult := checkSum(msg)
	msg[2], msg[3] = byte(checkResult>>8), byte(checkResult&255)

	return msg
}

func genIdentifier() (byte, byte) {
	return options.Destination[0], options.Destination[1]
}

func genSequence(v int16) (byte, byte) {
	return byte(v >> 8), byte(v & 255)
}

func checkSum(msg []byte) uint16 {
	sum := 0

	for i := 0; i < len(msg)-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if len(msg)%2 == 1 {
		sum += int(msg[len(msg)-1]) * 256
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	return uint16(^sum)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
