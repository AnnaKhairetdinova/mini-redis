package mini_redis

import "github.com/AnnaKhairetdinova/mini-redis/store"

func main() {

	s := store.New()
	s.Get("hello")
}
