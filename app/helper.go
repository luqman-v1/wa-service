package app

import (
	"os"
	"strconv"
	"strings"
)

const FILENAME_SESSION = "whatsappSession.gob"
const PATH_SESSION = "./storage/session/" + FILENAME_SESSION
const PATH_FILE = "./storage/file/"

func StoI(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func WidToJid(w string) string {
	return strings.Split(w, "@")[0] + "@s.whatsapp.net"
}

func FilePath(filename string) string {
	return "https://" + os.Getenv("AWS_BUCKET") + "." + "s3-" + os.Getenv("AWS_REGION") + ".amazonaws.com/" + filename
}
