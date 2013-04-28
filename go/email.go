/* email.go - Email functionality for tempest
 * Copyright (C) 2013  Blake Mitchell 
 */

package main


import (
        "net/smtp"
)


type EmailInfo struct {
        Serv string
        Auth smtp.Auth
}
        

func NewEmailInfo(addr, user, pass string) *EmailInfo {
        auth := smtp.PlainAuth("", user pass, addr)
        return &EmailInfo { addr, auth }
}
        
        
