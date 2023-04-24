package validators

import (
	"sdu.store/server/model"
)

type CommentValidator struct {
	*model.Comment
	Validator
}

func NewCommentValidator(comment *model.Comment) CommentValidator {
	return CommentValidator{comment, Validator{errors: []string{}}}
}

func (v *CommentValidator) Check() {
	if v.Text == "" {
		v.errors = append(v.errors, "Comment shouldn't be empty")
	}
	if _, err := model.GetUserByID(int64(v.UserID)); err != nil {
		v.errors = append(v.errors, "User doesn't exist")
	}
	if _, err := model.GetProductByID(v.ProductID); err != nil {
		v.errors = append(v.errors, "Product doesn't exist")
	}
}
