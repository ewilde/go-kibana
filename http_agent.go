package kibana

import "github.com/parnurzeal/gorequest"

type HttpAgent struct {
	client      *gorequest.SuperAgent
	authHandler AuthenticationHandler
}

type AuthenticationHandler interface {
	Initialize(agent *gorequest.SuperAgent)
}

type NoAuthenticationHandler struct {
}

type BasicAuthenticationHandler struct {
	userName string
	password string
}

type LogzAuthentication struct {
}

func NewHttpAgent(config *Config, authHandler AuthenticationHandler) *HttpAgent {
	return &HttpAgent{
		client:      gorequest.New(),
		authHandler: authHandler,
	}
}

func (authClient *HttpAgent) Auth(handler AuthenticationHandler) *HttpAgent {
	authClient.authHandler = handler
	return authClient
}

func (authClient *HttpAgent) Get(targetUrl string) *HttpAgent {
	authClient.client.Get(targetUrl)
	return authClient
}

func (authClient *HttpAgent) Delete(targetUrl string) *HttpAgent {
	authClient.client.Delete(targetUrl)
	return authClient
}

func (authClient *HttpAgent) Put(targetUrl string) *HttpAgent {
	authClient.client.Put(targetUrl)
	return authClient
}

func (authClient *HttpAgent) Post(targetUrl string) *HttpAgent {
	authClient.client.Post(targetUrl)
	return authClient
}

func (authClient *HttpAgent) Query(content interface{}) *HttpAgent {
	authClient.client.Query(content)
	return authClient
}

func (authClient *HttpAgent) Set(param string, value string) *HttpAgent {
	authClient.client.Set(param, value)
	return authClient
}

func (authClient *HttpAgent) Send(content interface{}) *HttpAgent {
	authClient.client.Send(content)
	return authClient
}

func (authClient *HttpAgent) End(callback ...func(response gorequest.Response, body string, errs []error)) (gorequest.Response, string, []error) {
	authClient.authHandler.Initialize(authClient.client)
	return authClient.client.End(callback...)
}

func NewBasicAuthentication(userName string, password string) *BasicAuthenticationHandler {
	return &BasicAuthenticationHandler{userName: userName, password: password}
}

func (auth *BasicAuthenticationHandler) Initialize(agent *gorequest.SuperAgent) {
	agent.SetBasicAuth(auth.userName, auth.password)
}

func (auth *NoAuthenticationHandler) Initialize(agent *gorequest.SuperAgent) {
}
