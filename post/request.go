package post

import "strings"

type ValidationErrors map[string]string

func (e ValidationErrors) HasErrors() bool { return len(e) > 0 }

type postRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (req postRequest) validate() ValidationErrors {
	errs := make(ValidationErrors)

	title := strings.TrimSpace(req.Title)
	switch {
	case title == "":
		errs["title"] = "required"
	case len(title) > 255:
		errs["title"] = "must not exceed 255 characters"
	}

	if body := strings.TrimSpace(req.Body); len(body) > 1000 {
		errs["body"] = "must not exceed 1000 characters"
	}

	return errs
}
