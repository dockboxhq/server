package server

func Init() {
	r := NewRouter()
	r.Run(":8000")
}
