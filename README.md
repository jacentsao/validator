# validator
go validator

#Installation
go get -u github.com/jacentsao/validator

#Usage

```golang
func main() {
	type User struct {
		Id    int    `json:"id" validate:"number,min=1,max=1000"`
		Name  string `validate:"string,min=2,max=10"`
		Bio   string `validate:"string"`
		Email string `validate:"email"`
	}

	user := User{
		Id:    111110,
		Name:  "superlongstring",
		Bio:   "",
		Email: "foobar",
	}

	fmt.Println("Errors:")
	for i, err := range validator.ValidateStruct(user) {
		fmt.Printf("\t%d. %s\n", i+1, err.Error())
	}
}
```