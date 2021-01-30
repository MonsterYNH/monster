package utils

// MultiServe 用于开启多个监听端口服务的类型
type MultiServe func() error

// MultiRun 以分组的方式并行运行多个程序，发生不可恢复的错误时停止整个组的GOROUTING
func MultiRun(serves ...MultiServe) error {
	c := make(chan error)
	for _, serve := range serves {
		go func(s MultiServe) {
			if err := s(); err != nil {
				c <- err
			}
		}(serve)
	}
	return <-c
}
