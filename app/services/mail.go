package services

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"telego/app/config"
)

func SendWelcomeMail(to, username string) {
	if config.Config.AppEnv == "development" {
		return
	}

	m := mail.NewV3Mail()
	m.SetTemplateID("d-77046c231a30486a9909c52d8acfa9b6")
	p := mail.Personalization{
		To:                  []*mail.Email{mail.NewEmail(username, to)},
		From:                mail.NewEmail("Teleporthq", config.Config.Email.TeleportEmail),
		DynamicTemplateData: map[string]interface{}{"username": username},
	}
	m.AddPersonalizations(&p)

	sendMail(m)
}

func SendCollaborationInvitationMail(to, invitedToName, invitedByName, projectName string) {
	if config.Config.AppEnv == "development" {
		return
	}

	m := mail.NewV3Mail()
	m.SetTemplateID("d-ff8708b3beeb441aaa4b39883894d873")
	if invitedToName == "" {
		// we send a different template because it's a new user
		m.SetTemplateID("d-0d8885a72c8443f6a97a9474449c9e8f")
	}
	p := mail.Personalization{
		To:                  []*mail.Email{mail.NewEmail(invitedToName, to)},
		From:                mail.NewEmail("Teleporthq", config.Config.Email.TeleportEmail),
		DynamicTemplateData: map[string]interface{}{"invitedToName": invitedToName, "invitedByName": invitedByName, "projectName": projectName},
	}
	m.AddPersonalizations(&p)

	sendMail(m)
}

func sendMail(m *mail.SGMailV3) {
	request := sendgrid.GetRequest(config.Config.Email.SendgridApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)
	if err != nil {
		fmt.Println(err)
	}
}
