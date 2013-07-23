package main

import (
        "./sensors"
)


type (
        
        PlotPoint struct {
                X int `json:"x"`
                Y int `json:"y"`
        }

        PlotData []PlotPoint
        
        TempestPlot map[string]PlotData

)


func NewTempestPlot() TempestPlot {
        return make(map[string]PlotData)
}


func (tp TempestPlot) AddToSeries(series string, dat PlotPoint) {
        if _, found := tp[series]; !found {
                tp[series] = make([]PlotPoint, 0)
        }
        tp[series] = append(tp[series], dat)
}
        


func ReadPlotData(hf HistFile, offset, interval int) (TempestPlot, error) {
        tp := NewTempestPlot() 

        addpoint := func(sname string, t, dat int) {
                tp.AddToSeries(sname, PlotPoint {t, dat})
        }

        //Adds a point to series sname record, provided
        //that the last point we added is beyond the interval
        addreading := func(t int, sname string, sdat sensors.SensorReading) {
                if sdat.Err == "" {
                        dat := sdat.Data
                        if pd, found := tp[sname]; found && len(pd) > 0 {
                                tlast := pd[len(pd) - 1].X
                                if t - tlast >= interval {
                                        addpoint(sname, t, dat)
                                }
                        } else {
                                addpoint (sname, t, dat)
                        }
                }
        }
        

        allreadings, rerr := hf.ReadAllRecords()
        if rerr != nil {
                return nil, rerr 
        }

        for _, rawreading := range(allreadings) {
                secs, rdg, cerr := sensors.ReadingsFromCSVRecord(rawreading)
                if cerr == nil && secs >= offset {
                        for sname, sdat := range(rdg) {
                                addreading(secs, sname, sdat)
                        }
                }
        }

        return tp, nil
}
