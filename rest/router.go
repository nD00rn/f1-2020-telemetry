package rest

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "github.com/nD00rn/f1-2020-telemetry/statemachine"
    "github.com/nD00rn/f1-2020-telemetry/websocket"
)

var sm *statemachine.StateMachine

type Options struct {
    Port uint
}

func DefaultOptions() Options {
    return Options{
        Port: 8000,
    }
}

func SetUpRestApiRouter(options Options, stateMachine *statemachine.StateMachine) {
    // Set up the state machine which we will need for the requests
    sm = stateMachine

    // set up REST API
    r := mux.NewRouter()
    r.HandleFunc("/", homeHandler).Methods("GET")
    r.HandleFunc("/socket", websocket.SocketHandler)
    http.Handle("/", r)

    srv := &http.Server{
        Handler: r,
        Addr:    fmt.Sprintf("%s:%d", "127.0.0.1", options.Port),

        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    // Run our server in a goroutine so that it doesn't block.
    go func() {
        fmt.Printf("starting REST server on port %v\n", srv.Addr)
        if err := srv.ListenAndServe(); err != nil {
            log.Println(err)
        }
    }()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(sm)
}
