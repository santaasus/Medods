package auth

import (
	"errors"
	"fmt"

	jwtHandler "github.com/santaasus/JWTToken-handler"
	errorsDomain "github.com/santaasus/errors-handler"

	domain "Medods/auth_service/inner_layer/domain"
	repository "Medods/auth_service/inner_layer/repository/user"

	"Medods/auth_service/inner_layer/security"
)

type Service struct {
	UserRepository repository.IRepository
}

func (s *Service) GetNewTokens(guid string, clientIP string) (*SecurityData, error) {
	if guid == "" {
		err := &errorsDomain.AppError{
			Err:  errors.New("guid is empty"),
			Type: errorsDomain.NotFound,
		}
		return nil, err
	}

	accessToken, refreshToken, err := generateTokens(guid, map[string]any{"ip": clientIP})
	if err != nil {
		return nil, err
	}

	hash, err := security.GeneratePasswordHash(refreshToken.Token[:72])
	if err != nil {
		return nil, err
	}

	domainUser, _ := s.UserRepository.GetUserByGuid(guid)
	if domainUser != nil {
		params := map[string]any{
			"refresh_hash": string(hash),
			"ip":           clientIP,
		}
		s.UserRepository.UpdateUser(params, domainUser.ID)
	} else {
		newUser := &domain.NewUser{Guid: guid, Hash: string(hash), IP: clientIP}
		s.UserRepository.CreateUser(newUser)
	}

	return &SecurityData{
		JWTAccessToken:            accessToken.Token,
		JWTRefreshToken:           refreshToken.Token,
		ExpirationAccessDateTime:  accessToken.ExpirationTime,
		ExpirationRefreshDateTime: refreshToken.ExpirationTime,
	}, nil
}

func (s *Service) AccessTokenByRefreshToken(refresh string, clientIP string) (*SecurityData, error) {
	claims, err := jwtHandler.VerifyTokenAndGetClaims(refresh, "refresh")
	if err != nil {
		if err.Error() == "token expired" {
			err := s.UserRepository.DeleteUserByHash(refresh[:72])
			fmt.Print(err)
		}

		return nil, err
	}

	guid, ok := claims["id"].(string)
	if !ok {
		return nil, &errorsDomain.AppError{
			Err:  errors.New("guid doesn't match: claims[guid].(string)"),
			Type: errorsDomain.ValidationError,
		}
	}

	accessToken, refreshToken, err := generateTokens(guid, nil)
	if err != nil {
		return nil, err
	}

	var warning string
	domainUser, _ := s.UserRepository.GetUserByGuid(guid)
	if domainUser != nil && domainUser.IP != clientIP {
		warning = "Seems your ip address was changed. If you didn't visit pornhub.com right now, open this instruction: lalala.ucoz/safe_this_world's_eyes_from_your_bdsm_content"
	}

	return &SecurityData{
		Warning:                   warning,
		JWTAccessToken:            accessToken.Token,
		JWTRefreshToken:           refreshToken.Token,
		ExpirationAccessDateTime:  accessToken.ExpirationTime,
		ExpirationRefreshDateTime: refreshToken.ExpirationTime,
	}, nil
}

func generateTokens(guid string, payload map[string]any) (accessToken *jwtHandler.AppToken, refreshToken *jwtHandler.AppToken, err error) {
	accessToken, err = jwtHandler.GenerateJWTToken(guid, jwtHandler.Access, nil)
	if err != nil {
		return
	}

	refreshToken, err = jwtHandler.GenerateJWTToken(guid, jwtHandler.Refresh, payload)
	if err != nil {
		return
	}

	return
}
