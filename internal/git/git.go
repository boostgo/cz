package git

import (
	"github.com/boostgo/lite/errs"
	"os"
	"os/exec"
)

const errType = "GIT"

func Commit(message string) (err error) {
	defer errs.Wrap(errType, &err, "Commit")
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
