package main

import (
	"fmt"
	"sss/IhomeWeb/utils"
)

func main(){
	o := []byte("123")
	eo := utils.EncryptAES(o)
	fmt.Println(string(eo))

	do := utils.DecryptAES(eo)
	fmt.Println(string(do))
}
