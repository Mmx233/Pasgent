# Pasgent

[![Lisense](https://img.shields.io/github/license/Mmx233/Pasgent)](https://github.com/Mmx233/Pasgent/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/Mmx233/Pasgent?color=blueviolet&include_prereleases)](https://github.com/Mmx233/Pasgent/releases)
[![GoReport](https://goreportcard.com/badge/github.com/Mmx233/Pasgent)](https://goreportcard.com/report/github.com/Mmx233/Pasgent)

Pasgent is a tool used to simulate Pageant on Windows, utilizing 1Password SSH Agent for SSH authentication.

Pasgent registers the Pageant Window Message Class `Pageant` and forwards ssh-agent requests to the 1Password SSH Agent's Named Pipe `\\.\pipe\openssh-ssh-agent` to implement the authentication simulation process.

**Please note that this tool cannot run simultaneously with Pageant.**

## :gear: How to use

Download the release and then run it.

![Command Line](screenshots/cmd.png)