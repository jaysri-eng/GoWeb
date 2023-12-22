A repository for a simple Golang project for storing username and passwords, hashing the passwords after passing and fetch all the details.

Please ensure to create a local SQL database and necessary tables to implement all functionalities, the table structure is given below:
type User struct {
	Id       int64
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}
