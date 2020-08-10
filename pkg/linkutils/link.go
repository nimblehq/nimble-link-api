package linkutils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/nimble-link/backend/models"
)

func IsDuplicateAlias(alias string) bool {
	link := models.FindLinkByAlias(alias)
	return link.ID != 0
}

func GenerateHashFrom(originalURL string, userID uint) string {
	currentTime := time.Now().String()

	md5Hashed := getMD5Hash(strconv.FormatUint(uint64(userID), 10) + currentTime + originalURL)
	base64Hashed := getBase64Hash(md5Hashed)
	return base64Hashed[0:6]
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func getBase64Hash(text string) string {
	return base64.URLEncoding.EncodeToString([]byte(text))
}
