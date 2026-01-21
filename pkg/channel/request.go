package channel

type AcquireChannelRequest struct {
	Id string `json:"id"`
}

type RenewChannelRequest struct {
	Id string `json:"id"`
}
