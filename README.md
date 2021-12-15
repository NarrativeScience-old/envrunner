# envrunner

Utility that resolves a list of AWS SecretsManager secrets and sets environment variables before executing the passed command.

This is handy when running ECS or Lambda containers that need secrets from many separate SecretsManager secrets.

## Usage

```bash
export SECRETS='[{"secret_name": "my/secret", "prefix": ""}]'
envrunner launch-my-server --arg foo
```

## Configuration

This script assumes that the `SECRETS` environment variable is set. Its value should be a JSON list of `{"secret_name": "...", "prefix": "..."}` where the `secret_name` is the name of the secret in SecretsManager to load into the environment. Each secret is a JSON object whose fields are to be merged into the subprocess's environment. `prefix` is prepended onto each environment variable's key before injecting the value.

For example, given a secret `foo/rds` with keys `HOST`, `PORT`, `USERNAME`, and `PASSWORD` and `SECRETS='[{"secret_name": "foo/rds", "prefix": "DB_"}]`, the subprocess's environment will have variables `DB_HOST`, `DB_PORT`, `DB_USERNAME`, and `DB_PASSWORD`.
