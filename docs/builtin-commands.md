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

## `/add-label <label>`

**Aliases:** `label`

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

## `/assign <username>`

Assign an issue or pull request to a specific user.

```
/assign johndoe
/assign "jane.smith"
```

**Required Permission**: `Assign` action on the resource

## `/self-assign`

Assign an issue or pull request to yourself (the comment author).

```
/self-assign
```

**Required Permission**: `Assign` action on the resource

## `/unassign <username>`

Unassign a specific user from an issue or pull request.

```
/unassign johndoe
/unassign "jane.smith"
```

**Required Permission**: `Unassign` action on the resource

## `/self-unassign`

Unassign yourself (the comment author) from an issue or pull request.

```
/self-unassign
```

**Required Permission**: `Unassign` action on the resource
