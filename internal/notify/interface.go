package notify

// Beeper interface wraps the beeep library functions for testing
type Beeper interface {
	Notify(title, message string, icon interface{}) error
	Alert(title, message string, icon interface{}) error
	Beep(freq float64, duration int) error
	SetAppName(name string)
}

// DefaultBeeper implements Beeper using the actual beeep library
type DefaultBeeper struct{}

func NewDefaultBeeper() *DefaultBeeper {
	return &DefaultBeeper{}
}

// These functions will be implemented to wrap the actual beeep calls
func (b *DefaultBeeper) Notify(title, message string, icon interface{}) error {
	return beepNotify(title, message, icon)
}

func (b *DefaultBeeper) Alert(title, message string, icon interface{}) error {
	return beepAlert(title, message, icon)
}

func (b *DefaultBeeper) Beep(freq float64, duration int) error {
	return beepBeep(freq, duration)
}

func (b *DefaultBeeper) SetAppName(name string) {
	beepSetAppName(name)
}