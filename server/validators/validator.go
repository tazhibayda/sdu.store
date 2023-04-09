package validators

type Validator struct{
	errors []string
}

func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

func (v *Validator) Errors() []string {
	return v.errors
}
