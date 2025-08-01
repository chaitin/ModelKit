package tee

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"
)

type ReqTee struct {
	src     io.ReadCloser
	buf     bytes.Buffer
	done    chan struct{}
	mu      sync.RWMutex
	closed  bool
	maxSize int64 // 最大缓冲区大小，0表示无限制
}

var _ io.ReadCloser = &ReqTee{}

func NewReqTee(src io.ReadCloser) *ReqTee {
	return NewReqTeeWithMaxSize(src, 0)
}

// NewReqTeeWithMaxSize 创建带最大缓冲区大小限制的ReqTee
func NewReqTeeWithMaxSize(src io.ReadCloser, maxSize int64) *ReqTee {
	return &ReqTee{
		src:     src,
		buf:     bytes.Buffer{},
		done:    make(chan struct{}),
		maxSize: maxSize,
	}
}

func (r *ReqTee) GetBody() ([]byte, error) {
	return r.GetBodyWithTimeout(30 * time.Second)
}

// GetBodyWithTimeout 获取缓冲的数据，支持自定义超时时间
func (r *ReqTee) GetBodyWithTimeout(timeout time.Duration) ([]byte, error) {
	select {
	case <-r.done:
		r.mu.RLock()
		defer r.mu.RUnlock()
		return r.buf.Bytes(), nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout waiting for data")
	}
}

// Close implements io.ReadCloser.
func (r *ReqTee) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return nil
	}
	r.closed = true

	close(r.done)
	return r.src.Close()
}

// Read implements io.ReadCloser.
func (r *ReqTee) Read(p []byte) (n int, err error) {
	n, err = r.src.Read(p)
	if n > 0 {
		r.mu.Lock()
		// 检查缓冲区大小限制
		if r.maxSize > 0 && int64(r.buf.Len()+n) > r.maxSize {
			r.mu.Unlock()
			return n, fmt.Errorf("buffer size limit exceeded: %d bytes", r.maxSize)
		}
		// 直接写入缓冲区，避免额外的内存分配和复制
		r.buf.Write(p[:n])
		r.mu.Unlock()
	}
	return n, err
}
