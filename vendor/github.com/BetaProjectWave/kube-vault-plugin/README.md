# Kube Vault plugin
This plugin uses the vault kube auth jwt token to login to vault server and load all the secrets into memory or a given service

[![CircleCI](https://ci.shared.astoapp.co.uk/gh/BetaProjectWave/kube-vault-plugin.svg?style=svg&circle-token=5cff9cb8d9b06e3eafa1ff22739cee37a200de34)](https://ci.shared.astoapp.co.uk/gh/BetaProjectWave/kube-vault-plugin)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/30803a2f9a05478eac1f77079da176d3)](https://www.codacy.com?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=BetaProjectWave/kube-vault-plugin&amp;utm_campaign=Badge_Grade)
[![codecov](https://codecov.io/gh/BetaProjectWave/kube-vault-plugin/branch/master/graph/badge.svg?token=N6ZnqHV0Fk)](https://codecov.io/gh/BetaProjectWave/kube-vault-plugin)

Here is an example of how this client can be used

```go
package main

import (
	"github.com/BetaProjectWave/kube-vault-plugin"
	"fmt"
)

func main(){
	config := vault.Config{Address:"https://vault.do"}
	client := vault.NewClient(config)
	secret := client.ReadSecret("secret1")
	fmt.Println(secret)
}

```
