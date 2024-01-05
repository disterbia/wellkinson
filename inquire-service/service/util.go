package service

import (
	"common/model"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func sendEmail(inquire model.Inquire, answer model.InquireReply) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	email := os.Getenv("WELLKINSON_SMTP_EMAIL")
	password := os.Getenv("WELLKINSON_SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", email, password, smtpHost)

	to := []string{inquire.Email}
	body := fmt.Sprintf(
		"<h2>작성자: </h2><span>%s</span><br>"+
			"<h2>날짜: </h2><span>%s</span><br>"+
			"<h2>제목: </h2><span>%s</span><br>"+
			"<h2>내용: </h2><span>%s</span><br>"+
			"<h2>답변: </h2><span>%s</span><br>"+
			"<h2>답변 날짜: </h2><span>%s</span><br>",
		inquire.Email, inquire.Created.Format("2006년 01월 02일 15시 04분 05초"),
		inquire.Title, inquire.Content, answer.Content,
		answer.Created.Format("2006년 01월 02일 15시 04분 05초"))

	msg := []byte("To: " + inquire.Email + "\r\n" +
		"Subject: 문의 주신 " + inquire.Title + "에 답변드립니다.\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + body)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, email, to, msg)
}
