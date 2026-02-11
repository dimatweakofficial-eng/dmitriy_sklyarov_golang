package jwt

import "github.com/golang-jwt/jwt"

type JwtData struct {
	Email string
}

type JWT struct {
	Secret string
}

func NewJwt(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(data JwtData) (string, error) {
	//Замапили емаил внутрь jwt payload
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
	})
	//Вернули токен используя ключ
	token, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (j *JWT) Parse(token string) (bool, *JwtData) {
	//достаем из токена его payload через ключ
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return false, nil
	}
	//достаем значение по ключу
	email := t.Claims.(jwt.MapClaims)["email"]
	return t.Valid, &JwtData{
		Email: email.(string),
	}
}
