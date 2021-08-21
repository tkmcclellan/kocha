package providers

type Provider struct {
	name string
}

func (provider Provider) ProviderName() string {
	return provider.name
}
