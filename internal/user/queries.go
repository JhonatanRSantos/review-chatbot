package user

var createUser = `
	INSERT INTO users (id, first_name, last_name, email)
	VALUES (:id, :first_name, :last_name, :email);
`

var findUserByEmail = `
	SELECT id, first_name, last_name, email FROM users WHERE email = :email;
`
