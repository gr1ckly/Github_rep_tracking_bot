package LinkTracker

import "fmt"

type NotificationService struct {
	client struct{}
}

func NewNotificationService() (*NotificationService, error) {
	return &NotificationService{struct{}{}}, nil
}

func (ns *NotificationService) Notify(msg any) {
	fmt.Println("Send Object")
}
