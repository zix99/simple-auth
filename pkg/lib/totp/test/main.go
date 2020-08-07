package main

import (
	"fmt"
	"os"
	"simple-auth/pkg/lib/totp"
)

const testSecret = "WHEWDNY="

func main() {
	var otp *totp.Totp
	args := os.Args[1:]
	if len(args) > 0 {
		fmt.Println("Parsing OTP...")
		otp, _ = totp.ParseTOTP(args[0])
	} else {
		fmt.Println("New OTP")
		//otp, _ = totp.NewTOTP(4, "coolco", "chris")
		otp, _ = totp.FromSecret(testSecret, "coolco", "chris")
	}
	fmt.Println(otp.String())
	fmt.Println(otp.GetTOTP())

	//f, _ := os.Create("test.png")
	//defer f.Close()
	//png, _ := otpimagery.GenerateQRCode(otp, 256)
	//f.Write(png)
}
