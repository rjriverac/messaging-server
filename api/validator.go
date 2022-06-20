package api

import "github.com/go-playground/validator/v10"

var validRequest validator.StructLevelFunc = func(sl validator.StructLevel) {
	info := sl.Current().Interface().(UpdateUserRequest)

	if len(info.Name) == 0 && len(info.Email) == 0 && len(info.Password) == 0 && len(info.Image) == 0 && len(info.Status) == 0 {
		sl.ReportError(info.Name, "name", "name", "empty request", "")
		sl.ReportError(info.Password, "Password", "empty request", "name", "")
		sl.ReportError(info.Image, "Image", "Image", "empty request", "")
		sl.ReportError(info.Email, "Email", "Email", "empty request", "")
		sl.ReportError(info.Status, "Status", "Status", "empty request", "")
	}

}
