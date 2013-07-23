/* email.go - Email functionality for tempest
 * Copyright (C) 2013  Blake Mitchell 
 */

package main


import (
        "net/smtp"
        "fmt"
)

import (
        "./config"
)


type EmailInfo struct {
        Serv string
        Port int
        Auth smtp.Auth
}
        

func NewEmailInfo(addr string, port int, user, pass string) *EmailInfo {
        auth := smtp.PlainAuth("", user, pass, addr)
        return &EmailInfo { addr, port, auth }
}
        

func NewEmailInfoFromConf(ss config.SmtpSettings, pass string) *EmailInfo {
        return NewEmailInfo(ss.Server, ss.Port, ss.User, pass)
}
        

func (ei *EmailInfo) FullServer() string {
        return fmt.Sprintf("%s:%d", ei.Serv, ei.Port)
}




