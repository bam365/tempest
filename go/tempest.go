/* tempest.go - Tempest data run operations 
 * Copyright (C) 2013  Blake Mitchell 
 */

package main


import (
        "fmt"
        "time"
        "errors"
        "os"
        "regexp"
        "path/filepath"
)

import (
        "./config"
        "./sensors"
)


const (
        //Make sure this ends with a slash
        RunHistDir = "runhist/"
        RunHistFileName = "runhist"
        RunHistFileExt = ".csv"
)

var ( 
        CurrentRunHistFile = RunHistDir + RunHistFileName + RunHistFileExt
)


type (
        TempestRunner interface {
                TimeStarted() time.Time
                TimeEnded() time.Time
                IsRunning() bool
                RunDuration() time.Duration 
                Hist() *HistFile
        }
        

        TempestRun struct {
                StartTime  time.Time
                EndTime    time.Time //Will be Time(0) if current
                History    HistFile
        }


        CurrentTempestRun struct {
                TempestRun
                Conf     config.TempestConf
                stop     chan bool 
                err      chan string
                alert    chan string
        }
)


func RunsList() []TempestRunner {
        ret := make([]TempestRunner, 0)
        pattern := fmt.Sprintf("%s%s*%s", RunHistDir, RunHistFileName, 
        	               RunHistFileExt)
        if files, ferr := filepath.Glob(pattern); ferr == nil {
        	for _, f := range(files) {
        		if tr, err := LoadTempestRun(f); err == nil {
        			ret = append(ret, tr)
        		}
        	}
        }
        return ret
}



func TempestEncodeTime(t time.Time) string {
        return fmt.Sprintf("%02d%02d%02d%02d%02d%02d", 
                           t.Year(), t.Month(), t.Day(),
                           t.Hour(), t.Minute(), t.Second())
}


func TempestDecodeTime(t string) (time.Time, error) {
        //The layout string is cryptic, it means YYMMDDhhmmss
        loc, _ := time.LoadLocation("Local")
        return time.ParseInLocation("20060102150405", t, loc)
}


func IsRunInProgress() bool {
        _, err := os.Stat(CurrentRunHistFile)
        return (err == nil)
}


func intervalTicker(interval int) <-chan time.Time {
        return time.Tick(time.Duration(interval) * time.Second)
}



func alertMessage(sname, msg string) {
        fmt.Sprintf("ALERT: %s sensor: %s", sname, msg)
}



// TempestRun stuff

func LoadTempestRun(fname string) (*TempestRun, error) {
        tr := new(TempestRun)
        curr_restr := fmt.Sprintf("%s%s%s", 
                             RunHistDir, RunHistFileName, RunHistFileExt)
        past_restr := fmt.Sprintf("%s%s-([0-9]{14})%s", 
                             RunHistDir, RunHistFileName, RunHistFileExt)
        curr_re := regexp.MustCompile(curr_restr)
        past_re := regexp.MustCompile(past_restr)
        if past_re.MatchString(fname) {
	        endtimestr := past_re.FindStringSubmatch(fname)[1]
	        if endtime, err := TempestDecodeTime(endtimestr); err != nil {
	                return nil, err
	        } else {
	                tr.EndTime = endtime
	        }
	} else if !curr_re.MatchString(fname) {
                return nil, errors.New("Not a run history file")
        }
        tr.History = OpenHistFile(fname)
        if starttime, err := tr.Hist().ReadStartTime(); err != nil {
                return nil, err
        } else {
                tr.StartTime = starttime
        }
        return tr, nil
}


func (tr TempestRun) TimeStarted() time.Time {
        return tr.StartTime
}


func (tr TempestRun) TimeEnded() time.Time {
        return tr.EndTime
}


func (tr TempestRun) RunDuration() time.Duration {
        var endtime time.Time
        if tr.IsRunning() {
                endtime = time.Now()
        } else {
                endtime = tr.TimeEnded()
        }
        return endtime.Sub(tr.TimeStarted())
}


func (tr TempestRun) IsRunning() bool {
        return tr.TimeEnded().IsZero()
}


func (tr TempestRun) Hist() *HistFile {
        return &(tr.History)
}

        
// CurrentTempestRun stuff

func newCurrentTempestRun(td *TempestData) *CurrentTempestRun {
        return &CurrentTempestRun { 
                TempestRun: TempestRun {
                        History: OpenHistFile(CurrentRunHistFile),
                },
                Conf: *td.Conf,
                stop: make(chan bool),
                err: make(chan string),
                alert: td.Alert,
        }
}


func StartNewTempestRun(td *TempestData) (*CurrentTempestRun, error) {
        if IsRunInProgress() {
                return nil, errors.New("Run already in progress")
        }
        ret := newCurrentTempestRun(td)
        ret.StartTime = time.Now()
        if err := ret.History.WriteStartTime(ret.StartTime); err != nil {
                return nil, err
        }
        if err := ret.resumeRun(); err != nil {
                return nil, err
        }
        return ret, nil
}


func ResumeCurrentTempestRun(td *TempestData) (*CurrentTempestRun, error) {
        if !IsRunInProgress() {
                return nil, errors.New("There is no current run to resume")
        }
        ret := newCurrentTempestRun(td)
        if starttime, err := ret.History.ReadStartTime(); err != nil {
                return nil, err
        } else {
                ret.StartTime = starttime
        }
        if err := ret.resumeRun(); err != nil {
                return nil, err
        }
        return ret, nil
}



func (tr *CurrentTempestRun) resumeRun() error {
        if !tr.IsRunning() {
                return errors.New("Not running") 
        }
        go tr.histRecorderProc()
        go tr.alerterProc()
        return nil
}



func (tr *CurrentTempestRun) StopRun() error {
        if (!tr.IsRunning()) {
                return errors.New("Not running")
        }
        tr.stop <- true
        tr.EndTime = time.Now()
        fn, et := CurrentRunHistFile, tr.EndTime
        os.Rename(fn, fmt.Sprintf("%s%s-%s%s", RunHistDir, RunHistFileName,
                                  TempestEncodeTime(et), RunHistFileExt))  
        return nil
}



func (tr *CurrentTempestRun) alerterProc() {
        tmr := intervalTicker(tr.Conf.AlertInterval)
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
                for sname, sdat := range(sensors.ReadSensors(tr.Conf.Sensors)) {
                        if amsg := alertmsg(tr.Conf.Sensors[sname].Alert, sdat); amsg != "" {
                                tr.alert <- fmt.Sprintf("Sensor %s: %s", sname, amsg )
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


func (tr *CurrentTempestRun) histRecorderProc() {
        tmr := intervalTicker(tr.Conf.HistInterval)
        writerec := func (t int) { 
                readings := sensors.ReadSensors(tr.Conf.Sensors)
                if err := tr.History.Write(readings.ToCSVRecord(t)); err != nil {
                        tr.err <- err.Error()
                }
        }
        tbegin := tr.TimeStarted() 
        record := func(t time.Time) {
        	writerec(int(t.Sub(tbegin).Seconds()))
        }
        record(time.Now())
        for {
                select {
                case quit := <-tr.stop:
                        if (quit) {
                                return
                        }
                case now := <-tmr:
                        record(now) 
                }
        }
}


