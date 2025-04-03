package sandbox

type CodeServerConfig struct {
	BindAddr string `yaml:"bind-addr"`
	Auth     string `yaml:"auth"`
	Password string `yaml:"password"`
	// HashedPassword []byte
	Cert bool `yaml:"cert"`
}

// const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// // Function to generate a random password
// func generateRandomPassword(length int) string {
// 	var sb strings.Builder
// 	for i := 0; i < length; i++ {
// 		sb.WriteByte(charset[rand.Intn(len(charset))])
// 	}
// 	return sb.String()
// }

// func getPasswordHash(password string) ([]byte, error) {
// 	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return []byte{}, err
// 	}
// 	return hash, nil
// }

func GetCodeServerConfig() (*CodeServerConfig, error) {
	// password := generateRandomPassword(10)
	// hashedPassword, err := getPasswordHash(password)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not get hashed password")
	// }
	return &CodeServerConfig{
		BindAddr: "0.0.0.0:9090",
		Auth:     "none",
		Cert:     false,
		// Password: password,
		// HashedPassword: hashedPassword,
	}, nil
}
