package visitor

type urlJoint struct {
	contentURL string
	containURL string
}
type contentJoint struct {
	URL     string
	content string
}

func ConvertToUrlJoint(ele interface{}) (urlJoint, bool) {
	joint, ok := ele.(urlJoint)
	if ok {
		return joint, true
	}
	return urlJoint{contentURL: "gakki", containURL: "gakki"}, false
}

func ConvertToContentJoint(ele interface{}) (contentJoint, bool) {
	joint, ok := ele.(contentJoint)
	if ok {
		return joint, true
	}
	return contentJoint{URL: "gakki", content: "gakki"}, false
}
