package main

type TerminalOptions struct {
    width uint
}

func defaultTerminalOptions() TerminalOptions {
    return TerminalOptions{
        width: 90,
    }
}
