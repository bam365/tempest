/* main.go - Main source and web server for tempest 
 * Copyright (C) 2013  Blake Mitchell 
 */
package main

import (
        "fmt"
        "encoding/json"
        "net/smtp"
        "bufio"
        "os"
        "time"
        "strings"
)

import (
        "github.com/howeyc/gopass"
)

import (
        "./config"
)


type TempestData struct {
        Conf *config.TempestConf
        Run *TempestRun
        EmailInf *EmailInfo
        Running bool
}


func main() {
        td := TempestData {}
        if conf, err := config.LoadConf("testconf"); err != nil {
                fmt.Printf("Error loading conf: %s\n", err)
        } else {
                td.Conf = &conf
                td.EmailInf = NewEmailInfoFromConf(conf.Smtp, GetEmailPassword(conf))
                td.Run = NewTempestRun("runhist.csv", conf)
                if td.Run.IsRunning() {
                        if rerr := td.Run.ResumeRun(); rerr != nil {
                                fmt.Printf("Run error: %s\n", rerr.Error())
                                return 
                        }
                }
                fmt.Print("Starting web server...")
                //TODO: WebServer() returns an err, which we shouldn't be ignoring.
                //This is going to take some doing.
                go WebServer(conf.Port)
                fmt.Println("DONE")
                td.Running = true
                RunConsole(&td)
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


func AlertHandler(amsg string, td *TempestData) { 
        tr, ei := td.Run, td.EmailInf
        sendmail := func(msg string) {
                err := smtp.SendMail(ei.FullServer(), ei.Auth, "Tempest Alerter", 
                tr.Conf.Emails, []byte("Subject: Tempest alert\n" + msg))
                if err != nil {
                        fmt.Printf("Could not send alert emails.\nReason: %s\n", 
                                    err.Error())
                }
        }
        alert := func(msg string) {
                layout := "1/2/06 3:04:05 PM"
                fmt.Printf("\n(%s) ALERT: %s\n", time.Now().Format(layout), msg)
                if tr.Conf.ShouldEmail() {
                        sendmail(msg)
                }
        }

        alert(amsg)
}


       
func RunConsole(td *TempestData) {
        cmdin := make(chan string)
        alerts := td.Run.alert
        prompt := func () {
                fmt.Print("tempest> ")
        }

        prompt()
        go GetCommand(cmdin)
        for td.Running {
                select {
                case amsg := <-alerts:
                        AlertHandler(amsg, td)
                        //Alert screwed up our prompt, redo it
                        prompt()
                case cmd := <-cmdin:
                        cmd = strings.ToLower(cmd)
                        RunCommand(cmd, td)
                        if (td.Running) {
                                prompt()
                                go GetCommand(cmdin)
                        }
                       
                }
        }
}


func GetCommand(cmdin chan<- string) {
        rdr := bufio.NewReader(os.Stdin)
        ret := ""
        //TODO: Don't ignore err here 
        if cmd, err := rdr.ReadString('\n'); err == nil {
                ret = strings.Trim(cmd, "\n")
        }
        cmdin <- ret
}


func RunCommand(cmd string, td *TempestData) int {
        //TODO: Will probably have to parse out args from cmd 
        ret := 0
        if c, exists := CommandMappings[cmd]; exists {
                ret = c.Run(td)
        } else {
                fmt.Println("No such command")
        }

        return ret
}
