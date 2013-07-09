/* main.go - Main source and web server for tempest 
 * Copyright (C) 2013  Blake Mitchell 
 */
package main

import (
        "fmt"
        "net/http"
        "io/ioutil"
        "encoding/json"
        "html/template"
)

import (
        "github.com/gorilla/mux"
        "./sensors" 
)


const (
        SensorFile = "html/sensors.html"
        ReadingsFile = "html/readings.html"
        WebHistFile = "html/hist.html"
        RunsTemplateFile = "html/templ/runs.html"
)

type (
        PageHandler  func (http.ResponseWriter, *http.Request) 
        AjaxHandler  func (string) (string, error)


        WebServer struct {
                TData *TempestData
                URLMappings map[string]PageHandler
                AjaxMappings map[string]AjaxHandler
        }

)


type ( //Ajax requests
        HistRequest struct {
                Offset   int `json:"offset,omitempty"`
                Interval int `json:"interval,omitempty"`
        }
)


type ( //View models
	RunModel struct {
		IsRunning bool
		Status     string
		StartTime  string
		EndTime    string
		Duration   string
		ButtonText string
		ButtonLink string
	}
)



func NewWebServer(td *TempestData) *WebServer {
        ws := new(WebServer)     

        ws.TData = td
        ws.URLMappings = map[string]PageHandler {
                "/sensors": StaticFileServer(SensorFile), 
                "/readings": StaticFileServer(ReadingsFile),
                "/hist": StaticFileServer(WebHistFile),
                "/ajax/{request}": ws.HandleAjax,
                "/js/tempest/{file}": StaticFileServerFromVar("file", "html/js/tempest"),
                "/stylesheets/{file}": StaticFileServerFromVar("file", 
                	                                       "html/stylesheets"),
                "/runs": ws.RunsHandler,
        }

        ws.AjaxMappings = map[string]AjaxHandler {
                "readings": ws.AjaxReadings,
                "hist":     ws.AjaxHist,
                "sensors":  ws.AjaxSensors,
        }

        http.Handle("/", ws.SetupUrlRouter())
        return ws
}


func (ws *WebServer) Run() error {
        addr := fmt.Sprintf(":%d", ws.TData.Conf.Port) 
        return http.ListenAndServe(addr, nil)
}

       
func (ws *WebServer) SetupUrlRouter() *mux.Router {
        router := mux.NewRouter()
        for url, hdlr := range(ws.URLMappings) {
                router.HandleFunc(url, hdlr)
        }
        return router
}


func ServeFile(fname string, w http.ResponseWriter) {
        if buf, err := ioutil.ReadFile(fname); err != nil {
                http.Error(w, err.Error(), http.StatusNotFound)
        } else {
                fmt.Fprint(w, string(buf))
        }
}


func StaticFileServer(fname string) PageHandler {
        return func(w http.ResponseWriter, r *http.Request) {
                ServeFile(fname, w)
        }
}


func StaticFileServerFromVar(varname, basepath string) PageHandler {
        return func(w http.ResponseWriter, r *http.Request) {
                vars := mux.Vars(r)
                fname := basepath + "/" + vars[varname]
                ServeFile(fname, w)
        }
}


func ParseRequestBody(r *http.Request) (string, error) {
        ret, rerr := "", error(nil)
        if buf, err := ioutil.ReadAll(r.Body); err != nil {
                rerr = err
        } else {
                ret = string(buf)
        }

        return ret, rerr
}


func (ws *WebServer) HandleAjax(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        ajaxreq := vars["request"]

        ajaxhdlr, found := ws.AjaxMappings[ajaxreq]
        if !found {
                http.Error(w, "Request not found", http.StatusNotFound)
                return
        } 
        body, perr := ParseRequestBody(r)  
        if perr != nil {
                http.Error(w, perr.Error(), http.StatusBadRequest)
                return
        }
        response, rerr := ajaxhdlr(body)
        if rerr != nil {
                http.Error(w, rerr.Error(), http.StatusInternalServerError)
                return 
        }
        
        fmt.Fprintln(w, response)
}


func (ws *WebServer) AjaxReadings(arg string) (ret string, rerr error) {
        readings := sensors.ReadSensors(ws.TData.Conf.Sensors)
        return JsonStr(readings), nil
}


func (ws *WebServer) AjaxHist(arg string) (ret string, rerr error) {
        hf := ws.TData.Run.History
        req := HistRequest {}
        json.Unmarshal([]byte(arg), &req)
        tp, err := ReadPlotData(hf, req.Offset, req.Interval)
        if err != nil {
                ret, rerr = "", err
        } else {
                if buf, merr := json.Marshal(tp); merr == nil {
                        ret = string(buf)
                } else {
                        rerr = merr 
                }
        }
        return
}


func (ws *WebServer) AjaxSensors(arg string) (ret string, rerr error) {
	sensors := ws.TData.Conf.Sensors
        return JsonStr(sensors), nil
}


func RunsListToModels(runslist []TempestRunner) []RunModel {
	//TODO: Should I make consts for some of these URLS?
	ret := make([]RunModel, 0)
	for _, run := range(runslist) {
		et := TimeStr(run.TimeEnded())
		stat := "Finsihed"
		link := fmt.Sprintf("/runs/%s", TempestEncodeTime(run.TimeEnded()))
		if run.IsRunning() {
			et = "-"
			stat = "Running"
			link = "/"
		}
		mdl := RunModel {
			IsRunning: run.IsRunning(),
			Status: stat,
			StartTime: TimeStr(run.TimeStarted()),
			EndTime: et,
			Duration: DurationStr(run.RunDuration()),
			ButtonText: "View",
			ButtonLink: link,

		}
		ret = append(ret, mdl)
	}
	return ret
}
                

func (ws *WebServer) RunsHandler(w http.ResponseWriter, r *http.Request) {
	if t, terr := template.ParseFiles(RunsTemplateFile); terr != nil {
		http.Error(w, "Could not load view", http.StatusInternalServerError)
	} else {
		runs := RunsListToModels(RunsList())
		t.Execute(w, runs)
	}
}
