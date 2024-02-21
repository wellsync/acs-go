package email

import "encoding/json"

type Message struct {
	Content       Content    `json:"content"`
	SenderAddress string     `json:"senderAddress"`
	Recipients    Recipients `json:"recipients"`
}

type Content struct {
	Html      string `json:"html"`
	PlainText string `json:"plainText"`
	Subject   string `json:"subject"`
}

type Recipients struct {
	To []Address `json:"to"`
}

type Address struct {
	Address     string `json:"address"`
	DisplayName string `json:"displayName"`
}

type EmailSendStatus string

const (
	EmailSendStatusCanceled   EmailSendStatus = "Canceled"
	EmailSendStatusFailed     EmailSendStatus = "Failed"
	EmailSendStatusNotStarted EmailSendStatus = "NotStarted"
	EmailSendStatusRunning    EmailSendStatus = "Running"
	EmailSendStatusSucceeded  EmailSendStatus = "Succeeded"
)

type ErrorDetail struct {
	AdditionalInfo []ErrorAdditionalInfo `json:"additionalInfo"`
	Code           string                `json:"code"`
	Details        []ErrorDetail         `json:"details"`
	Message        string                `json:"message"`
	Target         string                `json:"target"`
}

type ErrorAdditionalInfo struct {
	Info json.RawMessage `json:"info"`
	Type string          `json:"type"`
}

type SendResult struct {
	Error  ErrorDetail     `json:"error"`
	Id     string          `json:"id"`
	Status EmailSendStatus `json:"status"`
}
