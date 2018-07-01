package main

import (
    "os"
    "log"
	"fmt"
	"os/exec"
	"path/filepath"
    "net/http"

	"code.gitea.io/git"
)


func main() {
	os.MkdirAll("doc", 0755)

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		docdir := "doc/"+r.URL.Path[1:]
		repodir := "repo/"+r.URL.Path[1:]
		if _, err := os.Stat(docdir); !os.IsNotExist(err) {
			http.ServeFile(w, r, docdir)
		} else {
			if _, err := os.Stat(repodir); os.IsNotExist(err) {
				if err := os.MkdirAll(filepath.Dir(repodir), 0755); err != nil {
					fmt.Fprintf(w, "Error mkdirall: %v", err)
					return
				}
				if err := git.Clone("https://" + r.URL.Path[1:], repodir, git.CloneRepoOptions{Branch: "master"}); err != nil {
					fmt.Fprintf(w, "Error gitClone: %v", err)
					return
				}
			}
			if err := os.MkdirAll(filepath.Dir(docdir), 0755); err != nil {
				fmt.Fprintf(w, "Error mkdirall: %v", err)
				return
			}
			cmd := exec.Command("./phpDocumentor.phar", "-t", docdir, "-d", repodir)
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(w, "Error cmd: %v", err)
				return
			}
			http.ServeFile(w, r, docdir)
		}
	})

	log.Println("listening on", ":18080")
	log.Fatal(http.ListenAndServe(":18080", nil))
}
