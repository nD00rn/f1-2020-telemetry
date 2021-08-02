package terminal

import (
    "fmt"
)

var lastStreamBuffer string = fmt.Sprintf("%s%s%s%s",
    BgRed,
    "hello world\nnice",
    BgReset,
    "\n",
)

func GetLastBuffer() string {
    return lastStreamBuffer
}
