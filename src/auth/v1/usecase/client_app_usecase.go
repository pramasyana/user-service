package usecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
)

// CreateClientApp function
func (au *AuthUseCaseImpl) CreateClientApp(name string) <-chan ResultUseCase {
	ctx := "AuthUseCase-CreateClientApp"

	output := make(chan ResultUseCase)
	go func() {
		defer close(output)

		if len(name) == 0 {
			output <- ResultUseCase{Error: errors.New("name cannot be empty")}
			return
		}

		clientApp := model.NewClientApp(name)

		clientAppResult := <-au.ClientAppRepoWrite.Save(clientApp)

		if clientAppResult.Error != nil {
			helper.SendErrorLog(context.Background(), ctx, "create_client_app", clientAppResult.Error, name)
			output <- ResultUseCase{Error: clientAppResult.Error, HTTPStatus: http.StatusUnauthorized}
			return
		}

		output <- ResultUseCase{Result: clientApp}

	}()

	return output
}

// GetClientApp function for getting and validating client app
func (au *AuthUseCaseImpl) GetClientApp(clientID, clientSecret string) <-chan ResultUseCase {

	output := make(chan ResultUseCase)
	go func() {
		defer close(output)

		clientAppResult := <-au.ClientAppRepoRead.FindByClientID(clientID)

		if clientAppResult.Error != nil {
			output <- ResultUseCase{Error: clientAppResult.Error, HTTPStatus: http.StatusUnauthorized}
			return
		}

		clientApp, ok := clientAppResult.Result.(model.ClientApp)

		if !ok {
			output <- ResultUseCase{Error: errors.New("result is not basic auth"), HTTPStatus: http.StatusUnauthorized}
			return
		}

		// validate client app secret
		valid := clientApp.Authenticate(clientSecret)
		if !valid {
			output <- ResultUseCase{Error: errors.New("password does not match"), HTTPStatus: http.StatusUnauthorized}
			return
		}
		if !clientApp.IsActive() {
			output <- ResultUseCase{Error: errors.New("client status not active"), HTTPStatus: http.StatusUnauthorized}
			return
		}

		output <- ResultUseCase{Result: valid}

	}()

	return output
}
