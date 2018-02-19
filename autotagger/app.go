package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/log"
)

func main() {
	router := httprouter.New()
	router.GET("/hook/*project", handler)
	http.ListenAndServe(":8080", router)
}

func generateTagName() string {
	year, month, day := time.Now().UTC().Date()
	return fmt.Sprintf("%v%02v%02v", year, int(month), int(day))
}

func tagRepo(ctx context.Context, user string, repo string) error {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_AUTH_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	commits, _, err := client.Repositories.ListCommits(ctx, "rphillips", "kubernetes-base-image-builder", nil)
	if err != nil {
		return err
	}

	tagName := generateTagName()
	sha := commits[0].GetSHA()
	want := fmt.Sprintf("refs/tags/%v", tagName)
	t := "commit"
	ref := &github.Reference{
		Ref: &want,
		Object: &github.GitObject{
			Type: &t,
			SHA:  &sha,
		},
	}
	_, _, err = client.Git.CreateRef(ctx, "rphillips", "kubernetes-base-image-builder", ref)
	return err
}

func handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	project := ps.ByName("project")
	if project == "" {
		http.NotFound(w, r)
		return
	}

	// setup context
	ctx := context.Background()

	// get user/project name
	userRepo := strings.Split(project, "/")
	if len(userRepo) != 2 {
		err := fmt.Errorf("Project invalid format: %v", project)
		log.Errorf(ctx, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	err := tagRepo(ctx, userRepo[0], userRepo[1])
	if err != nil {
		log.Errorf(ctx, "Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "OK")
}
