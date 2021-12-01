package typefile

type Command struct {
	Id       int
	Name     string
	Comment  string
	ParentId int
	IsDir    bool
}
