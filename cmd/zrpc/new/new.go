package new

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"

	"github.com/spf13/cobra"
)

const (
	layoutRepoUrl = "https://github.com/mike955/zrpc-layout"
)

var (
	CmdNew = &cobra.Command{
		Use:   "new",
		Short: "Create a service template",
		Long:  longDescriber,
		Run:   run,
	}
	pName string
)

const (
	longDescriber = `
Create a http or grpc project using new command.
	
Default create a grpc and http project,
Example:
	- zrpc new demo		# create a new demo directory as project directory
	- zrpc new		# using current directory as project directory
`
)

func run(cmd *cobra.Command, args []string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if len(args) > 0 {
		pName = args[0]
	} else {
		dirs := strings.Split(dir, "/")
		if len(dirs) == 0 {
			panic("project name is nill")
		} else {
			pName = dirs[len(dirs)-1]
			dir = strings.Join(dirs[:len(dirs)-1], "/")
		}
	}
	if err := new(ctx, dir, pName); err != nil {
		panic(err)
	}
	fmt.Println("=========================")
	fmt.Println("===== init success ======")
	fmt.Println("=========================")
}

func new(ctx context.Context, dir, pName string) (err error) {
	to := path.Join(dir, pName)
	_, err = os.Stat(to)
	if !os.IsNotExist(err) {
		return fmt.Errorf("dir: %s is already exists", to)
	}
	dir = to

	if err = copyTemplate(ctx, dir, pName); err != nil {
		return fmt.Errorf("copy template error: %s", err.Error())
	}
	updateFile(dir, pName)
	// TODO(mike.cai) generate proto compile file
	return nil
}

func copyTemplate(ctx context.Context, dir, pName string) (err error) {
	_, err = git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
		URL:      layoutRepoUrl,
		Progress: os.Stdout,
		// Auth: &http.BasicAuth{
		// 	Username: "zrpc-layout",
		// 	Password: "enJwYy1sYXlvdXQ=",
		// },
	})
	return err
}

var removeFiles = []string{
	".idea",
	".git",
}

var replaceFiles = []string{
	"global.yml",
	"api/layout.proto",
	"go.mod",
	"cmd/layout/main.go",
	"internal/dao/init.go",
	"internal/dao/layout.go",
	"internal/data/layout.go",
	"internal/server/grpc.go",
	"internal/server/http.go",
	"internal/service/layout.go",
}

var renameFiles = []string{
	"api/layout.proto",
	"cmd/layout",
	"internal/service/layout.go",
	"internal/data/layout.go",
	"internal/dao/layout.go",
}

func updateFile(dir, pName string) {
	// remove file
	for _, f := range removeFiles {
		os.RemoveAll(path.Join(dir, f))
	}

	// replace ctx
	for _, f := range replaceFiles {
		fmt.Println(path.Join(dir, f))
		cmd := exec.Command("sed", "-i", "s/github.com\\/mike955\\/zrpc-layout/"+pName+"/g", path.Join(dir, f))
		if err := cmd.Run(); err != nil {
			os.RemoveAll(dir)
			panic(err)
		}

		cmd = exec.Command("sed", "-i", "s/layout/"+pName+"/g", path.Join(dir, f))
		if err := cmd.Run(); err != nil {
			os.RemoveAll(dir)
			panic(err)
		}

		upperName := strings.Replace(strings.Title(pName), "-", "", -1)
		cmd = exec.Command("sed", "-i", "s/Layout/"+upperName+"/g", path.Join(dir, f))
		if err := cmd.Run(); err != nil {
			os.RemoveAll(dir)
			panic(err)
		}
	}

	// rename files
	for _, f := range renameFiles {
		newName := strings.Replace(f, "layout", pName, 1)
		os.Rename(path.Join(dir, f), path.Join(dir, newName))
	}
}
