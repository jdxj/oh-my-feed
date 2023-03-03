package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	t := time.Time{}
	fmt.Println(t.Format(time.DateTime))
	h := md5.New()
	h.Write([]byte("hello"))
	res := h.Sum(nil)
	ss := hex.EncodeToString(res)
	fmt.Printf("%s\n", ss)
	ss = base64.StdEncoding.EncodeToString(res)
	fmt.Printf("%s\n", ss)
	res, err := bcrypt.GenerateFromPassword([]byte("abc"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s, %d\n", res, len(res))

	res, err = bcrypt.GenerateFromPassword([]byte("88**$/abc"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s, %d\n", res, len(res))
}
