package domain

type FileCloudRepository interface {
	Create() error
	Update() error
	Delete() error
	List() error
	Move() error
	Process() error
}
