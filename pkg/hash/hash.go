package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash tạo ra một chuỗi hash bcrypt từ mật khẩu.
// Hàm này được thiết kế để hash mật khẩu.
func Hash(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(bytes), err
}

// Verify so sánh một mật khẩu đã được hash bằng bcrypt với phiên bản gốc của nó.
// Trả về true nếu mật khẩu và hash khớp nhau.
func Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
