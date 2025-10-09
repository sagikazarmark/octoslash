permit(
	principal in Role::"Triager",
	action in [Action::"close", Action::"add-label", Action::"remove-label"],
	resource is Issue
);
