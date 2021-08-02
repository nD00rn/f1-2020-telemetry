package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "net"
    "os"
    "time"

    "github.com/nD00rn/f1-2020-telemetry/rest"
    "github.com/nD00rn/f1-2020-telemetry/statemachine"
    "github.com/nD00rn/f1-2020-telemetry/terminal"
    "github.com/nD00rn/f1-2020-telemetry/websocket"
)

var localSm *statemachine.StateMachine

func main() {
    restOptions, terminalOptions, f1Options := processOptions()
    fmt.Printf("terminal width is %d\n", terminalOptions.Width)

    addr := net.UDPAddr{
        IP:   net.ParseIP("0.0.0.0"),
        Port: int(f1Options.udpPort),
    }

    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    csm := statemachine.CreateCommunicationStateMachine()
    sm := statemachine.CreateStateMachine()
    localSm = &sm

    // Start the REST server
    rest.SetUpRestApiRouter(restOptions, &sm)

    // Set up websocket
    go websocket.Broadcast()
    go constantStreamWebSocket()

    for {
        // the number of incoming bytes lays around the 1500 bytes per request.
        // this should give some space whenever a bigger request comes in.
        buffer := make([]byte, 4096)
        n, _, err := conn.ReadFrom(buffer)
        if err != nil {
            panic(err)
        }

        csm.Process(buffer[:n], &sm)

        textBuf := terminal.CreateBuffer(sm, csm, terminalOptions)

        // Write the terminal output to file if enabled in the settings
        if len(terminalOptions.Path) > 0 {
            f, err := os.Create(
                fmt.Sprintf("%s/f1.screen.tmp", terminalOptions.Path),
            )

            if err != nil {
                panic(err)
            }
            _, _ = fmt.Fprint(f, textBuf)
            _ = f.Close()

            _ = os.Rename(
                fmt.Sprintf("%s/f1.screen.tmp", terminalOptions.Path),
                fmt.Sprintf("%s/f1.screen.txt", terminalOptions.Path),
            )
        }
    }
}

func processOptions() (
    rest.Options,
    terminal.Options,
    F1Options,
) {
    restOptions := rest.DefaultOptions()
    terminalOptions := terminal.DefaultOptions()
    f1Options := defaultF1Options()

    flag.UintVar(
        &restOptions.Port,
        "restport",
        restOptions.Port,
        "port to allow REST communication",
    )
    flag.UintVar(
        &terminalOptions.Width,
        "terminalwidth",
        terminalOptions.Width,
        "terminal width used to generate table",
    )
    flag.UintVar(
        &f1Options.udpPort,
        "udpport",
        f1Options.udpPort,
        "udp port to capture data on",
    )
    flag.StringVar(
        &terminalOptions.Path,
        "terminalpath",
        "",
        "path to store the terminal file (.tmp and .txt)",
    )

    flag.Parse()

    return restOptions, terminalOptions, f1Options
}

func constantStreamWebSocket() {
    for {
        time.Sleep(250 * time.Millisecond)

        marshal, err := json.Marshal(localSm)
        if err != nil {
            continue
        }

        websocket.BroadcastMessage(string(marshal))
    }
}
