package src

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"path"
	"strings"

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
	homePage        []byte
	fsHandler       http.Handler
	nodeModsHandler http.Handler
	state           *State
}

func (s server) serveFile(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if path == "/" || path == "/index.html" || path == "/index" {
		// TODO
	}
	if strings.HasPrefix(path, "/node_modules/") {
		log.Println("Serving node module", path)
		res.Header().Add("Content-Type", "text/javascript")
		s.nodeModsHandler.ServeHTTP(res, req)
	} else {
		log.Println("Serving file", path)
		s.fsHandler.ServeHTTP(res, req)
	}
}

// StartServer starts the main backend server.
//
// The /ws route is where the frontend should connect to the websocket. All other routes
// are served from the file system with a root at "<frontendDir>/modules/out".
//
// It returns a channel that sends the value `true` when the server stops, which only
// happens if it crashes.
func StartServer(frontendDir string, logger *slog.Logger, state *State) chan bool {
	modsDir := path.Join(frontendDir, ModulesDir, "out")
	mux := http.NewServeMux()

	s := server{
		fsHandler:       http.FileServer(http.Dir(modsDir)),
		nodeModsHandler: http.FileServer(http.Dir(frontendDir)),
		state:           state,
	}

	mux.Handle("/ws", websocket.Handler(wsServer))
	mux.HandleFunc("/", s.serveFile)

	// whether the process should be restarted
	result := make(chan bool)

	go func() {
		logger.Info(fmt.Sprintf("Running at http://localhost:%d\n", 8000))
		err := http.ListenAndServe(":8000", mux)
		logger.Warn("Server crashed", "error", err)
		result <- true
		close(result)
	}()

	return result
}
