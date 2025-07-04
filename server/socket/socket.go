package socket

import (
	"email_verify/respond"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsocket struct {
	w          http.ResponseWriter
	r          *http.Request
	conn       *websocket.Conn
	isAuth     bool
	eventMap   map[string]func([]byte)
	onceEvents map[string]struct{}
}

type Socket interface {
	On(string, func([]byte))
	Once(string, func([]byte))
	Close()
	Emit(string, string)
	EmitErr(string, string) interface{ Close() }
	Listen() error
}

func EmitWs[T any](s Socket, ev string, obj T) {
	if s == nil {
		return
	}

	o, err := json.Marshal(obj)

	if err != nil {
		log.Fatal(err.Error())
	}

	s.Emit(ev, string(o))
}

type socketMsg struct {
	EventName string `json:"eventName"`
	Data string `json:"data"`
}

func NewWebSocket(w http.ResponseWriter, r *http.Request) (Socket, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	s := &wsocket{
		w:          w,
		r:          r,
		conn:       conn,
		eventMap:   make(map[string]func([]byte)),
		onceEvents: make(map[string]struct{}),
	}

	return s, nil
}

func (s *wsocket) On(evName string, f func([]byte)) {
	s.eventMap[evName] = f
}

func (s *wsocket) Once(evName string, f func([]byte)) {
	s.onceEvents[evName] = struct{}{}
	s.On(evName, f)
}

func (s *wsocket) Close() {
	s.conn.Close()
}

func (s *wsocket) Emit(evName string, data string) {
	res, err := json.Marshal(socketMsg{EventName: evName, Data: data})
	if err != nil {
		log.Fatal(err.Error())
	}
	s.conn.WriteMessage(websocket.TextMessage, res)
}

func (s *wsocket) EmitErr(evName string, errMsg string) interface{ Close() } {
	data, _ := json.Marshal(respond.ResponseStruct{Err: true, Msg: errMsg})
	s.Emit(evName, string(data))
	return s
}

func (s *wsocket) Listen() error {
	for {
		var msg socketMsg
		err := s.conn.ReadJSON(&msg)
		if err != nil {
			return err
		}
		if f, ok := s.eventMap[msg.EventName]; ok {
			f([]byte(msg.Data))
			if _, ok = s.onceEvents[msg.EventName]; ok {
				delete(s.onceEvents, msg.EventName)
				delete(s.eventMap, msg.EventName)
			}
		}
	}
}
