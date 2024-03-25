package structs

type SignUpStruct struct {
	Id       string `bson:"_id"`
	Name     string
	Surname  string
	Login    string
	Password string
}

type PostStruct struct {
	Id          string `bson:"_id"`
	Owner_id    string
	Title       string
	Image       string
	Description string
	Like        int
}

type LikeStruct struct {
	Id      string `bson:"_id"`
	Post_id string
	User_id string
}
