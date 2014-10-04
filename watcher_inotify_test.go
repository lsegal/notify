// +build linux,!fsnotify

package notify_test

import (
	"io"
	"os"
	"testing"

	"github.com/rjeczalik/notify"
	"github.com/rjeczalik/notify/test"
)

// Access describes a set of access events.
var access = []notify.Event{
	notify.IN_OPEN,
	notify.IN_MODIFY,
	notify.IN_CLOSE_WRITE,
	notify.IN_OPEN,
	notify.IN_ACCESS,
	notify.IN_CLOSE_NOWRITE,
}

// InotifyActions extends the default feature with an action executor for
// an IN_ACCESS event.
var inotifyActions = test.Actions{
	notify.IN_ACCESS: func(p string) (err error) {
		f, err := os.OpenFile(p, os.O_RDWR, 0755)
		if err != nil {
			return
		}
		if _, err = f.WriteString(p); err != nil {
			f.Close()
			return
		}
		f.Close()
		if f, err = os.Open(p); err != nil {
			return
		}
		if _, err = f.Read(make([]byte, 1)); err != nil && err != io.EOF {
			f.Close()
			return
		}
		return f.Close()
	},
}

// TODO(rjeczalik): The following method should be autogenerated by reflection-based
// generator function from test (see TODO in test package for details).
func ExpectInotifyEvents(t *testing.T, wr notify.Watcher, e notify.Event,
	ei map[notify.EventInfo][]notify.Event) {
	w := test.W(t, inotifyActions)
	defer w.Close()
	if err := w.WatchAll(wr, e); err != nil {
		t.Fatal(err)
	}
	defer w.UnwatchAll(wr)
	w.ExpectEvents(wr, ei)
}

func TestInotify(t *testing.T) {
	ei := map[notify.EventInfo][]notify.Event{
		test.EI("github.com/rjeczalik/fs/fs.go", notify.IN_ACCESS): access,
		// test.EI("github.com/rjeczalik/fs/binfs/", notify.IN_MODIFY),
		// test.EI("github.com/rjeczalik/fs/binfs.go", notify.IN_ATTRIB),
		// test.EI("github.com/rjeczalik/fs/binfs_test.go", notify.IN_CLOSE_WRITE),
		// test.EI("github.com/rjeczalik/fs/binfs/", notify.IN_CLOSE_NOWRITE),
		// test.EI("github.com/rjeczalik/fs/binfs/", notify.IN_OPEN),
		// test.EI("github.com/rjeczalik/fs/fs_test.go", notify.IN_MOVED_FROM),
		// test.EI("github.com/rjeczalik/fs/binfs/", notify.IN_MOVED_TO),
		// test.EI("github.com/rjeczalik/fs/binfs.go", notify.IN_CREATE),
		// test.EI("github.com/rjeczalik/fs/binfs_test.go", notify.IN_DELETE),
		// test.EI("github.com/rjeczalik/fs/binfs/", notify.IN_DELETE_SELF),
		// test.EI("github.com/rjeczalik/fs/binfs/", notify.IN_MOVE_SELF),
	}
	ExpectInotifyEvents(t, notify.NewWatcher(), notify.IN_ALL_EVENTS, ei)
}
