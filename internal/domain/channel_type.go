package domain

type ChannelType string

const (
	ChannelEmail   ChannelType = "EMAIL"
	ChannelSlack   ChannelType = "SLACK"
	ChannelWebhook ChannelType = "WEBHOOK"
	ChannelSMS     ChannelType = "SMS"
)
