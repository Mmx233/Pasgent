package pageant

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Microsoft/go-winio"
	"github.com/Mmx233/Pasgent/system"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"time"
)

func AgentRequestHandler(data []byte) error {
	mapName := string(bytes.TrimRight(data, "\000"))
	shared, err := NewFileMapping(mapName)
	if err != nil {
		return fmt.Errorf("open file mapping failed: %v", err)
	}
	defer shared.Close()

	pipePath := `\\.\pipe\openssh-ssh-agent`
	pipeTimeout := time.Second * 15
	pipe, err := winio.DialPipe(pipePath, &pipeTimeout)
	if err != nil {
		return fmt.Errorf("open 1password ssh agent named pipe failed: %v", err)
	}
	defer pipe.Close()

	err = agent.ServeAgent(agent.NewClient(pipe), shared)
	if err != nil && errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func Run(logger logrus.FieldLogger) error {
	if err := system.CreateHiddenMessageWindow(
		"Pageant", "Pageant",
		system.NewAgentWndProc(func(data []byte) error {
			err := AgentRequestHandler(data)
			if err != nil {
				logger.Warnln("handle agent request failed:", err)
			} else {
				logger.Infoln("request succeeded:", string(data))
			}
			return err
		}),
	); err != nil {
		logger.Errorln("create message window failed:", err)
		return err
	}
	return nil
}
