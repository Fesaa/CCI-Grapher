package utils

import "regexp"

func MaybeDiscordID(maybe string) (bool, error) {
    return regexp.Match(`^\d{18,19}$`, []byte(maybe))
}

func Max(a int, b int) int {
    if a > b {
        return a
    }
    return b
}
