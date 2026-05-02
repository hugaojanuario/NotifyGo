package domain

type ChannelType string

const (
	Email   ChannelType = "EMAIL"
	Slack   ChannelType = "SLACK"
	Webhook ChannelType = "WEBHOOK"
	Sms     ChannelType = "SMS"
)
