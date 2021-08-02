package terminal

import (
    "fmt"
    "strings"

    "github.com/nD00rn/f1-2020-telemetry/statemachine"
)

func drawTrackProcess(
    sm statemachine.StateMachine,
    playerIndex uint8,
    terminalWidth uint,
    textBuf *string,
) {
    if playerIndex == 255 {
        return
    }
    name := string(sm.ParticipantData.Participants[playerIndex].Name[:])[0:3]

    lap := sm.LapData.LapData[playerIndex]
    trackLength := sm.SessionData.TrackLength

    lapDistance := lap.LapDistance
    if lapDistance < 0 {
        lapDistance = 0
    }
    lapPercentage := lapDistance / float32(trackLength)
    if lapPercentage >= 1 {
        lapPercentage = 0.99
    }

    availableBlockSpace := terminalWidth - 7
    current := getCorrectTrackPixel(lapPercentage, availableBlockSpace)

    passed := current
    remaining := availableBlockSpace - current - 1

    if remaining < 0 {
        remaining = 0
    }
    if passed < 0 || passed > terminalWidth {
        passed = 0
    }

    arrowHead := ">"
    *textBuf += fmt.Sprintf("[%-3s]|"+strings.Repeat(" ", int(passed))+arrowHead, name)
    *textBuf += fmt.Sprintf(strings.Repeat(" ", int(remaining)) + "|\n")
}

func drawMarshalZones(
    sm statemachine.StateMachine,
    terminalWidth uint,
    textBuf *string,
    showOtherPlayers bool,
) {
    numZones := sm.SessionData.NumMarshalZones
    zones := sm.SessionData.MarshalZones

    *textBuf += "[FLG]|"
    usableTerminalWidth := terminalWidth - 7

    for currentPixel := uint(0); currentPixel < usableTerminalWidth; currentPixel++ {
        zoneToUse := 0
        trackProgress := float32(1) / float32(usableTerminalWidth) * float32(currentPixel)

        trackLength := sm.SessionData.TrackLength
        trackText := " "

        if showOtherPlayers {
            for playerIndex := uint8(0); playerIndex < 22; playerIndex++ {
                if playerIndex == sm.LapData.Header.PlayerCarIndex || playerIndex == sm.LapData.Header.SecondPlayerCarIndex {
                    continue
                }
                if sm.LapData.LapData[playerIndex].ResultStatus == 0 {
                    continue
                }
                lapPercentage := sm.LapData.LapData[playerIndex].LapDistance / float32(trackLength)

                if isCorrectTrackPixel(lapPercentage, currentPixel, usableTerminalWidth) {
                    trackText = ">"
                    break
                }
            }
        }

        // Determine which zone data to use
        for i, _ := range zones {
            if zones[0].ZoneStart == 0 {
                break
            }

            if isCorrectZone(i, numZones, trackProgress, zones) {
                zoneToUse = i
                break
            }
        }

        trackColour := BgReset

        switch zones[zoneToUse].ZoneFlag {
        case 1:
            // trackText = " "
            trackColour = BgGreen
            break
        case 2:
            // trackText = " "
            trackColour = BgBlue
            break
        case 3:
            // trackText = " "
            trackColour = BgYellow
            break
        case 4:
            // trackText = " "
            trackColour = BgRed
            break
        }

        *textBuf += trackColour + trackText
    }
    *textBuf += BgReset + "|\n"
}
