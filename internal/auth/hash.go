package auth
import(
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) (string, error){
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ERROR generating hash: %v\n",err)
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
}