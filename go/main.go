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
        Run *CurrentTempestRun
        EmailInf *EmailInfo
        Running bool
        //Email and stdout
        Alert chan string
        //Just stdout
        Msg chan string
}


func main() {
        td := new(TempestData)
        if conf, err := config.LoadConf("testconf"); err != nil {
                fmt.Printf("Error loading conf: %s\n", err)
        } else {
                fmt.Println("Loaded conf")
                td.Conf = &conf
                td.Alert = make(chan string)
                td.Msg = make(chan string)
                td.EmailInf = NewEmailInfoFromConf(conf.Smtp, GetEmailPassword(conf))
                TryResumeRun(td)
                fmt.Print("Starting web server...")
                ws := NewWebServer(td)
                //TODO: WebServer.Run() returns an err, which we shouldn't be ignoring.
                //This is going to take some doing.
                go ws.Run()
                fmt.Println("DONE")
                td.Running = true
                RunConsole(td)
        }
}


func JsonStr(v interface{}) string {
        b, _ := json.MarshalIndent(v, "", "\t")
        return string(b)
}


func TimeStr(t time.Time) string {
	return t.Format("1/2/2006 3:04:05pm")
}

func GetEmailPassword(conf config.TempestConf) string {
        ret := ""
        if conf.ShouldEmail() {
                fmt.Printf("Enter password for SMTP account %s at %s: ",
                            conf.Smtp.User, conf.Smtp.Server)
                ret = string(gopass.GetPasswd())
        }
        fmt.Println("\n*** I didn't try to authenticate with your password")
        fmt.Println("*** You should definitely run the 'testemail' command to verify\n")

        return ret
}


func TryResumeRun(td *TempestData) {
        if IsRunInProgress() {
                if r, rerr := ResumeCurrentTempestRun(td); rerr != nil {
                        fmt.Printf("Run error: %s\n", rerr.Error())
                } else {
                        td.Run = r
                        fmt.Print("Continuing run started on ")
                        fmt.Println(TimeStr(td.Run.TimeStarted()))
                }
        } else {
                fmt.Println("No run is in progress");
        }
}



//Will print a message to stdout if there's an error
func EmailEveryone(td TempestData, subject, body string) {
        ei, tc := td.EmailInf, td.Conf
        msg := fmt.Sprintf("Subject: %s\n%s\n", subject, body)
        err := smtp.SendMail(ei.FullServer(), ei.Auth, "Tempest Alerter", 
                             tc.Emails, []byte(msg))
        if err != nil {
                fmt.Printf("Could not send emails.\nReason %s\n", err.Error())
        }

}


func RunConsole(td *TempestData) {
        cmdin := make(chan string)
        prompt := func () {
                fmt.Print("tempest> ")
        }
        alert := func(amsg string) {
                layout := "1/2/06 3:04:05 PM"
                pmsg := fmt.Sprintf("\n(%s) ALERT: %s", 
                                    time.Now().Format(layout), amsg)
                fmt.Println(pmsg) 
                if td.Conf.ShouldEmail() {
                        EmailEveryone(*td, "Tempest Alert", amsg)
                }
        }

        prompt()
        go GetCommand(cmdin)
        for td.Running {
                select {
                case pmsg := <-td.Msg:
                        fmt.Println(pmsg)
                        //Message screwed up our prompt, redo it
                        prompt()
                case amsg := <-td.Alert:
                        alert(amsg)
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



