package authorization

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/repository/userRepo"
	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

const (
	jwtExpire = 1 * time.Hour
)

type jwtClaims struct {
	Role string
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

func (s *AuthorizationService) DummyLogin(ctx context.Context, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtExpire).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	return token.SignedString(jwtSecret)
}

func (s *AuthorizationService) Login(ctx context.Context, req dto.UserCreditionalsDto) (string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	if user.Password != req.Password {
		return "", errors.New("bad creditionals")
	}

	return s.DummyLogin(ctx, user.Role)
}

func (s *AuthorizationService) Register(ctx context.Context, req dto.UserCreditionalsDto) (dto.UserDto, error) {
	_, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return dto.UserDto{}, errors.New("user already exists")
	}

	user := req.ToModel()

	newUser, err := s.userRepo.AddUser(ctx, user)
	if err != nil {
		return dto.UserDto{}, errors.New("failed to add user")
	}

	res := dto.UserDto{}
	res.FromModel(newUser)

	return res, nil
}

func (s *AuthorizationService) ValidateToken(ctx context.Context, token string) (string, error) {
	// made for matching tokens with and without bearer prefix
	tmpToken := token
	if match, _ := regexp.Match("bearer\\s.*", []byte(strings.ToLower(tmpToken))); match {
		tmpToken = strings.Split(token, " ")[1]
	}

	tmp := jwtClaims{}
	parsedToken, err := jwt.ParseWithClaims(tmpToken, &tmp, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !parsedToken.Valid {
		slog.Error("failed token validation")
		return "", errors.New("bad token")
	}

	return tmp.Role, nil
}
