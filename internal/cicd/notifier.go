package cicd

import "log"

// Event types for CI/CD notifications.
const (
	EventBuildFailed  = "build_failed"
	EventBuildSuccess = "build_success"
	EventDeployed     = "deployed"
)

// Notification represents a CI/CD event to be stored or forwarded.
type Notification struct {
	Event  string
	Repo   string
	Branch string
	Message string
}

// Notifier handles CI/CD event dispatch.
type Notifier struct{}

// NewNotifier creates a new Notifier.
func NewNotifier() *Notifier {
	return &Notifier{}
}

// Notify logs and dispatches a CI/CD notification.
func (n *Notifier) Notify(notif Notification) {
	// TODO: persist to Firestore, forward to webhook or Pub/Sub
	log.Printf("[cicd] event=%s repo=%s branch=%s msg=%s",
		notif.Event, notif.Repo, notif.Branch, notif.Message)
}
