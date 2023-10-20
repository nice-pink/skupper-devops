package deploy

import (
	"fmt"
	"time"

	filesystem "github.com/nice-pink/gm_filesystem"
	git "github.com/nice-pink/gm_git"
	"github.com/nice-pink/skupper-devops/pkg/kynetes"
)

const (
	skupperSecret string = "skupper-secret-default"
)

// prepare

func Prepare(src string, dest string, gitPushDest bool) {
	// copy src to dest if dest is defined
	if dest != "" {
		err := filesystem.CopyDir(src, dest, false, true)
		if err != nil {
			panic(err)
		}
	} else {
		// src was manually copied to dest. Use src as dest.
		fmt.Println("src is already dest!")
		dest = src
	}

	// push dest via git
	if gitPushDest {
		gitMessage := "Add skupper to " + dest
		git.CommitPushLocalRepo(".", gitMessage, true)
	}
}

// token

// func CreateTokenRequest(path string, namespace string) string {
// 	// token request
// 	tokenRequest := `apiVersion: v1
// kind: Secret
// metadata:
//   labels:
//     skupper.io/type: connection-token-request
//   name: skupper-secret-default
//   namespace: ` + namespace

// 	return tokenRequest
// }

func CreateTokenRequest(namespace string) string {
	// token request
	tokenRequest := `apiVersion: v1
kind: Secret
metadata:
  labels:
    nice.io/type: blalbla
  name: test-secret
  namespace: ` + namespace

	return tokenRequest
}

func SendTokenRequest(namespace string) error {
	token := CreateTokenRequest(namespace)

	err := kynetes.CreateSecret(namespace, []byte(token))
	if err != nil {
		return err
	}

	return nil
}

func TokenHasData(namespace string, retries int) bool {
	for index := 1; index <= retries; index++ {
		if kynetes.SecretHasData(skupperSecret, namespace) {
			return true
		} else {
			fmt.Print("Is empty...")
		}

		if index < retries-1 {
			time.Sleep(1 * time.Second)
			fmt.Println("Retry!")
		} else {
			fmt.Println("Failed!")
		}
	}
	return false
}
