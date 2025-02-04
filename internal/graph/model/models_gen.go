// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AccessToken struct {
	ID             string  `json:"id"`
	UserID         string  `json:"userId"`
	RefreshTokenID string  `json:"refreshTokenId"`
	CreateTime     string  `json:"createTime"`
	UpdateTime     *string `json:"updateTime,omitempty"`
}

type CreateNoteInput struct {
	Title       *string   `json:"title,omitempty"`
	Content     *string   `json:"content,omitempty"`
	ObjectNames []*string `json:"objectNames,omitempty"`
}

type CreatePresignedUrlsResponse struct {
	Urls []*PresignedURL `json:"Urls,omitempty"`
}

type File struct {
	ID            string  `json:"id"`
	NoteID        string  `json:"noteId"`
	OriginalFile  string  `json:"originalFile"`
	ProcessedFile *string `json:"processedFile,omitempty"`
	URL           string  `json:"url"`
	CreateTime    string  `json:"createTime"`
	UpdateTime    *string `json:"updateTime,omitempty"`
}

type Mutation struct {
}

type Note struct {
	ID         string  `json:"id"`
	UserID     string  `json:"userId"`
	Title      *string `json:"title,omitempty"`
	Content    *string `json:"content,omitempty"`
	Files      []*File `json:"files,omitempty"`
	CreateTime string  `json:"createTime"`
	UpdateTime *string `json:"updateTime,omitempty"`
}

type NotesInput struct {
	Cursor *string `json:"cursor,omitempty"`
	Trash  *bool   `json:"trash,omitempty"`
}

type NotesResponse struct {
	Notes  []*Note `json:"notes,omitempty"`
	Cursor string  `json:"cursor"`
}

type PresignedURL struct {
	URL      string `json:"Url"`
	File     string `json:"File"`
	ObjectID string `json:"ObjectId"`
}

type Query struct {
}

type RefreshToken struct {
	ID         string  `json:"id"`
	UserID     string  `json:"userId"`
	CreateTime string  `json:"createTime"`
	UpdateTime *string `json:"updateTime,omitempty"`
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UpdateNoteInput struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
}

type User struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	CreateTime string  `json:"createTime"`
	UpdateTime *string `json:"updateTime,omitempty"`
}
