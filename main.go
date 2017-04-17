package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		// isolated process
		run()
	case "child":
		// don't fork
		child()
	default:
		panic("What?")
	}
}

func run() {
	// Creating the chroot
	// 1. sudo vi /etc/schroot/schroot.conf
	// [xenial]
	//description=Ubuntu Xenial
	//location=/home/chamilad/rootfs
	//priority=3
	//users=chamilad
	//groups=sbuild
	//root-groups=root
	// 2. sudo debootstrap --variant=buildd --arch amd64 xenial /home/chamilad/rootfs/ http://mirror.cc.columbia.edu/pub/linux/ubuntu/archive/

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		// new namespace and process id for my command, then create a new filesystem
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}

func child() {
	fmt.Printf("running [%v] as PID %v\n", os.Args[2], os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Chroot and chdir to /, then mount /proc
	must(syscall.Chroot("/home/chamilad/rootfs"))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
