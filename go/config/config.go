/* conf.go - Configuration for tempest
 * Copyright (C) 2013  Blake Mitchell 
 */

package config


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
                AlertInterval int  `json:"alertdelay"`
                HistInterval int  `json:"histdelay"`
        }
)


func NewTempestConf() TempestConf {
        return TempestConf { 
                Sensors: make(map[string]SensorData), 
                Emails: make([]string, 0),
                AlertInterval: 60,
                HistInterval: 60,
        }
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
