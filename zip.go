package podule

import "sync"

// Zip struct
type Zip struct{}

var instance *Zip
var once sync.Once

func GetInstance() *Zip {
	once.Do(func() {
		instance = &Zip{}
	})
	return instance
}
