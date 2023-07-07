package connect

import (
	"bufio"
	"net"
)

type Conner struct {
	net.Conn
	br *bufio.Reader
	bw *bufio.Writer
}

func Dial(network, address string) (*Conner, error) {
	netConn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	con := &Conner{
		Conn: netConn,
		br:   bufio.NewReader(netConn),
		bw:   bufio.NewWriter(netConn),
	}
	return con, nil
}

func (c *Conner) Read(b []byte) (n int, err error) {
	byteLen := len(b)
	readLen := 0
	for byteLen > readLen { // 循环一直读完
		n, err = c.br.Read(b[readLen:])
		readLen += n
		if err != nil {
			return readLen, err
		}
	}
	return readLen, nil
}

func (c *Conner) Write(b []byte) (n int, err error) {
	return c.bw.Write(b)
}

func (c *Conner) Flush() error {
	if err := c.bw.Flush(); err != nil {
		return c.Fatal(err)
	}
	return nil
}

func (c *Conner) Fatal(err error) error {
	_ = c.Conn.Close()
	return err
}

func (c *Conner) Close() {
	_ = c.Conn.Close()
}
