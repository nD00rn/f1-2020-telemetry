package statemachine

import (
  "testing"
)

func TestProcessPacketEvent(t *testing.T) {
  state := StateMachine{
    []byte{
      0xe4, 0x07, 0x01, 0x13, 0x01, 0x03, 0x5d, 0xc9,
      0xb5, 0x0a, 0x53, 0x0c, 0xf3, 0x11, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x13,
      0x53, 0x53, 0x54, 0x41, 0xc8, 0xd8, 0x12, 0x0c,
      0x00, 0x00, 0x00,
    },
  }

  // 20: c8d8 120c 0000 00xx
  // 21: 5894 4036 0000 0000

  ProcessPacketEvent(&state)
}
