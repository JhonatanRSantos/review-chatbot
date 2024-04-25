package chat

var createChat = `
	INSERT INTO chats (id, user_id) VALUES (:id, :user_id);
`

var createMessage = `
	INSERT INTO messages (id, chat_id, author, message)
	VALUES (:id, :chat_id, :author, :message);
`
