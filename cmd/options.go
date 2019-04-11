package cmd

type Options struct {
	AppName    string
	ModuleName string
	Aws        *Aws
	Terraform  *Terraform
}

type Aws struct {
	Region string
}

type Terraform struct {
	ModuleVersion string
}
