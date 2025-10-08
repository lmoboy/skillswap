package swaps

import (
	"skillswap/backend/database"
	"skillswap/backend/structs"
)

func IncreaseSwaps(usr structs.UserInfo, amount int) {
	swapsAvailable := GetSwaps(usr)
	database.Execute("UPDATE users SET swaps = ? WHERE id = ?", swapsAvailable+amount, usr.ID)
}

func DecreaseSwaps(usr structs.UserInfo, amount int) {
	swapsAvailable := GetSwaps(usr)
	database.Execute("UPDATE users SET swaps = ? WHERE id = ?", swapsAvailable-amount, usr.ID)
}

func ExchangeSwaps(swapper structs.UserInfo, swappee structs.UserInfo) {
	DecreaseSwaps(swappee, 1)
	IncreaseSwaps(swapper, 2)
}

func GetSwaps(usr structs.UserInfo) int {
	swapsAvailable := 0
	swps := database.QueryRow("SELECT swaps FROM users WHERE id = ?", usr.ID)
	swps.Scan(&swapsAvailable)
	return swapsAvailable
}
