package resolvers

type key int

const (
	RequestKey key = iota
)

type RootResolver struct {
	MeResolver
	CreateUserResolver
}

func NewRootResolver(meResolver *MeResolver, createUserResolver *CreateUserResolver) *RootResolver {
	return &RootResolver{
		MeResolver:         *meResolver,
		CreateUserResolver: *createUserResolver,
	}
}
