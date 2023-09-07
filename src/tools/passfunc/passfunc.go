package passfunc

import "golang.org/x/crypto/bcrypt"

// BcryptCreatePassword создание хэша
func BcryptCreatePassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// BcryptCheckPassword проверка пароля
func BcryptCheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
