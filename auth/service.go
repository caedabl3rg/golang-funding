package auth

import "github.com/dgrijalva/jwt-go"

type Service interface {
	GenerateToken(userID int) (string, error)
}

type jwtService struct{}

var SCRET_KEY = []byte("hanyabolehdiketahuiolehkita")

func NewService() *jwtService{
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userID int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(SCRET_KEY)
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}
