package internal

import (
	"context"
	"crypto/tls"
	"downloader_email/configs"
	"downloader_email/pkg"
	"downloader_email/rabbitmq"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	amqp "github.com/rabbitmq/amqp091-go"
	gomail "gopkg.in/gomail.v2"
)

type IEmailService interface {
}

type EmailService struct {
	rabbitmq rabbitmq.RabbitMQ
	configs  configs.ConfigStruct
	dialer   *gomail.Dialer
	dialLock *sync.RWMutex
	open     bool
}

const emailConsumerCount = 10

func NewEmailService(rabbit rabbitmq.RabbitMQ) *EmailService {
	conf := configs.GetConfigs()

	d := gomail.NewDialer(conf.MailServerHost, conf.MailServerPort, "", "")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} //not secure in production environment :: https://github.com/go-gomail/gomail

	emailSvc := EmailService{
		rabbitmq: rabbit,
		configs:  conf,
		dialer:   d,
		dialLock: &sync.RWMutex{},
		open:     true,
	}

	emailConfig := rabbitmq.NewConfigConsume(rabbitmq.EmailQueue, "")
	for i := 0; i < emailConsumerCount; i++ {
		ctx, _ := context.WithCancel(context.Background())
		go func() {
			openConChan := make(chan struct{})
			rabbitmq.NotifySetupDone(openConChan)
			<-openConChan
			if err := rabbit.Consume(ctx, emailConfig, &emailSvc, EmailConsumer); err != nil {
				message := fmt.Sprintf("error consuming from queue %s: %s", rabbitmq.EmailQueue, err)
				pkg.SaveError(message, err)
			}
		}()
	}

	return &emailSvc
}

//------------------------------------------
//------------------------------------------

type EmailType string

const (
	UserRegistration EmailType = "registration email"
	UserLogin        EmailType = "login email"
	PasswordUpdated  EmailType = "password updated"
	ResetPassword    EmailType = "reset password"
	VerifyEmail      EmailType = "verify email"
	DeleteAccount    EmailType = "delete account"
)

type EmailQueueData struct {
	Type        EmailType   `json:"type"`
	UserId      int64       `json:"userId"`
	RawUsername string      `json:"rawUsername"`
	Email       string      `json:"email"`
	Token       string      `json:"token"`
	Host        string      `json:"host"`
	Url         string      `json:"url"`
	DeviceInfo  *DeviceInfo `json:"deviceInfo"`
	IpLocation  string      `json:"ipLocation"`
}

type DeviceInfo struct {
	AppName     string `json:"appName" validate:"required"`
	AppVersion  string `json:"appVersion" validate:"required"` //format: ^\d\d?\.\d\d?\.\d\d?$
	Os          string `json:"os" validate:"required"`
	DeviceModel string `json:"deviceModel" validate:"required"`
	NotifToken  string `json:"notifToken"`
	Fingerprint string `json:"fingerprint"`
}

//------------------------------------------
//------------------------------------------

func EmailConsumer(d *amqp.Delivery, extraConsumerData interface{}) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		if os.Getenv("LOG_PANIC_TRACE") == "true" {
			log.Println(
				"level:", "error",
				"err: ", err,
				"trace", string(debug.Stack()),
			)
		} else {
			log.Println(
				"level", "error",
				"err", err,
			)
		}

		sentry.CurrentHub().Recover(err)
		sentry.Flush(time.Second * 5)

		//if err = d.Nack(false, true); err != nil {
		//	log.Printf("error nacking [email] message: %s\n", err)
		//}
		if err = d.Ack(false); err != nil {
			message := fmt.Sprintf("error acking [email] message: %s", err)
			pkg.SaveError(message, nil)
		}
	}()
	// run as rabbitmq consumer
	emailSvc := extraConsumerData.(*EmailService)
	var channelMessage *EmailQueueData
	err := json.Unmarshal(d.Body, &channelMessage)
	if err != nil {
		return
	}

	var m *gomail.Message
	switch channelMessage.Type {
	case UserRegistration:
		html := fmt.Sprintf("<a href=\"%v\">user %v, welcome to our wonderful app, click this link to verify your account<a/>",
			channelMessage.Url, channelMessage.RawUsername)
		m = createEmailMessage(
			emailSvc.configs.MailServerUsername,
			channelMessage.Email,
			"Thanks for registering",
			html,
		)
	case UserLogin:
		//todo : use config.userSessionsPage
		html := fmt.Sprintf("<div>\n <p>new device login: \n<p/>\n <p>appName: %v \n<p/>\n <p>appVersion: %v \n<p/>\n <p>deviceModel: %v \\<p/>\n <p>deviceModel: %v \n<p/>\n <p>ipLocation: %v \n<p/>\n </div>",
			channelMessage.DeviceInfo.AppName,
			channelMessage.DeviceInfo.AppVersion,
			channelMessage.DeviceInfo.Os,
			channelMessage.DeviceInfo.DeviceModel,
			channelMessage.IpLocation,
		)
		m = createEmailMessage(
			emailSvc.configs.MailServerUsername,
			channelMessage.Email,
			"New login detected",
			html,
		)
	case VerifyEmail:
		html := fmt.Sprintf("<a href=\"%v\">user %v, click this link to verify your account<a/>\n<p>if you don't know what's this, ignore.</p>",
			channelMessage.Url,
			channelMessage.RawUsername)
		m = createEmailMessage(
			emailSvc.configs.MailServerUsername,
			channelMessage.Email,
			"Verify Account",
			html,
		)
	case PasswordUpdated:
		//todo : use config.userSessionsPage
		html := fmt.Sprintf("<div>\n <p>Your password has been changed<p/>\n </div>")
		m = createEmailMessage(
			emailSvc.configs.MailServerUsername,
			channelMessage.Email,
			string(channelMessage.Type),
			html,
		)
	case DeleteAccount:
		// todo : delete account page design
		html := fmt.Sprintf("<a href=\"%v\">user %v, click this link to delete your account<a/>\n                        <p>if you don't know what's this, ignore.</p>",
			channelMessage.Url,
			channelMessage.RawUsername)
		m = createEmailMessage(
			emailSvc.configs.MailServerUsername,
			channelMessage.Email,
			"Delete Account",
			html,
		)
	case ResetPassword:
		//todo : implement
		return
	}

	err = emailSvc.dialer.DialAndSend(m)

	if err != nil {
		message := fmt.Sprintf("error on sending [%s] email: %s", channelMessage.Type, err)
		pkg.SaveError(message, err)
		if err = d.Nack(false, true); err != nil {
			message := fmt.Sprintf("error nacking [email] message: %s", err)
			pkg.SaveError(message, err)
		}
	} else {
		if err = d.Ack(false); err != nil {
			message := fmt.Sprintf("error acking [email] message: %s", err)
			pkg.SaveError(message, err)
		}
	}
}

func createEmailMessage(from string, to string, subject string, body string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return m
}

//------------------------------------------
//------------------------------------------
