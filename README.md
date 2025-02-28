# List Agent CLI

## Detailed Command Documentation
Detailed command documentation can be found here is auto generated and published to GitHub Pages:
https://rad12000.github.io/list-agent.

## Installing
**NOTE**: This CLI tool hasn't been tested in Powershell or Fish.
1) Install the [latest release](https://github.com/rad12000/list-agent/releases/latest) from GitHub.
2) Unzip the release artifact and move the `listagent` executable to a location within your `$PATH`.
   (Note: you may need to grant executable permissions like so: `chmod +x listagent`)
3) Recommended: [Setup shell completions](#setting-up-completions)

## Usage
The listagent CLI is made up of subcommands, each with their own unique arguments.
In general, usage will look like the following:
```bash
listagent [subcommand] <subcommand args>
```

## Setting up completions
What is a CLI tool without autocompletion and intelligent suggestions?!

### Bash
```bash
# Make sure you have the 'bash-completion' installed. Odds are you already do.
source <(listagent completion bash) # set up autocomplete in bash into the current shell, bash-completion package should be installed first.
echo "source <(listagent completion bash)" >> ~/.bashrc # add autocomplete permanently to your bash shell.
```

### Zsh
```zsh
# set up autocomplete in zsh into the current shell
source <(listagent completion zsh)
# add autocomplete permanently to your zsh shell
echo 'autoload -U compinit; compinit' >> ~/.zshrc
echo '[[ $commands[listagent] ]] && source <(listagent completion zsh)' >> ~/.zshrc
```

### Powershell
1) Add the following line to your powershell profile: `listagent completion powershell | Out-String | Invoke-Expression`
2) Restart shell session.

### Fish and more
For more information on setting up Fish, or for more granular control over your completions, run `listagent completion [zsh|bash|fish|powershell] --help`