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


