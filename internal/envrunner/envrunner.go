package envrunner

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretConfig struct {
	SecretName string `json:"secret_name"`
	Prefix     string
}

// Parse secret configuration from the SECRETS environment variable
func ParseSecretConfigs() []SecretConfig {
	var secretConfigs []SecretConfig
	secretConfigsJson := os.Getenv("SECRETS")
	json.Unmarshal([]byte(secretConfigsJson), &secretConfigs)
	return secretConfigs
}

// Get key=value environment variables from a secret
//
// This will fetch the JSON string from SecretsManager and then transform the map into a list of pairs.
func GetSecretEnv(client *secretsmanager.Client, sc SecretConfig) []string {
	output, err := client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: &sc.SecretName,
	})
	if err != nil {
		log.Fatal(err)
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(*output.SecretString), &result)
	var env []string
	for key, value := range result {
		env = append(env, key+"="+value.(string))
	}
	return env
}

// Get key=value environment variables from all secrets
func GetAllSecretEnv() []string {
	// Load AWS configuration using the credential chain
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := secretsmanager.NewFromConfig(cfg)
	var env []string
	var wg sync.WaitGroup
	for _, sc := range ParseSecretConfigs() {
		// Increment the WaitGroup counter.
		wg.Add(1)
		// Launch a goroutine to fetch the secret
		go func(sc SecretConfig) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			// Get a list of key=value pairs for the secret
			env = append(env, GetSecretEnv(client, sc)...)
		}(sc)
	}
	// Wait for all secrets to be fetched
	wg.Wait()
	return env
}

// Run the command with environment variables sourced from secrets
func Run(command []string) {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = append(os.Environ(), GetAllSecretEnv()...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
