package client

import (
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/gamelee/robot/app/msg"
	"github.com/gamelee/robot/connect"
)

type Client struct {
	conf *Config
	conn *connect.Conner // 链接
	seq  uint64          // 发包序列
	stop int32
}

func (cli *Client) NewMessage(id, data interface{}) *connect.Message {
	mid, ok := id.(msg.ID)
	if !ok {
		panic("消息id 只能是 msg.ID 类型")
	}

	req, ok := data.(*msg.Req)
	if !ok {
		panic("只能发送 *msg.Req 类型的消息")
	}

	message := connect.NewMessage(cli.Name(), id, data)

	req.Seq = atomic.AddUint64(&cli.seq, 1)
	message.WaitId = req.Seq
	message.From = cli.Name()
	message.Type = connect.ReqMessage
	message.IDPretty = mid.String()
	return message
}

func (cli *Client) Running() bool {
	return atomic.LoadInt32(&cli.stop) == 0
}

func (cli *Client) Name() string {
	return cli.conf.Name
}

func (cli *Client) Close() {
	if cli.Running() {
		if cli.conn != nil && cli.Running() {
			cli.conn.Close()
		}
		atomic.AddInt32(&cli.stop, 1)
	}
}

func (cli *Client) KeepAlive() (time.Duration, *connect.Message) {
	message := cli.NewMessage(msg.ID_ID_HEART, &msg.Req{Heart: &msg.ReqHeart{Time: time.Now().Unix()}})
	return time.Second * 14, message
}

// New
// @Desc 创建服务实例
// @Date 10:36 2020/9/25
func New(conf *Config) connect.Connector {
	this := &Client{
		conf: conf,
		stop: 0,
	}
	return this
}

// Connect 连接
func (cli *Client) Connect() (err error) {
	defer func() {
		if err != nil {
			cli.Close()
		}
	}()
	cli.conn, err = connect.Dial("tcp", cli.conf.Addr)
	return err
}

func (cli *Client) Write(message *connect.Message) error {
	id := message.ID.(msg.ID)
	data := message.Data.(proto.Message)
	return cli.write(id, data)
}

// write 写入一条消息
func (cli *Client) write(id msg.ID, req proto.Message) error {
	bodyData, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	head := new(msg.Head)
	head.ID = id
	head.BodyLen = uint32(len(bodyData))
	headData, err := proto.Marshal(head)
	if err != nil {
		return err
	}

	m := make([]byte, 1+len(headData)+len(bodyData))
	m[0] = byte(len(headData))
	copy(m[1:], headData)
	copy(m[1+len(headData):], bodyData)
	// 发送消息
	_, err = cli.conn.Write(m)
	if err != nil {
		cli.Close()
		return err
	}
	err = cli.conn.Flush()
	if err != nil {
		return err
	}
	return nil
}

// Read 读取一条消息
func (cli *Client) Read() *connect.Message {
	message := cli.read()
	message.From = cli.Name()
	return message
}

// read 读取一条消息
func (cli *Client) read() *connect.Message {
	message := new(connect.Message)
	bufHeadLen := make([]byte, 1)
	if _, message.Error = cli.conn.Read(bufHeadLen); message.Error != nil {
		return message
	}
	// read head
	bufHead := make([]byte, uint32(bufHeadLen[0]))
	if _, message.Error = cli.conn.Read(bufHead); message.Error != nil {
		return message
	}
	head := &msg.Head{}
	if message.Error = proto.Unmarshal(bufHead, head); message.Error != nil {
		return message
	}
	message.ID = head.ID
	if head.BodyLen > 0 {
		// read body
		bufMsg := make([]byte, head.BodyLen)
		if _, message.Error = cli.conn.Read(bufMsg); message.Error != nil {
			return message
		}

		message.ID = head.ID
		if message.Error = proto.Unmarshal(bufMsg, message.Data.(proto.Message)); message.Error != nil {
			return message
		}
	}
	if message.Data == nil {
		message.Data = &msg.Rsp{}
	}
	message.Type = connect.RspMessage

	message.IDPretty = head.ID.String()
	return message
}
