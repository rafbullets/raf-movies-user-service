package user

// User - model for user table
type User struct {
	ID                   int64  `json:"id" db:"id"`
	Username             string `json:"username" db:"username"`
	Password             string `json:"password" db:"password"`
	Role                 int64  `json:"role" db:"role"`
	FirstName            string `json:"firstName" db:"first_name"`
	LastName             string `json:"lastName" db:"last_name"`
	Email                string `json:"email" db:"email"`
	NumberOfReservations int64  `json:"numberOfReservations" db:"number_of_reservations"`
}
