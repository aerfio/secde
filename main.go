package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

func main() {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		err := decodeSecret(os.Stdin)
		if err != nil {
			panic(err)
		}

	} else {
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err)
		}
		coreClient, err := typedcorev1.NewForConfig(config)
		if err != nil {
			panic(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		secret, err := coreClient.Secrets(os.Args[1]).Get(ctx, os.Args[2], metav1.GetOptions{})
		if err != nil {
			panic(err)
		}
		buf := bytes.NewBuffer(nil)
		if err := json.NewEncoder(buf).Encode(secret); err != nil {
			panic(err)
		}
		if err := decodeSecret(buf); err != nil {
			panic(err)
		}
	}
}

func decodeSecret(r io.Reader) error {
	readerBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	sec := &corev1.Secret{}
	err = yaml.Unmarshal(readerBytes, sec)
	if err != nil {
		return err
	}

	dckrCfg, ok := sec.Data[corev1.DockerConfigJsonKey]
	if !ok {
		return fmt.Errorf("no %q key found", corev1.DockerConfigJsonKey)
	}
	registries := Registries{}
	if err := json.Unmarshal(dckrCfg, &registries); err != nil {
		return err
	}

	for key, val := range registries.Auths {
		fmt.Printf("registry: %s, username: %s, pass: %s\n", key, val.Username, val.Password)
	}

	return nil
}

type Registries struct {
	Auths map[string]Auth `json:"auths"`
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
