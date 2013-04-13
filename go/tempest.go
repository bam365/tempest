/* tempest.go - Tempest data run operations 
 * Copyright (C) 2013  Blake Mitchell 
 */

package main


import (
        "fmt"
        "time"
        "errors"
        "os"
)

import (
        "./config"
        "./sensors"
)


const DefaultRunHistFile = "runhist.csv"


type TempestRun struct {
        filename string
        histfile HistFile
        conf     config.TempestConf
        start    time.Time
        stop     chan bool 
        err      chan string
}


func NewTempestRun(fname string, conf config.TempestConf) *TempestRun {
        return &TempestRun { 
                filename: fname,
                conf: conf,
                stop: make(chan bool),
                err: make(chan string),
        }
}


func (tr *TempestRun) IsRunning() bool {
        _, err := os.Stat(tr.filename)
        return (err == nil)
}


func (tr *TempestRun) ResumeRun() error {
        if (!tr.IsRunning()) {
                return errors.New("Not running")
        }
        if st, err := tr.histfile.ReadStartTime(); err != nil {
                return err
        } else {
                tr.start = st
        }
        go tr.histRecorderProc()
        go tr.alerterProc()
        return nil
}


func (tr *TempestRun) StartRun() error {
        if (tr.IsRunning()) {
                return errors.New("Already running") 
        }
        tr.histfile = OpenHistFile(tr.filename)
        tr.histfile.WriteStartTime(time.Now())
        return tr.ResumeRun()
}
        

func (tr *TempestRun) StopRun() error {
        if (!tr.IsRunning()) {
                return errors.New("Not running")
        }
        tr.stop <- true
        fn, st := tr.filename, tr.start
        os.Rename(fn, fmt.Sprintf("%s-%02d%02d%02d%02d%02d%02d", fn, 
                                  st.Year(), st.Month(), st.Day(),
                                  st.Hour(), st.Minute(), st.Second()))
        return nil
}



func (tr *TempestRun) alerterProc() {
        tmr := intervalTicker(tr.conf.AlertInterval)
        alertmsg := func(arange config.SensorRange, sdat sensors.SensorReading) string {
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
                for sname, sdat := range(sensors.ReadSensors(tr.conf.Sensors)) {
                        if amsg := alertmsg(tr.conf.Sensors[sname].Alert, sdat); amsg != "" {
                                alert(sname, amsg)
                        }
                }
        }
                                
        for {
                select {
                case quit := <-tr.stop:
                        if (quit) {
                                return
                        }
                case <-tmr:
                        checkalerts()
                }
        }
}


func (tr *TempestRun) histRecorderProc() {
        tmr := intervalTicker(tr.conf.HistInterval)
        writerec := func (t int) { 
                readings := sensors.ReadSensors(tr.conf.Sensors)
                if err := tr.histfile.Write(readings.ToCSVRecord(t)); err != nil {
                        tr.err <- err.Error()
                }
        }
        tbegin := time.Now()
        writerec(0)
        for {
                select {
                case quit := <-tr.stop:
                        if (quit) {
                                return
                        }
                case now := <-tmr:
                        writerec(int(now.Sub(tbegin).Seconds())) 
                }
        }
}


func intervalTicker(interval int) <-chan time.Time {
        return time.Tick(time.Duration(interval) * time.Second)
}


func alert(sname, msg string) {
        fmt.Printf("ALERT: %s sensor: %s\n", sname, msg)
        //TODO:  Send some emails too
}

