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
        "github.com/howeyc/gopass"
)

import (
        "./config"
)


func main() {
        if conf, err := config.LoadConf("testconf"); err != nil {
                fmt.Printf("Error loading conf: %s\n", err)
        } else {
                emailinf := NewEmailInfoFromConf(conf.Smtp, GetEmailPassword(conf))
                fmt.Println(JsonStr(conf.Sensors))
                //TODO: Don't do this here, start web server instead
                run := NewTempestRun("runhist.csv", conf)
                if rerr := run.StartRun(); rerr != nil {
                        fmt.Printf("Run error: %s\n", rerr.Error())
                } else {
                        defer run.StopRun()
                        go AlertListener(run, emailinf)
                        time.Sleep(60 * time.Second)
                }
        }
}


func JsonStr(v interface{}) string {
        b, _ := json.MarshalIndent(v, "", "\t")
        return string(b)
}


func GetEmailPassword(conf config.TempestConf) string {
        ret := ""
        if conf.ShouldEmail() {
                fmt.Printf("Enter password for SMTP account %s at %s: ",
                            conf.Smtp.User, conf.Smtp.Server)
                ret = string(gopass.GetPasswd())
        }
        return ret
}


func AlertListener(tr *TempestRun, ei *EmailInfo) { 
        sendmail := func(msg string) {
                err := smtp.SendMail(ei.FullServer(), ei.Auth, "Tempest Alerter", 
                tr.Conf.Emails, []byte("Subject: Tempest alert\n" + msg))
                if err != nil {
                        fmt.Printf("Could not send alert emails.\nReason: %s\n", 
                                    err.Error())
                }
        }
        alert := func(msg string) {
                fmt.Println(msg)
                if tr.Conf.ShouldEmail() {
                        sendmail(msg)
                }
        }

        //DBG
        fmt.Println(tr.Conf.Emails)
        for {
                amsg := <-tr.alert
                alert(amsg)
        }
}


