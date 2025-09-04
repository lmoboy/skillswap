package GoBackend

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type User struct {
	Host bool
	conn *websocket.Conn
	UID  string
}

type Room struct {
	mutext sync.Mutex
	Users  map[string][]User
	UID    string
}

func randUID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = letters[rand.Int()%len(letters)]
	}
	return string(b)
}

func (r *Room) init() {
	r.mutext.Lock()
	r.Users = make(map[string][]User)
	r.UID = randUID()
	r.mutext.Unlock()

	fmt.Println("Room created with UID:", r.UID)
}
func (r *Room) getUsers(roomUID string) []User {
	r.mutext.Lock()
	defer r.mutext.Unlock()
	return r.Users[roomUID]
}
func (r *Room) addUser(roomUID string, host bool, conn *websocket.Conn) {
	r.mutext.Lock()
	defer r.mutext.Unlock()
	uid := randUID()
	newUser := User{
		Host: host,
		conn: conn,
		UID:  uid,
	}
	r.Users[roomUID] = append(r.Users[roomUID], newUser)
	fmt.Println("User added with UID:", uid, "to room:", roomUID)

}

func (r *Room) removeRoom(roomUID string) {
	r.mutext.Lock()
	defer r.mutext.Unlock()
	delete(r.Users, roomUID)
	fmt.Println("Room removed with UID:", roomUID)
}

func setupVideo() {
	fmt.Println("Starting server on :8080")
	fmt.Println("Random string:", randUID())
	http.ListenAndServe(":8080", nil)

}
