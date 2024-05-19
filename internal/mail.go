package internal

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/mail"
	// "crypto/tls"
	"net/smtp"
	"net/url"
    "gopkg.in/gomail.v2"
	"strings"
)

type Destination struct {
	ToAddresses  []string `json:"ToAddresses"`
	CcAddresses  []string `json:"CcAddresses"`
	BccAddresses []string `json:"BccAddresses"`
}

type Content struct {
	Data    string `json:"Data"`
	CharSet string `json:"CharSet"`
}

type Body struct {
	Text Content `json:"Text"`
	Html Content `json:"Html"`
}

type Subject struct {
	Data string `json:"Data"`
}

type Message struct {
	Body    Body    `json:"Body"`
	Subject Subject `json:"Subject"`
}

type SendEmailRequest struct {
	Action           string      `json:"Action"`
	Destination      Destination `json:"Destination"`
	Message          Message     `json:"Message"`
	Source           string      `json:"Source"`
	ReplyToAddresses []string    `json:"ReplyToAddresses"`
}

func deserializeSendEmailRequest(reqBody string) (*SendEmailRequest, error) {
	queryValues, err := url.ParseQuery(reqBody)
	if err != nil {
		return nil, err
	}

	toAddresses := []string{queryValues.Get("Destination.ToAddresses.member.1")}

	// Then, initialize the struct fields using the map values
	sendEmailRequest := SendEmailRequest{
		Action: queryValues.Get("Action"),
		Destination: Destination{
			ToAddresses: toAddresses,
		},
		Message: Message{
			Body: Body{
				Html: Content{
					Data: queryValues.Get("Message.Body.Html.Data"),
				},
				Text: Content{
					Data: queryValues.Get("Message.Body.Text.Data"),
				},
			},
			Subject: Subject{
				Data: queryValues.Get("Message.Subject.Data"),
			},
		},
		Source: queryValues.Get("Source"),
	}

	for _, address := range toAddresses {
		if isEmailInvalid(address) {
			return nil, errors.New("To-Address is invalid: " + address)
		}
	}

	// Optional fields
	if ccAddresses, ok := queryValues["Destination.CcAddresses.member.1"]; ok {
		sendEmailRequest.Destination.CcAddresses = ccAddresses
		for _, address := range ccAddresses {
			if isEmailInvalid(address) {
				return nil, errors.New("CC-Address is invalid: " + address)
			}
		}
	}

	if bccAddresses, ok := queryValues["Destination.BccAddresses.member.1"]; ok {
		sendEmailRequest.Destination.BccAddresses = bccAddresses
		for _, address := range bccAddresses {
			if isEmailInvalid(address) {
				return nil, errors.New("BCC-Address is invalid: " + address)
			}
		}
	}

	if replyToAddresses, ok := queryValues["ReplyToAddresses.member.1"]; ok {
		sendEmailRequest.ReplyToAddresses = replyToAddresses
		for _, address := range replyToAddresses {
			if isEmailInvalid(address) {
				return nil, errors.New("Reply-To-Address is invalid: " + address)
			}
		}
	}

	return &sendEmailRequest, nil
}

func isEmailInvalid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err != nil
}

