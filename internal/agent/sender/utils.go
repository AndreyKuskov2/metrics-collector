package sender

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net"
	"net/http"
)

const IPAPI = "https://api.ipify.org"

// Функция вычисления хэша
func calculateHash(data, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func getPublicIP() (string, error) {
	resp, err := http.Get(IPAPI)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP.String(), nil
}
