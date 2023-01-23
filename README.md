# secde

Image pull secret decoder

## Usage

To decode image pull secret in cluster:

```bash
secde ${SECRET_NAMESPACE} ${SECRET_NAME}
```

To decode secret in local yaml file:

```bash
cat regcred.yaml | secde
```

## Alternatives

Good old bash scripts using jq/yq:
- To decode any k8s secret in YAML format use
    ```bash
    yqd () {
        yq ".data | map_values(@base64d)"
    }
    ```
  Example:
  ```
    kubectl get secret -n ${SECRET_NAMESPACE} ${SECRET_NAME} -oyaml | yqd
  ```
- To decode k8s secret in JSON format use
    ```bash
    jqd () {
        jq '.data | map_values(@base64d)'
    }
    ```
    or use `yqd` script, because JSON is also YAML.
    Example:
    ```bash
    kubectl get secret -n ${SECRET_NAMESPACE} ${SECRET_NAME} -ojson | jqd
    ```
- To decode image pull secret use:
    ```bash
    jqdockersec () {
        jq '.data' | jq '.".dockerconfigjson"' | xargs | base64 -D | jq
    }
    ```
    Usage is the same as with previous scripts, pipe the whole secret to this script
