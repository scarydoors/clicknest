package auth

import (
	kratos "github.com/ory/kratos-client-go"
)

var SessionCookieName = "ory_kratos_session"

func NewKratosClient(url string) *kratos.APIClient {
	configuration := kratos.NewConfiguration()
	configuration.Servers = []kratos.ServerConfiguration{
		{
			URL: url,
		},
	}

	apiClient := kratos.NewAPIClient(configuration);

	return apiClient
}
