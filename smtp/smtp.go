package smtp

import (
	"io"
	"net"

	"github.com/jerry-yuan/mail-blackhole/config"
	"github.com/jerry-yuan/mail-blackhole/data"
	"github.com/jerry-yuan/mail-blackhole/storage"
	"github.com/sirupsen/logrus"
)

func Listen(cfg *config.SMTPConfig, msgStorage storage.Storage, messageCh chan *data.Message, exitCh chan struct{}) *net.TCPListener {
	logrus.Infof("[SMTP] Binding to address: %s\n", cfg.BindAddr)
	ln, err := net.Listen("tcp", cfg.BindAddr)
	if err != nil {
		logrus.Fatalf("[SMTP] Error listening on socket: %s\n", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Errorf("[SMTP] Error accepting connection: %s\n", err)
			continue
		}

		go Accept(
			conn.(*net.TCPConn).RemoteAddr().String(),
			io.ReadWriteCloser(conn),
			msgStorage,
			messageCh,
			cfg.Hostname)
	}
}
