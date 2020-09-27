package utils

var NopIO = NIO{}

type NIO struct{}

func (n2 NIO) Read(p []byte) (n int, err error) {
	return len(p), nil
}

func (n2 NIO) Write(p []byte) (n int, err error) {
	return len(p), nil
}
