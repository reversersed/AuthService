package smtp

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	mock_smtp "github.com/reversersed/AuthService/internal/smtp/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSendEmailWarning(t *testing.T) {
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       true,
		LogServerActivity: true,
	})

	err := server.Start()
	assert.NoError(t, err)

	defer server.Stop()

	ctrl := gomock.NewController(t)
	logger := mock_smtp.NewMocklogger(ctrl)
	cfg := &SmtpConfig{
		Host:     "127.0.0.1",
		Port:     server.PortNumber(),
		User:     "server@example.com",
		Password: "",
	}
	logger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	service := New(cfg, logger)

	service.SendEmailWarning("ip address")

	time.Sleep(500 * time.Millisecond)

	msg := server.Messages()
	assert.True(t, msg[len(msg)-1].Mailfrom())
	assert.Equal(t, "From: server@example.com \r\nTo: server@example.com \r\nSubject: Выполнен вход в аккаунт \r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n\r\n\r\n<html>\r\n\t<head> </head>\r\n\t<body>\r\n\t\t<table\r\n\t\t\tstyle=\"\r\n\t\t\t\twidth: fit-content;\r\n\t\t\t\tbackground-color: navy;\r\n\t\t\t\tborder-radius: 10px;\r\n\t\t\t\tcolor: white;\r\n\t\t\t\tfont-family: sans-serif;\r\n\t\t\t\ttext-align: center;\r\n\t\t\t\tmax-width: 500px;\r\n\t\t\t\tborder-spacing: 30px;\r\n\t\t\t\"\r\n\t\t>\r\n\t\t\t<tr>\r\n\t\t\t\t<td>\r\n\t\t\t\t\t<h2>Здравствуйте, client name!</h2>\r\n\t\t\t\t</td>\r\n\t\t\t</tr>\r\n\t\t\t<tr>\r\n\t\t\t\t<td>\r\n\t\t\t\t\tУведомляем Вас, что в Ваш аккаунт был выполнен вход с IP ip address.\r\n\t\t\t\t</td>\r\n\t\t\t</tr>\r\n\t\t\t<tr>\r\n\t\t\t\t<td>\r\n\t\t\t\t\tЕсли это были не Вы, настоятельно рекомендуем изменить пароль к\r\n\t\t\t\t\tаккаунту и проверить активные сессии.\r\n\t\t\t\t</td>\r\n\t\t\t</tr>\r\n\t\t</table>\r\n\t</body>\r\n</html>\r\n", msg[len(msg)-1].MsgRequest())
}
