/* conf.go - Configuration for tempest
 * Copyright (C) 2013  Blake Mitchell 
 */

package conf


import (
        "io/ioutil"
        "encoding/json"
)

type (

        SensorRange struct {
                Low int  `json:"lo"`
                High int `json:"hi"`
        }


        SensorData struct {
                File string       `json:"file"`
                Range SensorRange `json:"range"` 
                Alert SensorRange `json:"alert"`
        }

        SensorConf map[string]SensorData


        TempestConf struct {
                Sensors SensorConf `json:"sensors"`
                Emails []string    `json:"emails"`
        }
)


func NewTempestConf() TempestConf {
        return TempestConf{ make(map[string]SensorData) }
}


func LoadConf(fname string) (TempestConf, error) {
        var rerr error = nil
        conf := NewTempestConf() 
        if bytes, err := ioutil.ReadFile(fname); err != nil {
                rerr = err 
        } else {
                if jerr := json.Unmarshal(bytes, &conf); jerr != nil {
                        rerr = jerr
                }
        }

        return conf, rerr 
}

