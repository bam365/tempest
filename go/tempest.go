/* tempest.go - Web server for tempest
 * Copyright (C) 2013  Blake Mitchell 
 */

package main


import (
        "fmt"
        //"net/http"
        "encoding/json"
        "time"
)

import (
        "./conf"
        "./sensors"
)


const RunHistFile = "runhist.csv"


func JsonStr(v interface{}) string {
        b, _ := json.MarshalIndent(v, "", "\t")
        return string(b)
}


func Alert(sname, msg string) {
        fmt.Printf("ALERT: %s sensor: %s\n", sname, msg)
}


func IntervalTicker(interval int) <-chan time.Time {
        return time.Tick(time.Duration(interval) * time.Second)
}


func AlerterProc(c conf.TempestConf) {
        tmr := IntervalTicker(c.AlertInterval)
        alertmsg := func(arange conf.SensorRange, sdat sensors.SensorReading) string {
                msg := ""
                if sdat.Err != "" {
                        msg =  sdat.Err
                } else if sdat.Data < arange.Low {
                        msg =  fmt.Sprintf("Reading %d is below %d", sdat.Data, 
                                           arange.Low)
                } else if sdat.Data > arange.High {
                        msg =  fmt.Sprintf("Reading %d is above %d", sdat.Data, 
                                           arange.High)
                }
                return msg
        } 
        checkalerts := func() {
                for sname, sdat := range(sensors.ReadSensors(c.Sensors)) {
                        if amsg := alertmsg(c.Sensors[sname].Alert, sdat); amsg != "" {
                                Alert(sname, amsg)
                        }
                }
        }
                                
        for {
                <-tmr
                checkalerts()
        }
}


func HistRecorderProc(c conf.TempestConf, hist HistFile) {
        tmr := IntervalTicker(c.HistInterval)
        writerec := func (t int) { 
                readings := sensors.ReadSensors(c.Sensors)
                if err := hist.Write(readings.ToCSVRecord(t)); err != nil {
                        fmt.Printf("ERROR: %s", err.Error())
                }
        }
        tbegin := time.Now()
        writerec(0)
        for now := range(tmr) {
                writerec(int(now.Sub(tbegin).Seconds())) 
        }
}



func main() {
        if conf, err := conf.LoadConf("testconf"); err != nil {
                fmt.Printf("Error loading conf: %s\n", err)
        } else {
                fmt.Println(JsonStr(conf.Sensors))
                go HistRecorderProc(conf, OpenHistFile(RunHistFile))
                AlerterProc(conf)
        }
}


