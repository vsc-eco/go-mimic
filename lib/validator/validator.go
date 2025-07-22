package validator

import v "github.com/go-playground/validator/v10"

// NOTE: a single instance for the library to cache struct parsing results
var validatorEngine = v.New()

func New() *v.Validate { return validatorEngine }
