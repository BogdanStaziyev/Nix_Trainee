package response

type CommentResponse struct {
	ID     int64  `json:"id"`
	PostID int64  `json:"post_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}
