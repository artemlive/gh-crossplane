package field

type FieldGroup struct {
	TabName     string
	FieldPaths  []string // Dot-separated, e.g., "Spec.Visibility"
	GroupLevel  bool     // true = group, false = repo
	Description string   // optional
}

var FieldGroups = []FieldGroup{
	{
		TabName: "Group: General",
		FieldPaths: []string{
			"Spec.Visibility",
			"Spec.DefaultBranch",
			"Spec.Topics",
			"Spec.ArchiveOnDestroy",
			"Spec.AutoInit",
			"Spec.IsTemplate",
			"Spec.ManagementPolicies",
			"Spec.DeletionPolicy",
		},
		GroupLevel: true,
	},
	{
		TabName: "Group: Features",
		FieldPaths: []string{
			"Spec.HasIssues",
			"Spec.HasDownloads",
			"Spec.HasWiki",
			"Spec.HasDiscussions",
			"Spec.AllowAutoMerge",
			"Spec.AllowSquashMerge",
			"Spec.AllowMergeCommit",
			"Spec.AllowRebaseMerge",
			"Spec.AllowUpdateBranch",
			"Spec.DeleteBranchOnMerge",
			"Spec.VulnerabilityAlerts",
		},
		GroupLevel: true,
	},
	{
		TabName: "Group: Merge Messages",
		FieldPaths: []string{
			"Spec.MergeCommitMessage",
			"Spec.MergeCommitTitle",
			"Spec.SquashMergeCommitMessage",
			"Spec.SquashMergeCommitTitle",
		},
		GroupLevel: true,
	},
	{
		TabName: "Group: Protections",
		FieldPaths: []string{
			"Spec.Protections",
		},
		GroupLevel: true,
	},
	{
		TabName: "Group: Security",
		FieldPaths: []string{
			"Spec.SecurityAndAnalysis",
		},
		GroupLevel: true,
	},
	{
		TabName: "Group: Autolinks",
		FieldPaths: []string{
			"Spec.AutolinkReferences",
		},
		GroupLevel: true,
	},
	{
		TabName: "Group: Permissions",
		FieldPaths: []string{
			"Spec.Permissions",
		},
		GroupLevel: true,
	},
	{
		TabName: "Repositories",
		FieldPaths: []string{
			"Spec.Repositories",
		},
		GroupLevel: false,
	},
}

var RepoEditableFields = []string{
	"Name",
	"Description",
	"Archived",
	"Visibility",
	"DefaultBranch",
	"AllowAutoMerge",
	"DeleteBranchOnMerge",
}
