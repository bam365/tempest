/* sensors.go - Utils for reading sensor data 
 * Copyright (C) 2013  Blake Mitchell 
 */
package sensors


import (
        "io/ioutil"
        "strconv"
        "strings"
        "errors"
)


import (
        "../config"
)


const RecordErrorField = "?"


type (

        SensorReading struct {
                Data int    `json:"data"`
                Err  string `json:"err"`
        }

        SensorReadings map[string]SensorReading
)


func ReadSensors(sensors config.SensorConf) SensorReadings {
        ret := make(map[string]SensorReading)
        for sname, sdat := range(sensors) {
               ret[sname] = ReadSensor(sdat.File)
        }
        return ret
}


func ReadSensor(fname string) SensorReading {
        ret := SensorReading { 0, "" }
        if in, err := ioutil.ReadFile(fname); err != nil {
                ret.Err = err.Error()
        } else {
                in_cleaned := strings.TrimSpace(string(in))
                if d, cerr := strconv.Atoi(in_cleaned); cerr != nil {
                        ret.Err = cerr.Error() 
                } else {
                        ret.Data = d
                }
        }
        return ret
}


func (sr SensorReadings) ToCSVRecord(secs int) []string {
        rec := make([]string, 0)
        rec = append(rec, strconv.Itoa(secs))
        for sname, sdat := range(sr) {
                rec = append(rec, sname)
                if sdat.Err == "" {
                        rec = append(rec, strconv.Itoa(sdat.Data))
                } else {
                        rec = append(rec, RecordErrorField)
                }
        }
        return rec
}


func ReadingsFromCSVRecord(rec []string) (int, SensorReadings, error) {
        ret := make(map[string]SensorReading)
        recl := len(rec)
        var rtime int
        var cerr error
        if recl < 1 || (recl - 1) % 2 != 0 {
                return 0, ret, errors.New("Malformed record")
        }
        if rtime, cerr = strconv.Atoi(rec[0]); cerr != nil {
                return 0, ret, cerr
        }        
        for i := 1; i < recl; i += 2 {
                if rec[i+1] == RecordErrorField {
                        ret[rec[i]] = SensorReading {0, "Value not known"}
                } else {
                        if reading, cerr := strconv.Atoi(rec[i+1]); cerr != nil {
                                return 0, ret, cerr
                        } else {
                                ret[rec[i]] = SensorReading {reading, ""}
                        }
                }
        }
        return rtime, ret, cerr
}
