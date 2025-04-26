package custom_errors

type InvalidLinkError struct {
	Link string
}

func (e InvalidLinkError) Error() string {
	return "Invalid link: " + e.Link
}

type InvalidRepoEventsError struct {
	Events []string
}

func (e InvalidRepoEventsError) Error() string {
	ans := ""
	for _, event := range e.Events {
		ans += event + " "
	}
	return "Invalid Events: " + ans
}
