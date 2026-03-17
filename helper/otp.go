package helper

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOTP() string {

	rand.Seed(time.Now().UnixNano())

	otp := rand.Intn(900000) + 100000 //generate random 6 digit number 

	return fmt.Sprintf("%d", otp) //convert to string 
}
