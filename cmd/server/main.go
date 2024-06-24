package main

import (
	"bytes"
	"log"
	"log/slog"
	"net"
)

type Config struct {
	ListenAddr string
}

type Message struct {
	data []byte
	peer *Peer
}

type Server struct {
	Config
	peers     map[*Peer]struct{}
	addPeerCh chan *Peer
	ln        net.Listener
	quit      chan struct{}
	msgCh     chan Message

	kv *KV
}

func NewServer(cfg Config) *Server {
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = ":5000"
	}

	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]struct{}),
		addPeerCh: make(chan *Peer),
		quit:      make(chan struct{}),
		msgCh:     make(chan Message),

		kv: NewKV(),
	}
}

func (srv *Server) Start() error {
	ln, err := net.Listen("tcp", srv.ListenAddr)
	if err != nil {
		return err
	}
	srv.ln = ln

	go srv.loop()

	slog.Info("server running", "address", srv.ListenAddr)
	return srv.acceptLoop()
}

func (srv *Server) handleRawMessage(rawMsg []byte) error {
	cmd, err := parseCommand(bytes.NewReader(rawMsg))
	if err != nil {
		return err
	}

	switch v := cmd.(type) {
	case *SetCommand:
		slog.Info("somebody want to set a key", "key", v.Key, "value", v.Val)
		// fmt.Println("AAAAAAAA", srv.kv.data)
		return srv.kv.Set(v.Key, v.Val)
	case *GetCommand:
		slog.Info("somebody want to get a key", "key", v.Key)
		_, err = srv.kv.Get(v.Key)
		return err
		// TODO: write back to conn
	}

	return nil
}

func (srv *Server) loop() {
	for {
		select {
		case <-srv.quit:
			return
		case rawMsg := <-srv.msgCh:
			if err := srv.handleRawMessage(rawMsg); err != nil {
				slog.Info("raw msg error", "error", err)
				continue
			}
		case peer := <-srv.addPeerCh:
			srv.peers[peer] = struct{}{}
		}
	}
}

func (srv *Server) acceptLoop() error {
	for {
		conn, err := srv.ln.Accept()
		if err != nil {
			slog.Error("accept error", "error", err)
			continue
		}

		go srv.handleConn(conn)
	}
}

func (srv *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, srv.msgCh)
	srv.addPeerCh <- peer
	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "error", err, "remoteAddr", conn.RemoteAddr())
	}
}

func main() {
	srv := NewServer(Config{})
	log.Fatal(srv.Start())
}
