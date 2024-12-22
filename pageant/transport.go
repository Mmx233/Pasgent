package pageant

import (
	"encoding/binary"
	"fmt"
	"github.com/Mmx233/Pasgent/system"
	"github.com/Mmx233/Pasgent/tools"
	"golang.org/x/sys/windows"
	"io"
)

const (
	agentMaxMsglen = 8192
)

type FileMappingConn struct {
	sharedFile  windows.Handle
	sharedMem   uintptr
	request     []byte
	readOffset  int
	writeOffset int
}

func NewFileMapping(mapName string) (*FileMappingConn, error) {
	sharedFile, err := system.OpenFileMapping(
		windows.FILE_MAP_WRITE|windows.FILE_MAP_READ,
		0,
		mapName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create shared file: %s", err)
	}

	sharedMem, err := windows.MapViewOfFile(
		sharedFile,
		windows.FILE_MAP_WRITE|windows.FILE_MAP_READ,
		0,
		0,
		agentMaxMsglen,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to map file into shared memory: %s", err)
	}

	messageSize := int(binary.BigEndian.Uint32(tools.PtrToSlice(sharedMem, 4)))
	request := make([]byte, messageSize+4)
	copy(request, tools.PtrToSlice(sharedMem, messageSize+4))

	return &FileMappingConn{
		sharedFile: sharedFile,
		sharedMem:  sharedMem,
		request:    request,
	}, nil
}

// Close frees resources used by Conn.
func (c *FileMappingConn) Close() error {
	if c.sharedMem == 0 {
		return nil
	}
	errUnmap := windows.UnmapViewOfFile(c.sharedMem)
	errClose := windows.CloseHandle(c.sharedFile)
	if errUnmap != nil {
		return errUnmap
	} else if errClose != nil {
		return errClose
	}
	c.sharedMem = 0
	c.sharedFile = windows.InvalidHandle
	return nil
}

func (c *FileMappingConn) Read(p []byte) (n int, err error) {
	if c.readOffset >= len(c.request) {
		return 0, io.EOF
	}
	src := c.request[c.readOffset : c.readOffset+len(p)]
	copy(p, src)
	c.readOffset += len(p)
	return len(src), nil
}

func (c *FileMappingConn) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, fmt.Errorf("message to send is empty")
	}
	dst := tools.PtrToSlice(c.sharedMem+uintptr(c.writeOffset), len(p))
	copy(dst, p)
	c.writeOffset += len(p)
	return len(p), nil
}
