package main

import (
        "os"
        "encoding/csv"
        "sync"
)


type HistFile struct {
        hflock *sync.Mutex
        FileName string
}


func OpenHistFile(fname string) HistFile {
        return HistFile { new(sync.Mutex), fname }
}


func (hf *HistFile) Write(rec []string) error {
        hf.hflock.Lock()
        defer hf.hflock.Unlock()

        flags := os.O_WRONLY | os.O_APPEND | os.O_CREATE
        file, err := os.OpenFile(hf.FileName, flags, 0666) 
        if err != nil {
                return err
        }
        defer file.Close()

        wtr := csv.NewWriter(file)
        werr := wtr.Write(rec)
        wtr.Flush()
        return werr
}








