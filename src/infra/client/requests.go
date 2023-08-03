package client

type PostRequest struct {
	Path  string
	Body  []byte
	Token string
	// todo: header
}

func (p PostRequest) GetPath() string {
	return p.Path
}

func (p PostRequest) GetBody() []byte {
	return p.Body
}

type GetRequest struct {
	Path string
}

func (g GetRequest) GetBody() []byte {
	//TODO implement me
	panic("implement me")
}

func (g GetRequest) GetPath() string {
	return g.Path
}
