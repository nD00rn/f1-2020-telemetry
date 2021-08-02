package main

type F1Options struct {
    udpPort uint
}

func defaultF1Options() F1Options {
    return F1Options{
        udpPort: 20777,
    }
}
