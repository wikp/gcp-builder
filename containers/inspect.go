package containers

type inspect []inspectedImage

type inspectedImage struct {
	Id string `json:"Id"`
}
