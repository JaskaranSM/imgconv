package errors

import "errors"

var ErrorImageNotFound error = errors.New("image not found")
var IdParamNotFound error = errors.New("id parameter not found")
var ConversionTaskNotFound error = errors.New("conversion task not found")
