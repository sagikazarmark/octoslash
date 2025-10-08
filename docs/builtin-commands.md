# Built-in Commands

Octoslash comes with several built-in commands for common GitHub operations:

## `/close [reason]`

Close an issue or pull request with an optional reason.

```
/close
/close completed
/close not_planned
```

**Required Permission**: `Close` action on the resource

## `/label <label>`

Add a label to an issue or pull request.

```
/label bug
/label "needs review"
```

**Required Permission**: `Label` action on the resource

## `/remove-label <label>`

Remove a label from an issue or pull request.

```
/remove-label bug
/remove-label "needs review"
```

**Required Permission**: `RemoveLabel` action on the resource
