package testhelper

// structs

type Person struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Area string `json:"area"`
	Age  int    `json:"age"`
}

type Areas struct {
	Name string `json:"name"`
}
