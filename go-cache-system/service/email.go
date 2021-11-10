/**
 * @Author: djh
 * @Description:
 * @File:  email
 * @Version: 1.0.0
 * @Date: 2021/10/18 19:46
 */

package service

import (
	"cache-system/config"
	"fmt"
	"log"
	"net/smtp"
)

func sendSMTPMail(body string, emailChan chan byte) {
	conf := config.Configs.Email
	// 通常身份应该是空字符串，填充用户名.
	auth := smtp.PlainAuth("", conf.User, conf.Password, conf.Host)
	// (text/plain)纯文本 , (text/html)html
	contentType := "Content-Type: text/html; charset=UTF-8"

	for _, mailAddress := range conf.To {
		msg := []byte("To: " + mailAddress + "\r\nFrom: " + conf.Subject + "<" + conf.User + ">\r\nSubject: " + conf.Subject +
			"\r\n" + contentType + "\r\n\r\n" + body)
		err := smtp.SendMail(fmt.Sprintf("%s:%d", conf.Host, conf.Port), auth, conf.User, conf.To, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
	emailChan <- 1
}

func SendSMTPMail(body string) {
	if config.Configs.Select.Env != "dev" {
		conf := config.Configs.Email
		// 通常身份应该是空字符串，填充用户名.
		auth := smtp.PlainAuth("", conf.User, conf.Password, conf.Host)
		// (text/plain)纯文本 , (text/html)html
		contentType := "Content-Type: text/html; charset=UTF-8"

		for _, mailAddress := range conf.To {
			msg := []byte("To: " + mailAddress + "\r\nFrom: " + conf.Subject + "<" + conf.User + ">\r\nSubject: " + conf.Subject +
				"\r\n" + contentType + "\r\n\r\n" + body)
			err := smtp.SendMail(fmt.Sprintf("%s:%d", conf.Host, conf.Port), auth, conf.User, conf.To, msg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
