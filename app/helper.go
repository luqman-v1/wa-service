package app

import (
    "strconv"
    "strings"
)

const PATH_SESSION = "./storage/session/whatsappSession.gob"
const PATH_FILE = "./storage/file/"

func StoI(s string) int{
    i ,_ := strconv.Atoi(s)
    return i
}

func WidToJid(w string) string {
    return strings.Split(w, "@")[0] + "@s.whatsapp.net"
}