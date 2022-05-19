package main

import "gopkg.in/toast.v1"

func ShowNotificationWithTitle(title string, message string) {

	notification := toast.Notification{
		AppID:   "window-manager",
		Title:   title,
		Message: message,
	}
	Check(notification.Push())
}

func ShowNotification(message string) {
	notification := toast.Notification{
		AppID:   "window-manager",
		Message: message,
	}
	Check(notification.Push())
}
