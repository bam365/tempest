/* web.go - Main source and web server for tempest 
 * Copyright (C) 2013  Blake Mitchell 
 */
package main

import (
        "fmt"
        "time"
        "encoding/json"
        //"net/http"
        "net/smtp"
)

import (
        //"github.com/gorilla/mux"
)

import (
        "./config"
)


func main() {
        if conf, err := config.LoadConf("testconf"); err != nil {
                fmt.Printf("Error loading conf: %s\n", err)
        } else {
                fmt.Println(JsonStr(conf.Sensors))
                //TODO: Don't do this here, start web server instead
                run := NewTempestRun("runhist.csv", conf)
                if rerr := run.StartRun(); rerr != nil {
                        fmt.Printf("Run error: %s", rerr.Error())
                } else {
                        defer run.StopRun()
                        time.Sleep(60 * time.Second)
                }
        }
}


func JsonStr(v interface{}) string {
        b, _ := json.MarshalIndent(v, "", "\t")
        return string(b)
}


func GetSmtpAuth(conf *conf.TempestConf) *smtp.Auth {
        return &smtp.PlainAuth(

func AlertListener(tr *TempestRun, smtpinf *smtp.ServerInfo) { 
        
        sendmail := func(msg string) {
                smtpauth := smtp.PlainAuth
                Semtp

        alert := func(msg string) {
                fmt.Println(msg)
                if smtpcli != nil && tr.Conf.ShouldEmail() {
                        sendmail(msg)
                }
        }

        for {
                amsg := <-

