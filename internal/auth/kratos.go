package auth

import (
	kratos "github.com/ory/kratos-client-go"
)

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
