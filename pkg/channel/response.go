package channel

type Channel struct {
	Id    string `json:"id"`
	Index int    `json:"index"`
}

type CreateChannelResponse struct {
	Channels []*Channel `json:"channels"`
}

type RenewChannelResponse struct {
	Channels []*Channel `json:"channels"`
}

type ChannelListResponse struct {
	Channels []*Channel `json:"channels"`
}
