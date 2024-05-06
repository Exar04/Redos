package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
)

const defaultListenAddress = "0.0.0.0:8888"

type Message struct {
	cmd  Command
	peer *Peer
}

type Config struct {
	ListenAddress string
}

type Server struct {
	Config                      // the configuration data is unfold in another structure to keep it clean
	peers        map[*Peer]bool // peers are basically the client connections that are getting connected to the server
	ln           net.Listener   // tcp server connecton
	addPeerCh    chan *Peer     // channel to add new clients
	deletePeerCh chan *Peer     // channel to add delete clients
	quitCh       chan struct{}  // idk wtf is this
	msgCh        chan Message   // idk this either
	kv           *KeyVal        // our key value database
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultListenAddress
	}

	return &Server{
		Config:       cfg,
		peers:        make(map[*Peer]bool),
		addPeerCh:    make(chan *Peer),
		deletePeerCh: make(chan *Peer),
		quitCh:       make(chan struct{}),
		msgCh:        make(chan Message),
		kv:           NewKeyVal(),
	}
}

func (s *Server) start() error {
	ln, err := net.Listen("tcp", s.ListenAddress)
	if err != nil {
		fmt.Println("error while starting server")
		return err
	}

	s.ln = ln

	go s.loop()

	slog.Info("sever running", "listenAddr", s.ListenAddress)

	return s.acceptLoop()
}

func (s *Server) handleRawMessage(msg Message) error {
	switch v := msg.cmd.(type) {
	case SetCommand:
		return s.kv.Set(v.key, v.val)

	case GetCommand:
		val, ok := s.kv.Get(v.key)
		if !ok {
			return fmt.Errorf("key not found")
		}
		_, err := msg.peer.Send(val)
		if err != nil {
			slog.Error("peer send error", "err :", err)
		}
	}
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleRawMessage(msg); err != nil {
				slog.Error("raw message error ", "err", err)
			}
		case <-s.quitCh:
			return
		case peer := <-s.addPeerCh:
			slog.Info("new peer connected", "remoteAddr", peer.conn.RemoteAddr())
			s.peers[peer] = true
		case peer := <-s.deletePeerCh:
			slog.Info("peer disconnected", "remoteAddr", peer.conn.RemoteAddr())
			delete(s.peers, peer)
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error :: ", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh, s.deletePeerCh)
	s.addPeerCh <- peer
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}

func main() {
	listenAddr := flag.String("listenAddr", defaultListenAddress, "listen address of the redos server")
	flag.Parse()
	server := NewServer(Config{
		ListenAddress: *listenAddr,
	})
	log.Fatal((server.start()))
}
