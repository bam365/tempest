/* hist.go - Defines operations on a CSV database history file
 * Copyright (C) 2013  Blake Mitchell 
 */
package main

import (
        "os"
        "encoding/csv"
        "sync"
        "time"
        "errors"
)


const DateTimeFormat = time.RFC1123Z
//According to golang docs, setting CSV reader's FieldsPerRecord
//to a negative number will cause it to read variable-length
//records
const CSV_VariableRecordLength = -1


type HistFile struct {
        hflock *sync.Mutex
        FileName string
}


func OpenHistFile(fname string) HistFile {
        return HistFile { new(sync.Mutex), fname }
}


func (hf *HistFile) WriteStartTime(t time.Time) error {
        return hf.Write([]string { t.Format(DateTimeFormat) })
}


func (hf *HistFile) readAllRecords() ([][]string, error) {
        hf.hflock.Lock()
        defer hf.hflock.Unlock()

        flags := os.O_RDONLY
        file, err := os.OpenFile(hf.FileName, flags, 0666) 
        if err != nil {
                return nil, err
        }
        defer file.Close()

        rdr := csv.NewReader(file)
        rdr.FieldsPerRecord = CSV_VariableRecordLength 
        return rdr.ReadAll()
}


func (hf *HistFile) ReadStartTime() (time.Time, error) {
        ret, rerr := time.Now(), error(nil)
        if recs, err := hf.readAllRecords(); err != nil {
                rerr = err
        } else {
                if len(recs) < 1 {
                        rerr = errors.New("No time record present")
                } else {
                        trec := recs[0][0]
                        ret, rerr = time.Parse(DateTimeFormat, trec)
                }
        }
        return ret, rerr
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
