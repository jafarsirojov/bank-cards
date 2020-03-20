package auth

type Auth struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Expired int64    `json:"exp"`
}

var AnonymousAuth = Auth{
	Id: -1,
	Name: "Anonymous",
	Expired: 0,
}
