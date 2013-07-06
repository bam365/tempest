/* command.go - Tempest console commands
 * Copyright (C) 2013  Blake Mitchell 
 */
package main

import (
        "fmt"
)

import (
        "./sensors"
)


const (
        RunSuccess = iota
        RunFailed
)


type Commander interface {
        Describe() string
        Run(*TempestData) int
}
         

var CommandMappings = map[string]Commander {
        "help": CmdHelp {},
        "?": CmdHelp {},
        "bye": CmdQuit {},
        "quit": CmdQuit {},
        "conf": CmdConfDump {},
        "read": CmdReadSensors {},
        "testemail": CmdTestEmail {},
        "startrun": CmdStartRun {},
        "stoprun": CmdStopRun {},
        "runstat": CmdRunStat {},
}


type CmdHelp struct {}

func (s CmdHelp) Describe() string {
        return "List commands"
}

func (s CmdHelp) Run(td *TempestData) int {
        for cmdstr, cmd := range(CommandMappings) {
                fmt.Printf("%10s - %s\n", cmdstr, cmd.Describe())
        }
        return RunSuccess
}


type CmdQuit struct {}

func (s CmdQuit) Describe() string {
        return "Quit tempest"
}

func (s CmdQuit) Run(td *TempestData) int {
        //TODO: Do other shutdown stuff?
        td.Running = false
        return RunSuccess
}


type CmdConfDump struct {}

func (s CmdConfDump) Describe() string {
        return "Dump config object to console"
}

func (s CmdConfDump) Run(td *TempestData) int {
        fmt.Println(JsonStr(td.Conf))
        return RunSuccess
}


type CmdReadSensors struct {}

func (s CmdReadSensors) Describe() string {
        return "Give current sensor readings"
}

func (s CmdReadSensors) Run(td *TempestData) int {
        for sname, sdat := range(sensors.ReadSensors(td.Conf.Sensors)) {
                fmt.Printf("%s: %s\n", sname, JsonStr(sdat))
        }
        return RunSuccess
}


type CmdTestEmail struct {}

func (s CmdTestEmail) Describe() string {
        return "Send a test email to everyone"
}

func (s CmdTestEmail) Run(td *TempestData) int {
        fmt.Println("Sending email to:")
        for _, emailaddr := range(td.Conf.Emails) {
                fmt.Printf("\t%s\n", emailaddr)
        }
        EmailEveryone(*td, "Tempest Test", "Tempest test")
        return RunSuccess
}


type CmdStartRun struct {}

func (s CmdStartRun) Describe() string {
        return "Start new run (if there isn't one already)"
}

func (s CmdStartRun) Run(td *TempestData) int {
        ret := RunSuccess
        if !IsRunInProgress() {
                fmt.Print("Starting run...")
                if r, rerr := StartNewTempestRun(td); rerr != nil {
                        fmt.Println("FAILED")
                        fmt.Printf("Reason: %s\n", rerr.Error())
                        ret = RunFailed
                } else {
                        td.Run = r 
                        fmt.Println("DONE")
                }
        } else {
                fmt.Println("There's already a run in progress")
        }
        return ret
}


type CmdStopRun struct {}

func (s CmdStopRun) Describe() string {
        return "Stop current run (if there is one)"
}

func (s CmdStopRun) Run(td *TempestData) int {
        if IsRunInProgress() {
                fmt.Print("Stopping run...")
                td.Run.StopRun()
                fmt.Println("DONE")
        } else {
                fmt.Println("No run in progress")
        }
        return RunSuccess
}


type CmdRunStat struct {}

func (s CmdRunStat) Describe() string {
        return "Give information on current run"
}

func (s CmdRunStat) Run(td *TempestData) int {
        if IsRunInProgress() {
                run := td.Run
                fmt.Printf("Run started: %s\nRun duration: %s\n",
                           run.TimeStarted(), run.RunDuration())
        } else {
                fmt.Println("No run in progress")
        }
        return RunSuccess
}


