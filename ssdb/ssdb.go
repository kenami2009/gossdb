package ssdb

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type Error string

func (err Error) Error() string { return string(err) }

var okReply interface{} = "OK"
var nilReply interface{} = nil

var connErr = errors.New("can not connect ssdb")

var ssdbErr = errors.New("ssdb error.")

//字节切片转换为整数
func parseInt(p []byte) (interface{}, error) {
	if len(p) == 0 {
		return 0, errors.New("malformed integer")
	}

	var negate bool
	if p[0] == '-' {
		negate = true
		p = p[1:]
		if len(p) == 0 {
			return 0, errors.New("malformed integer")
		}
	}

	var n int64
	for _, b := range p {
		n *= 10
		if b < '0' || b > '9' {
			return 0, errors.New("illegal bytes in length")
		}
		n += int64(b - '0')
	}

	if negate {
		n = -n
	}
	return n, nil
}

type Conn struct {
	conn net.Conn
	bw   *bufio.Writer
	br   *bufio.Reader
}

func Connect(addr string) (conn *Conn, err error) {
	netConn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(fmt.Errorf("%s:%s", connErr.Error(), addr))
	}
	conn = &Conn{
		conn: netConn,
		bw:   bufio.NewWriter(netConn),
		br:   bufio.NewReader(netConn),
	}
	return
}

func (c *Conn) Close() {
	c.conn.Close()
}

func (c *Conn) writeCRLF() {
	c.bw.WriteString("\r\n")
}

func (c *Conn) writeInt64(n int64) {
	s := make([]byte, 0)
	c.writeBytes(strconv.AppendInt(s, n, 10))
}
func (c *Conn) writeLen(prefix byte, n int) {
	c.bw.Write([]byte{prefix, byte('0' + n%10)})
	c.writeCRLF()
}
func (c *Conn) writeBytes(p []byte) {
	c.writeLen('$', len(p))
	c.bw.Write(p)
	c.writeCRLF()
}
func (c *Conn) writeString(s string) {
	c.writeLen('$', len(s))
	c.bw.WriteString(s)
	c.writeCRLF()
}

func (c *Conn) writeArgs(args ...interface{}) {
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			c.writeString(arg)
		case []byte:
			c.writeBytes(arg)
		case int:
			c.writeInt64(int64(arg))
		}
	}
}

func (c *Conn) Flush() error {
	return c.bw.Flush()
}

func (c *Conn) Do(cmd string, args ...interface{}) (interface{}, error) {
	c.writeLen('*', len(args)+1)
	c.writeString(cmd)
	c.writeArgs(args...)

	if err := c.Flush(); err != nil {
		return nil, err
	}
	res, err := c.ReadLine()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Conn) ReadLine() (interface{}, error) {
	line, err := c.br.ReadSlice('\n')
	if err != nil {
		return nil, err
	}

	fmtLine := make([]byte, len(line)-2)

	copy(fmtLine, line[:len(line)-2])
	switch fmtLine[0] {
	//字符串
	case '+':
		if len(fmtLine) == 3 && fmtLine[1] == 'O' {
			return okReply, nil
		}
	//错误
	case '-':
		return nil, errors.New(string(fmtLine[1:]))
	//整数
	case ':':
		if fmtLine[1] == '0' {
			return okReply, nil
			//
		} else {
			return parseInt(fmtLine[1:])
		}
	//批量字符串
	case '$':
		//$-1 (nil)
		if len(line) == 3 && line[1] == '-' && line[2] == '1' {
			return nilReply, nil
		}
		if l, err := c.br.ReadSlice('\n'); err == nil {
			return l[:len(l)-2], nil
		}
	//数组
	case '*':
	default:

	}
	return nil, nil
}
