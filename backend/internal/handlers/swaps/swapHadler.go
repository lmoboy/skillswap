package swaps

import (
	"skillswap/backend/internal/database"
	"skillswap/backend/internal/models"
)

func IncreaseSwaps(usr models.UserInfo, amount int) {
	swapsAvailable := GetSwaps(usr)
	database.Execute("UPDATE users SET swaps = ? WHERE id = ?", swapsAvailable+amount, usr.ID)
}

func DecreaseSwaps(usr models.UserInfo, amount int) {
	swapsAvailable := GetSwaps(usr)
	database.Execute("UPDATE users SET swaps = ? WHERE id = ?", swapsAvailable-amount, usr.ID)
}

func ExchangeSwaps(swapper models.UserInfo, swappee models.UserInfo) {
	DecreaseSwaps(swappee, 1)
	IncreaseSwaps(swapper, 2)
}

func GetSwaps(usr models.UserInfo) int {
	swapsAvailable := 0
	swps := database.QueryRow("SELECT swaps FROM users WHERE id = ?", usr.ID)
	swps.Scan(&swapsAvailable)
	return swapsAvailable
}
