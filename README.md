# Octoslash

![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/sagikazarmark/octoslash/ci.yaml?style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sagikazarmark/octoslash?style=flat-square&color=61CFDD)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/sagikazarmark/octoslash/badge?style=flat-square)](https://deps.dev/go/github.com%252Fsagikazarmark%252Foctoslash)

**A Go library and CLI for building custom GitHub slash commands.**

> [!WARNING]
> Make sure to read the [Security](#security) section to learn about potential attack vectors and how to use octoslash securely.

## Features

- **Library**: Build custom slash commands using the octoslash Go library
- **Built-in Commands**: Ready-to-use commands for common GitHub operations
- **CLI tool**: Standalone binary with built-in commands for immediate use
- **GitHub Action**: Easy integration via [sagikazarmark/octoslash-action](https://github.com/sagikazarmark/octoslash-action)
- **Authorization**: Fine-grained access control using [Cedar](https://www.cedarpolicy.com) policies

## Quickstart

Check out [this](https://github.com/sagikazarmark/octoslash-demo) repository for a quickstart guide.

## Installation

### GitHub Action

See [sagikazarmark/octoslash-action](https://github.com/sagikazarmark/octoslash-action).

### Using the CLI

Download the latest release from the [releases page](https://github.com/sagikazarmark/octoslash/releases) or install using Go:

```bash
go install github.com/sagikazarmark/octoslash/cmd/octoslash@latest
```

### Using the library

Add octoslash to your Go project:

```bash
go get github.com/sagikazarmark/octoslash
```

## Built-in Commands

Check out [this](docs/builtin-commands.md) page for a list of built-in commands.

## Authorization

The octoslash _binary_ uses the [Cedar](https://www.cedarpolicy.com/) policy language for fine-grained authorization.
By default, **all commands are denied** unless explicitly allowed by a policy.

The _library_ allows alternative authorization mechanisms by implementing the appropriate interface.

### Configuration

Create authorization configuration in `.github/octoslash/`:

```
.github/octoslash/
├── principals.json         # User and role mappings
└── policies/
    ├── collaborator.cedar  # Policies for collaborators
    └── triager.cedar       # Policies for triagers
```

### Principals

Map GitHub users to roles in `principals.json`:

```json
[
    {
        "uid": { "type": "User", "id": "1226384" },
        "attrs": { "login": "sagikazarmark" },
        "parents": [{ "type": "Role", "id": "Collaborator" }]
    },
    {
        "uid": { "type": "User", "id": "987654321" },
        "attrs": { "login": "triager" },
        "parents": [{ "type": "Role", "id": "Triager" }]
    }
]
```

> [!NOTE]
> For the moment, repository members also have to be added as principals to assign roles to them.
>
> See [#9](https://github.com/octoslash/octoslash/issues/9) for details.

### Policy Examples

**Collaborator Policy** (`policies/collaborator.cedar`):
```cedar
// Collaborators can perform all actions
permit(
    principal in Role::"Collaborator",
    action,
    resource
);
```

**Triager Policy** (`policies/triager.cedar`):
```cedar
// Triagers can only close, label, and remove labels on issues (not PRs)
permit(
    principal in Role::"Triager",
    action in [Action::"Close", Action::"Label", Action::"RemoveLabel"],
    resource is Issue
);
```

## Building Custom Commands

TODO

## Security

### Attack Vectors and Mitigations

Octoslash implements several security measures to protect against common attack vectors:

#### 1. **Unauthorized Command Execution**
- **Risk**: Malicious users executing privileged commands
- **Mitigation**: Authorization with deny-by-default policy
- **Protection**: All commands require explicit policy permissions

#### 2. **Policy Bypass via PR Modifications**
- **Risk**: Malicious PRs modifying GitHub Action workflows or policies
- **Mitigation**: GitHub workflows for `issue_comment` events only run on the default branch
- **Protection**: Malicious actors cannot gain privileges by submitting malicious PRs

#### 3. **Command Injection**
- **Risk**: Malicious command arguments causing unintended behavior
- **Mitigation**: Structured command parsing using [mvdan/sh](https://github.com/mvdan/sh)
- **Protection**: Argument validation and type safety

### Security Best Practices

1. **Use minimal GitHub token permissions**:
   ```yaml
   permissions:
     issues: write
     # Disable (or rather don't enable) unnecessary permissions
     # pull-requests: write
   ```

1. **Implement granular Cedar policies**:
   ```cedar
   // Prefer specific permissions over broad access
   permit(
     principal in Role::"Triager",
     action == Action::"Label",
     resource is Issue
   );
   ```

1. **Regular policy audits**: Review and update authorization policies regularly

## Development

TODO

## License

The project is licensed under the [MIT License](LICENSE).
