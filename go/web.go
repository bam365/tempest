/* main.go - Main source and web server for tempest 
 * Copyright (C) 2013  Blake Mitchell 
 */
package main

import (
        "fmt"
        "net/http"
        "io/ioutil"
)

import (
        "github.com/gorilla/mux"
)


const (
        SensorFile = "html/sensors.html"
)

type (
        PageHandler  func (http.ResponseWriter, *http.Request) 
)


func WebServer(port int) error {
        http.Handle("/", SetupUrlRouter())
        addr := fmt.Sprintf(":%d", port) 
        return http.ListenAndServe(addr, nil)
}

       
func SetupUrlRouter() *mux.Router {
        mappings := map[string]PageHandler {
                "/sensors": HandleSensors, 
        }
                
        router := mux.NewRouter()
        for url, hdlr := range(mappings) {
                router.HandleFunc(url, hdlr)
        }
        return router
}



func HandleSensors(w http.ResponseWriter, r *http.Request) {

        if buf, err := ioutil.ReadFile(SensorFile); err != nil {
                http.Error(w, err.Error(), http.StatusNotFound)
        } else {
                fmt.Fprint(w, string(buf))
        }
}
