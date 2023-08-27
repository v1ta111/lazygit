package helpers

import (
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/status"
)

type AppStatusHelper struct {
	c *HelperCommon

	statusMgr func() *status.StatusManager
}

func NewAppStatusHelper(c *HelperCommon, statusMgr func() *status.StatusManager) *AppStatusHelper {
	return &AppStatusHelper{
		c:         c,
		statusMgr: statusMgr,
	}
}

func (self *AppStatusHelper) Toast(message string) {
	self.statusMgr().AddToastStatus(message)

	self.renderAppStatus()
}

// withWaitingStatus wraps a function and shows a waiting status while the function is still executing
func (self *AppStatusHelper) WithWaitingStatus(message string, f func(gocui.Task) error) {
	self.c.OnWorker(func(task gocui.Task) {
		self.statusMgr().WithWaitingStatus(message, func() {
			self.renderAppStatus()

			if err := f(task); err != nil {
				self.c.OnUIThread(func() error {
					return self.c.Error(err)
				})
			}
		})
	})
}

func (self *AppStatusHelper) WithWaitingStatusSync(message string, f func() error) {
	self.statusMgr().WithWaitingStatus(message, func() {
		stop := make(chan struct{})
		defer func() { stop <- struct{}{} }()
		self.renderAppStatusSync(stop)

		if err := f(); err != nil {
			_ = self.c.Error(err)
		}
	})
}

func (self *AppStatusHelper) HasStatus() bool {
	return self.statusMgr().HasStatus()
}

func (self *AppStatusHelper) GetStatusString() string {
	return self.statusMgr().GetStatusString()
}

func (self *AppStatusHelper) renderAppStatus() {
	self.c.OnWorker(func(_ gocui.Task) {
		ticker := time.NewTicker(time.Millisecond * 50)
		defer ticker.Stop()
		for range ticker.C {
			appStatus := self.statusMgr().GetStatusString()
			self.c.OnUIThread(func() error {
				self.c.SetViewContent(self.c.Views().AppStatus, appStatus)
				return nil
			})

			if appStatus == "" {
				return
			}
		}
	})
}

func (self *AppStatusHelper) renderAppStatusSync(stop chan struct{}) {
	ticker := time.NewTicker(time.Millisecond * 50)
	go func() {
	outer:
		for {
			select {
			case <-ticker.C:
				appStatus := self.statusMgr().GetStatusString()
				self.c.SetViewContent(self.c.Views().AppStatus, appStatus)
				_ = self.c.GocuiGui().ForceRedrawView(self.c.Views().AppStatus)
			case <-stop:
				break outer
			}
		}
		ticker.Stop()
	}()
}