func SendEmail(bodyString string) error {
	request, err := deserializeSendEmailRequest(bodyString)

	if err != nil {
		return err
	}

	// Validation
	if !(request.Source != "" &&
		request.Message.Subject.Data != "" &&
		(request.Message.Body.Html.Data != "" || request.Message.Body.Text.Data != "") &&
		len(request.Destination.ToAddresses) > 0) {

		LogValidationErrors(request)

		return errors.New("one or more required fields was not sent")
	}

	return sendMail(request)

	// Mkdir dataDir and logDir
	// err = os.Mkdir(dataDir, 0755)
	// if err != nil && os.IsNotExist(err) {
	// 	return err
	// }

	// err = os.Mkdir(logDir, 0755)
	// if err != nil && os.IsNotExist(err) {
	// 	return err
	// }

	// // Write html data to dataDir/body.html
	// err = writeFileContent(filepath.Join(logDir, "body.html"), []byte(request.Message.Body.Html.Data))
	// if err != nil {
	// 	return err
	// }

	// // Write body to dataDir/body.txt
	// err = writeFileContent(filepath.Join(logDir, "body.txt"), []byte(request.Message.Body.Text.Data))
	// if err != nil {
	// 	return err
	// }

	// // Write headers to dataDir/headers.txt
	// headers := fmt.Sprintf("Subject: %s\nTo: %s\nCc: %s\nBcc: %s\nReply-To: %s\nFrom: %s\n",
	// 	request.Message.Subject.Data,
	// 	strings.Join(request.Destination.ToAddresses, ","),
	// 	strings.Join(request.Destination.CcAddresses, ","),
	// 	strings.Join(request.Destination.BccAddresses, ","),
	// 	strings.Join(request.ReplyToAddresses, ","),
	// 	request.Source,
	// )
	// err = writeFileContent(filepath.Join(logDir, "headers.txt"), []byte(headers))
	// if err != nil {
	// 	return err
	// }

	// Read file from templates/success.txt

	// return nil
}


func sendMail(req *SendEmailRequest) error {
	m := gomail.NewMessage()
    m.SetHeader("From", req.Source)
    m.SetHeader("To", strings.Join(req.Destination.ToAddresses, ","))
    m.SetHeader("Subject", req.Message.Subject.Data)

	if len(strings.TrimSpace(req.Message.Body.Html.Data)) > 0 {
		m.SetBody("text/html", req.Message.Body.Html.Data)
	} else {
		m.SetBody("text/plain", req.Message.Body.Text.Data)
	}

    // Send the email
    d := newDialer()
    return d.DialAndSend(m)
}

func newDialer() *gomail.Dialer {
	d := &gomail.Dialer{Host: Config.SmtpHost, Port: Config.SmtpPort}
	if Config.SmtpUser != "" && Config.SmtpPass != "" {
		d.Auth = smtp.CRAMMD5Auth(Config.SmtpUser, Config.SmtpPass)
	}
	// return gomail.NewPlainDialer(Config.SmtpHost, Config.SmtpPort, Config.SmtpUser, Config.SmtpPass)
	// d := gomail.NewDialer(Config.SmtpHost, Config.SmtpPort, Config.SmtpUser, Config.SmtpPass)
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d
}
// func sendMail(req *SendEmailRequest) error {
// 	from := req.Source
// 	to := strings.Join(req.Destination.ToAddresses, ",")
// 	subject := req.Message.Subject.Data

// 	msg := "From: " + from + "\n" +
// 		"To: " + to + "\n" +
// 		"Subject: " + subject + "\n\n" +
// 		"----- HTML -----" + 
// 		req.Message.Body.Html.Data +
// 		"----- TEXT -----" + 
// 		req.Message.Body.Text.Data

// 	return  smtp.SendMail(Config.SmtpHost, smtpAuth(),
// 		from, []string{to}, []byte(msg))
// }

// func smtpAuth() smtp.Auth {
// 	if Config.SmtpUser == "" || Config.SmtpPass == "" {
// 		return nil
// 	}
// 	return smtp.CRAMMD5Auth(Config.SmtpUser, Config.SmtpPass)
// 	// return smtp.PlainAuth("", Config.SmtpUser, Config.SmtpPass, Config.SmtpHost)
// }

func LogValidationErrors(request *SendEmailRequest) {
	// Check if ToAddresses is provided
	if len(request.Destination.ToAddresses) < 0 {
		logrus.Info("ToAddresses is not provided")
	}

	if request.Source == "" {
		logrus.Error("Source was not provided")
	}

	// Check if Subject is provided
	if request.Message.Subject.Data == "" {
		logrus.Error("Subject.Data was not provided")
	}

	// Check if Body.Html.Data or Body.Text.Data is provided
	if request.Message.Body.Html.Data == "" && request.Message.Body.Text.Data == "" {
		logrus.Error("Body.Html.Data or Body.Text.Data was not provided")
	}
}

func SendRawEmail(c *gin.Context, dateDir string, logFilePath string) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Not implemented",
	})
}
