package confluence

type API struct{}

func New() *API {
	return &API{}
}

func (a *API) CreatePage() error {
	return nil
}

func (a *API) UpdatePage() error {
	return nil
}

func (a *API) FindPage() error {
	return nil
}
