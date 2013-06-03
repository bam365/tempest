/* main.go - Main source and web server for tempest 
 * Copyright (C) 2013  Blake Mitchell 
 */
package main

import (
        "fmt"
        "encoding/json"
        "time"
        "net/smtp"
        "bufio"
        "os"
        "strings"
)

import (
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
                //TODO: Don't do this here, start web server instead
                run := NewTempestRun("runhist.csv", conf)
                if rerr := run.StartRun(); rerr != nil {
                        fmt.Printf("Run error: %s\n", rerr.Error())
                } else {
                        defer run.StopRun()
                        go AlertListener(run, emailinf)
                        time.Sleep(60 * time.Second)
                }

                fmt.Print("Starting web server...")
                //TODO: WebServer() returns an err, which we shouldn't be ignoring.
                //This is going to take some doing.
                go WebServer(conf.Port)
                fmt.Println("DONE")
                RunConsole(&conf, run.alert)
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

       
type Command func(*config.TempestConf) int


var CommandMappings = map[string]Command {
        "conf": CmdDumpConf,
}


func RunConsole(conf *config.TempestConf, alerts <-chan string) {
        quit := false
        cmdin := make(chan string)
        prompt := func () {
                fmt.Print("tempest> ")
        }

        prompt()
        go GetCommand(cmdin)
        for !quit {
                select {
                case <-alerts:
                        //Alert screwed up our prompt, redo it
                        prompt()
                case cmd := <-cmdin:
                        cmd = strings.ToLower(cmd)
                        if cmd == "quit" || cmd == "bye" || cmd == "exit" {
                                quit = true
                        } else {
                                RunCommand(cmd, conf)
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


func RunCommand(cmd string, conf *config.TempestConf) int {
        //TODO: Will probably have to parse out args from cmd 
        ret := 0
        if c, exists := CommandMappings[cmd]; exists {
                ret = c(conf)
        } else {
                fmt.Println("No such command")
        }

        return ret
}


func CmdDumpConf(conf *config.TempestConf) int {
        fmt.Println(JsonStr(conf))
        return 0
}




