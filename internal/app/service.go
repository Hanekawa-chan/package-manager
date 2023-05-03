package app

type service struct {
	client Client
}

func New(client Client) Service {
	return &service{
		client: client,
	}
}
