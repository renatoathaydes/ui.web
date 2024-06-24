package src

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"golang.org/x/net/websocket"
)

// Websocket Server handler.
func wsServer(ws *websocket.Conn) {
	fmt.Println("Starting Websocket loop")
	defer fmt.Println("Exiting Websocket loop")
	for {
		msg, err := WsMessageFromJson(ws)
		if err != nil {
			res := msg.WsMessageResponse(err, false)
			res.WriteJson(ws)
			return
		}
		var res WsMessage
		if msg.MsgType == "rpc" {
			if cmd, ok := msg.Data.(string); ok {
				result, verr := runCommand(cmd, "js")
				if verr != nil {
					res = msg.WsMessageResponse(fmt.Sprintf("%v", verr), false)
				} else if result.Error != nil {
					res = msg.WsMessageResponse(result.Error, false)
				} else {
					res = msg.WsMessageResponse(result, true)
				}
			}
		} else if msg.MsgType == "notify" {
			fmt.Printf("Received notification: %s", msg.Data)
			// notify requires no response
			continue
		} else {
			res = msg.WsMessageResponse(fmt.Sprintf("Go got data: %s", msg.Data), true)
		}
		res.WriteJson(ws)
	}
}

type server struct {
	homePage  []byte
	fsHandler http.Handler
	state     *State
}

func (s server) serveFile(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if path == "/" || path == "/index.html" || path == "/index" {
		res.Write(s.homePage)
	} else {
		s.fsHandler.ServeHTTP(res, req)
	}
}

func readHomePage(dir string) []byte {
	res, err := os.ReadFile(path.Join(dir, "index.html"))
	if err != nil {
		log.Fatal("Could not read index.html in dir:", dir)
	}
	return res
}

func StartServer(path string, state *State) {
	home := readHomePage(path)
	mux := http.NewServeMux()

	s := server{
		fsHandler: http.FileServer(http.Dir(path)),
		homePage:  home,
		state:     state,
	}

	mux.Handle("/ws", websocket.Handler(wsServer))
	mux.HandleFunc("/", s.serveFile)

	fmt.Printf("Running at http://localhost:%d\n", 8000)

	// does not return
	_ = http.ListenAndServe(":8000", mux)
}
