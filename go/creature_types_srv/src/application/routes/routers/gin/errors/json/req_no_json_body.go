package gin_errors_json

type NoJsonBody struct{}

func (err *NoJsonBody) Error() string {
	return "Request miss json body"
}
