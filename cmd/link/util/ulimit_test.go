// +build !windows

package util

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
)

func TestManageFdLimit(t *testing.T) {
	t.Log("Testing file descriptor count")
	if _, _, err := ManageFdLimit(); err != nil {
		t.Errorf("Cannot manage file descriptors")
	}

	if maxFds != uint64(8192) {
		t.Errorf("Maximum file descriptors default value changed")
	}
}

func TestManageInvalidNFds(t *testing.T) {
	t.Logf("Testing file descriptor invalidity")
	var err error
	if err = os.Unsetenv("LINK_FD_MAX"); err != nil {
		t.Fatal("Cannot unset the LINK_FD_MAX env variable")
	}

	rlimit := syscall.Rlimit{}
	if err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		t.Fatal("Cannot get the file descriptor count")
	}

	value := rlimit.Max + rlimit.Cur
	if err = os.Setenv("LINK_FD_MAX", fmt.Sprintf("%d", value)); err != nil {
		t.Fatal("Cannot set the LINK_FD_MAX env variable")
	}

	t.Logf("setting ulimit to %d, max %d, cur %d", value, rlimit.Max, rlimit.Cur)

	if changed, new, err := ManageFdLimit(); err == nil {
		t.Errorf("ManageFdLimit should return an error: changed %t, new: %d", changed, new)
	} else {
		flag := strings.Contains(err.Error(),
			"failed to raise ulimit to LINK_FD_MAX")
		if !flag {
			t.Error("ManageFdLimit returned unexpected error", err)
		}
	}

	// unset all previous operations
	if err = os.Unsetenv("LINK_FD_MAX"); err != nil {
		t.Fatal("Cannot unset the LINK_FD_MAX env variable")
	}
}

func TestManageFdLimitWithEnvSet(t *testing.T) {
	t.Logf("Testing file descriptor manager with LINK_FD_MAX set")
	var err error
	if err = os.Unsetenv("LINK_FD_MAX"); err != nil {
		t.Fatal("Cannot unset the LINK_FD_MAX env variable")
	}

	rlimit := syscall.Rlimit{}
	if err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		t.Fatal("Cannot get the file descriptor count")
	}

	value := rlimit.Max - rlimit.Cur + 1
	if err = os.Setenv("LINK_FD_MAX", fmt.Sprintf("%d", value)); err != nil {
		t.Fatal("Cannot set the LINK_FD_MAX env variable")
	}

	if _, _, err = ManageFdLimit(); err != nil {
		t.Errorf("Cannot manage file descriptor count")
	}

	// unset all previous operations
	if err = os.Unsetenv("LINK_FD_MAX"); err != nil {
		t.Fatal("Cannot unset the LINK_FD_MAX env variable")
	}
}
