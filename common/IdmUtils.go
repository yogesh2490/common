package common

import (
	"errors"
	"gopkg.in/resty.v0"
	"time"
	"utils/models"
)

/*
	This Utils is responsible for handling getting User info
*/

type IIdmUtils interface {
	GetAllUsers() (*[]models.UserData, error)
}

type IdmUtils struct {
	token  string
	config IConfigGetter
}

func GetIdmUtils(
	fullTokenString string,
	configs IConfigGetter,
) IIdmUtils {
	return &IdmUtils{
		fullTokenString,
		configs,
	}
}

func (this IdmUtils) GetAllUsers() (*[]models.UserData, error) {
	result := MakeAPIResult(this.config)
	defer result.Flush()

	data := []models.UserData{}
	serviceAddress := this.config.SafeGetConfigVar("IDM_URL") + "users/"

	resty.SetTimeout(time.Duration(10 * time.Second))

	resp, err := resty.R().
		SetHeader("Authorization", "Bearer "+this.token).
		SetHeader("Content-Type", "application/json").
		SetResult(&[]models.UserData{}).
		SetError(&models.NoUser{}).
		Get(serviceAddress)

	if err != nil {
		result.Errorf("operation error while getting user data from CRM endpoint: %s", err.Error())
		return &data, err
	}
	if resp.StatusCode() == 200 {
		result.Infof("successful user lookup")
		response := resp.Result().(*[]models.UserData)
		return response, nil
	} else if resp.StatusCode() == 404 {
		result.Errorf("user not found")
		return &data, errors.New("no user found")
	} else if resp.StatusCode() == 401 {
		result.Errorf("unauthorized")
		return &data, errors.New("unauthorized for CRM access")
	}

	result.Errorf("unknown error occured getting user data")

	return &data, errors.New("an unknown error has occured")
}
