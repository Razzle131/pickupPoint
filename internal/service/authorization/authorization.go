package authorization

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Razzle131/pickupPoint/api"
	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/Razzle131/pickupPoint/internal/repository/userRepo"
	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

const (
	jwtExpire = 10 * time.Minute
)

type jwtClaims struct {
	Role api.UserRole
	jwt.StandardClaims
}

type AuthorizationService struct {
	userRepo userRepo.UserRepo
}

func New(ur userRepo.UserRepo) *AuthorizationService {
	return &AuthorizationService{
		userRepo: ur,
	}
}

func (s *AuthorizationService) DummyLogin(role api.UserRole) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtExpire).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	return token.SignedString(jwtSecret)
}

func (s *AuthorizationService) Login(req dto.LoginDto) (string, error) {
	user, err := s.userRepo.GetUserByEmail(context.TODO(), req.Email)
	if err != nil {
		return "", err
	}

	if user.Password != req.Password {
		return "", errors.New("bad creditionals")
	}

	return s.DummyLogin(user.Role)
}

// TODO: mb change internal model to dto in return
func (s *AuthorizationService) Register(req dto.UserDto) (model.User, error) {
	_, err := s.userRepo.GetUserByEmail(context.TODO(), req.Email)
	if err == nil {
		return model.User{}, errors.New("user already exists")
	}

	user, err := s.userRepo.AddUser(context.TODO(), req)
	if err != nil {
		return model.User{}, errors.New("failed to add user")
	}

	return user, nil
}

func (s *AuthorizationService) ValidateToken(token api.Token) (api.UserRole, error) {
	splittedToken := strings.Split(token, " ")
	if len(splittedToken) < 2 {
		slog.Debug("bad token")
		return "", errors.New("bad token")
	}

	tmp := jwtClaims{}
	parsedToken, _ := jwt.ParseWithClaims(splittedToken[1], &tmp, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if !parsedToken.Valid {
		slog.Debug("bad token")
		return "", errors.New("bad token")
	}

	return tmp.Role, nil
}

// func (s *AuthorizationService) AuthenticateUser(login, password string) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
// 		login,
// 		jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(jwtExpire).Unix(),
// 			IssuedAt:  time.Now().Unix(),
// 		},
// 	})

// 	user, err := s.userRepo.GetUserByLogin(login)
// 	if err != nil && err == serverErrors.ErrUserNotFound {
// 		err = s.userRepo.AddUser(login, password)
// 		if err != nil {
// 			return "", serverErrors.ErrInternal
// 		}

// 		return token.SignedString(jwtSecret)
// 	} else if err != nil {
// 		return "", serverErrors.ErrInternal
// 	}

// 	if user.Password != password {
// 		return "", serverErrors.ErrBadCreditonals
// 	}

// 	return token.SignedString(jwtSecret)
// }

// func (s *AuthorizationService) AuthorizeUser(token string) (model.User, error) {
// 	splittedToken := strings.Split(token, " ")
// 	if len(splittedToken) < 2 {
// 		slog.Debug("bad token")
// 		return model.User{}, serverErrors.ErrBadToken
// 	}

// 	tmp := jwtClaims{}
// 	parsedToken, _ := jwt.ParseWithClaims(splittedToken[1], &tmp, func(token *jwt.Token) (interface{}, error) {
// 		return jwtSecret, nil
// 	})

// 	slog.Debug("AuthorizeUser: " + fmt.Sprint(tmp.Login))

// 	if parsedToken.Valid {
// 		slog.Debug("token valid")
// 		return s.userRepo.GetUserByLogin(tmp.Login)
// 	}

// 	return model.User{}, serverErrors.ErrBadToken
// }
