package main

import (
	"html/template"
	"bytes"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"

)

// Mailer Config
type Mail struct {
	Domain 			string
	Host				string
	Port				int
	Username		string
	Password		string
	Encryption	string
	FromAddress	string
	FromName		string
}

// indvidual Message
type Message struct {
	From				string
	FromName		string
	To					string
	Subject			string
	Attachments	[]string
	Data				any
	DataMap			map[string]any
}

// Send the email function
func (m *Mail) SendSMTPMessage(msg Message) error {

	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.Name == "" {
		msg.FromName = m.FromName
	}

	// HTML MAIL 
	data := map[string]any {
		"message": msg.Data,
	}

	msg.DataMap = data

	// HTML Text Message
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return "",err
	}
	// Plain Text Message
	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return "",err
	}

	// Create a Mail server
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return "",nil
	}

	// Now send an email
	email := mail.NewMSG()
	email.SetFrom(msg.From).
				AddTo(msg.To).
				SetSubject(msg.Subject)
	
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML,formattedMessage)

	// Add Attachments
	if len(msg.Attachments) > 0 {
		for _,x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	// Send the file
	err = email.Send(smtpClient)
	if err != nil {
		return "",err
	}

	// Success in email generation
	return nil

}

func (m *Mail) buildHTMLMessage(msg Message) (string error) {
	templateToRender := "./templates/mail.html.gohtml"
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "",err
	}

	// Error in the template file
	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "",err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "",nil
	}

	return formattedMessage, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string error) {
	templateToRender := "./templates/mail.plain.gohtml"
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}
	
	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl,"body",msg.Data); err != nil {
		return "", nil
	}

	plainMessage := tpl.String()

	return plainMessage,nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses: false,
		CssToAttributes: false,
		KeepBangImportant: true
	}
	prem, err := premailer.NewPremailerFromString(s,&options)
	if err != nil {
		return "",err
	}

	html, err := prem.Transform()
	if err != nil {
		return "",err
	}

	return html,nil
}

func (m *Mail) getEncryption(s string) mail.Encryption {

	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none","":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}

}