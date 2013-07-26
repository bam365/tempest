package main


import (
	"fmt"
	"encoding/json"
	"strconv"
	"os"
	"os/exec"
	"io/ioutil"
	"math"
	"time"
	"strings"
)


const MCP3008Ticks = 1024


type (
	Thermistor struct {
		Channel byte    `json:"channel"`
		A float32       `json:"a"`
		B float32       `json:"b"`
		C float32       `json:"c"`
		R float32	`json:"r"`
		File string	`json:"file"`
	}


	Config struct {
		Thermistors []Thermistor `json:"thermistors"`
		ReadSpi string `json:"readspi"`
		Fahrenheit bool `json:"fahrenheit",omitempty`
	}

)



func main() {
	conffile := "therm.conf"
	if conf, err := LoadConfig(conffile); err != nil {
		ReportError(err)
	} else {
		for _, therm := range(conf.Thermistors) {
			go ThermistorDriver(therm, conf)
		}
		BlockForever()
	}
}


func NewConfig() Config {
	return Config {
		Thermistors: make([]Thermistor, 0),
		Fahrenheit: false,
	}
}


func LoadConfig(fname string) (Config, error) {
        var rerr error = nil
        conf := NewConfig() 
        if bytes, err := ioutil.ReadFile(fname); err != nil {
                rerr = err 
        } else {
                if jerr := json.Unmarshal(bytes, &conf); jerr != nil {
                        rerr = jerr
                }
        }

        return conf, rerr 
}


func (t Thermistor) CalcThermResistance(rdg int) float32 {
	V_in := MCP3008Ticks
	V_ratio := float32(V_in) / float32(rdg)
	return (t.R * V_ratio - t.R)
}


func (t Thermistor) ResistanceToKelvin(r float32) float32 {
	var ret float32 = -1.0 
	if r > 0 {	
		R_log := math.Log(float64(r))
		R_inv := t.A + t.B*float32(R_log) + t.C*float32(math.Pow(R_log, 3.0))
		ret = (1.0 / R_inv)
	}
	return ret
}


func (t Thermistor) GetRawReading(conf Config) (int, error) {
	rdg, rerr := -1, error(nil)
	buf, err := exec.Command(conf.ReadSpi, strconv.Itoa(int(t.Channel))).Output()
	if err == nil {
		instr := strings.Trim(string(buf), " \n\t\r")
		if n, cerr := strconv.Atoi(instr); cerr == nil {
			rdg = n
		}
	} else {
		rerr = err
	}
	return rdg, rerr
}


func ReportError(err error) {
	fmt.Printf("thermdrv: ERROR: %s\n", err.Error())
}


func BlockForever() {
	foo := make(chan bool)
	<-foo
}


func KtoC(k float32) float32 {
	return (k - 273.15)
}


func CtoF(c float32) float32 {
	return (1.8 * c + 32.0)
}


func Round(n float32) int {
	return int(math.Floor(float64(n) + 0.5))
}


func UpdateFile(fname string, rdg int) error {
	f, err := os.OpenFile(fname, os.O_TRUNC | os.O_CREATE | os.O_WRONLY, 0666)
	if err == nil {
		fmt.Fprintf(f, "%d\n", rdg)
		f.Close()
	} else {
		ReportError(err)
	}
	return err
}


func ThermistorDriver(therm Thermistor, conf Config) {
	gettemp := func() int {
		ret := -1
		if rdg, err := therm.GetRawReading(conf); err == nil {
			tempK := therm.ResistanceToKelvin(therm.CalcThermResistance(rdg))
			temp := KtoC(tempK)
			if conf.Fahrenheit {
				temp = CtoF(temp)	
			}
			ret = Round(temp)
		} else {
			ReportError(err)
		}
		return ret 
	}

	readAndUpdate := func() {
		if err := UpdateFile(therm.File, gettemp()); err != nil {
			ReportError(err)	
		}
	}

	for {
		readAndUpdate()
		time.Sleep(3.0 * time.Second)
	}
}
