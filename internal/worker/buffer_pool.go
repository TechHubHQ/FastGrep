package worker

import (
	"fastgrep/internal/config"
	"sync"
)

var BufferPool = sync.Pool{
	New: func() any {
		return make([]byte, config.DefaultChunkSize)
	},
}

// GetBuffer retrieves a buffer from the pool. It always returns a slice of DefaultChunkSize.
func GetBuffer() []byte {
	buf := BufferPool.Get().([]byte)
	// Ensure the slice we return is exactly the default chunk size
	if len(buf) != config.DefaultChunkSize {
		return make([]byte, config.DefaultChunkSize)
	}
	return buf
}

// PutBuffer returns a buffer to the pool only if it's the correct size.
func PutBuffer(buf []byte) {
	if cap(buf) == config.DefaultChunkSize {
		// Reset to full capacity before returning
		BufferPool.Put(buf[:config.DefaultChunkSize])
	}
}
