/* sensors.go - Utils for reading sensor data 
 * Copyright (C) 2013  Blake Mitchell 
 */
package sensors


import (
        "io/ioutil"
        "strconv"
        "strings"
)

import (
        "../conf"
)


type (

        SensorReading struct {
                Data int    `json:"data"`
                Err  string `json:"err"`
        }

        SensorReadings map[string]SensorReading
)


func ReadSensors(sensors conf.SensorConf) SensorReadings {
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
