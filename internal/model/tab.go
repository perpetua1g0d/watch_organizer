package model


type Tab struct {
	Id   int
	Name string
}

type TabChildren struct {
	Id1 int
	Id2 int
}

type TabQueue struct {
	TabId    int
	PosterId int
	Position int
}
