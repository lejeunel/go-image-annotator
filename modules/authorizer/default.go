package authorizer

var DefaultRules = []AuthRule{
	{Method: "Annotate"},
	{Method: "CreateCollection"},
	{Method: "CreateLabel"},
	{Method: "DeleteCollection"},
	{Method: "DeleteImage"},
	{Method: "DeleteLabel"},
	{Method: "ImportImage"},
	{Method: "IngestImage"},
	{Method: "RenewToken"},
	{Method: "UpdateCollection"},
	{Method: "UpdateLabel"},

	{Method: "SetAdminRights", AdminOnly: true},
	{Method: "DeleteUser", AdminOnly: true},
	{Method: "AssignRoleToUser", AdminOnly: true},
	{Method: "AssignUserToGroup", AdminOnly: true},
	{Method: "ListUsers", AdminOnly: true},
	{Method: "FindUser", AdminOnly: true},
	{Method: "UnAssignRoleFromUser", AdminOnly: true},
	{Method: "UnAssignUserFromGroup", AdminOnly: true},
	{Method: "CreateUser", AdminOnly: true},
	{Method: "CreateGroup", AdminOnly: true},
	{Method: "DeleteGroup", AdminOnly: true},
}
