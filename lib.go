package robot

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"unsafe"
)

func Recovery(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
		if err != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("error: %w;\ntrace: %v", err, Bytes2String(buf))
		}
	}()
	err = f()
	return
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func String2bytes(s string) []byte {
	tmp := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{tmp[0], tmp[1], tmp[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func JsonEncode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func JsonDecode(stack []byte, v interface{}) error {
	return json.Unmarshal(stack, &v)
}

func JsonEncodeString(v interface{}) (string, error) {
	b, err := JsonEncode(v)
	if err != nil {
		return "", err
	}
	return Bytes2String(b), nil
}

func MustLoadJson(file string, val interface{}) {
	buf, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("read file failed: %v", err)
	}
	err = JsonDecode(buf, val)
	if err != nil {
		log.Fatalf("load json failed: %v, %v", err, Bytes2String(buf))
	}
}
