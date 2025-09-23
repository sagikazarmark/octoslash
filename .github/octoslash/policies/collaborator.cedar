permit(
	principal in Role::"Collaborator",
	action == [Action::"close", Action::"label", Action::"remove-label"],
	resource
);
