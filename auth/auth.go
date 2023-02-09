package auth

import (
	"encoding/json"
	"fmt"
	sqlconnection "server/db/mysql"
)

func AddNewUser(raw []byte) error {
	// parse raw
	var userData *SignUp
	fmt.Println(string(raw))
	err := json.Unmarshal(raw, &userData)
	fmt.Println("unmarshalled")
	if err != nil {
		return fmt.Errorf("error while unmarshaling user signup data: %v", err)
	}
	fmt.Println("parsed raw")
	fmt.Println(userData)
	fmt.Println(&userData)
	fmt.Println(*userData)
	err = sqlconnection.SearchUserInUsers(struct {
		Name     string
		Email    string
		Password string
		Phone    string
	}(*userData))
	fmt.Println("exiting AddNewUser")
	if err != nil {
		return fmt.Errorf("error while adding user to db: %v", err)
	}
	return nil
}
