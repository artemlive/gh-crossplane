package domain

// RepositoriesGroup is the root resource definition
// matching the YAML structure of the CRD.
type RepositoriesGroup struct {
	APIVersion string                `yaml:"apiVersion"`
	Kind       string                `yaml:"kind"`
	Metadata   Metadata              `yaml:"metadata"`
	Spec       RepositoriesGroupSpec `yaml:"spec"`
}

type Metadata struct {
	Name   string            `yaml:"name"`
	Labels map[string]string `yaml:"labels,omitempty"`
}

type RepositoriesGroupSpec struct {
	DeletionPolicy           string        `yaml:"deletionPolicy,omitempty" ui:"type=text,label=Deletion Policy"`
	ManagementPolicies       []string      `yaml:"managementPolicies,omitempty"` // omit: complex
	Repositories             []Repository  `yaml:"repositories" ui:"type=repository,label=Repositories"`
	Permissions              []Permission  `yaml:"permissions,omitempty"`         // omit: complex
	Topics                   []string      `yaml:"topics,omitempty"`              // omit: complex
	Protections              []Protection  `yaml:"protections,omitempty"`         // omit: complex
	SecurityAndAnalysis      []SecAnalysis `yaml:"securityAndAnalysis,omitempty"` // omit: complex
	DefaultBranch            string        `yaml:"defaultBranch,omitempty" ui:"type=text,label=Default Branch"`
	Visibility               string        `yaml:"visibility,omitempty" ui:"type=text,label=Visibility"`
	HasIssues                *bool         `yaml:"hasIssues,omitempty" ui:"type=checkbox,label=Has Issues"`
	HasDownloads             *bool         `yaml:"hasDownloads,omitempty" ui:"type=checkbox,label=Has Downloads"`
	HasWiki                  *bool         `yaml:"hasWiki,omitempty" ui:"type=checkbox,label=Has Wiki"`
	AllowAutoMerge           *bool         `yaml:"allowAutoMerge,omitempty" ui:"type=checkbox,label=Allow Auto-Merge"`
	AllowSquashMerge         *bool         `yaml:"allowSquashMerge,omitempty" ui:"type=checkbox,label=Allow Squash Merge"`
	DeleteBranchOnMerge      *bool         `yaml:"deleteBranchOnMerge,omitempty" ui:"type=checkbox,label=Delete Branch on Merge"`
	AutoInit                 *bool         `yaml:"autoInit,omitempty" ui:"type=checkbox,label=Auto Init"`
	ArchiveOnDestroy         *bool         `yaml:"archiveOnDestroy,omitempty" ui:"type=checkbox,label=Archive on Destroy"`
	HasDiscussions           *bool         `yaml:"hasDiscussions,omitempty" ui:"type=checkbox,label=Has Discussions"`
	AllowUpdateBranch        *bool         `yaml:"allowUpdateBranch,omitempty" ui:"type=checkbox,label=Allow Update Branch"`
	AllowMergeCommit         *bool         `yaml:"allowMergeCommit,omitempty" ui:"type=checkbox,label=Allow Merge Commit"`
	AllowRebaseMerge         *bool         `yaml:"allowRebaseMerge,omitempty" ui:"type=checkbox,label=Allow Rebase Merge"`
	IsTemplate               *bool         `yaml:"isTemplate,omitempty" ui:"type=checkbox,label=Is Template"`
	MergeCommitMessage       string        `yaml:"mergeCommitMessage,omitempty" ui:"type=text,label=Merge Commit Message"`
	MergeCommitTitle         string        `yaml:"mergeCommitTitle,omitempty" ui:"type=text,label=Merge Commit Title"`
	SquashMergeCommitMessage string        `yaml:"squashMergeCommitMessage,omitempty" ui:"type=text,label=Squash Commit Message"`
	SquashMergeCommitTitle   string        `yaml:"squashMergeCommitTitle,omitempty" ui:"type=text,label=Squash Commit Title"`
	VulnerabilityAlerts      *bool         `yaml:"vulnerabilityAlerts,omitempty" ui:"type=checkbox,label=Vulnerability Alerts"`
	AutolinkReferences       []AutolinkRef `yaml:"autolinkReferences,omitempty"` // omit: complex
}

type Repository struct {
	Name                string        `yaml:"name" ui:"type=text,label=Name"`
	Description         string        `yaml:"description,omitempty" ui:"type=text,label=Description"`
	Permissions         []Permission  `yaml:"permissions,omitempty"` // omit: complex
	Topics              []string      `yaml:"topics,omitempty"`      // omit: complex
	Archived            *bool         `yaml:"archived,omitempty" ui:"type=checkbox,label=Archived"`
	Visibility          string        `yaml:"visibility,omitempty" ui:"type=text,label=Visibility"`
	DefaultBranch       string        `yaml:"defaultBranch,omitempty" ui:"type=text,label=Default Branch"`
	AllowAutoMerge      *bool         `yaml:"allowAutoMerge,omitempty" ui:"type=checkbox,label=Allow Auto-Merge"`
	DeleteBranchOnMerge *bool         `yaml:"deleteBranchOnMerge,omitempty" ui:"type=checkbox,label=Delete Branch on Merge"`
	SecurityAndAnalysis []SecAnalysis `yaml:"securityAndAnalysis,omitempty"` // omit: complex
	Protections         []Protection  `yaml:"protections,omitempty"`         // omit: complex
}

type Permission struct {
	Team         string `yaml:"team,omitempty"`
	Collaborator string `yaml:"collaborator,omitempty"`
	Permission   string `yaml:"permission"`
}

type SecAnalysis struct {
	AdvancedSecurity             []Status `yaml:"advancedSecurity,omitempty"`
	SecretScanning               []Status `yaml:"secretScanning,omitempty"`
	SecretScanningPushProtection []Status `yaml:"secretScanningPushProtection,omitempty"`
}

type Status struct {
	Status string `yaml:"status"`
}

type Protection struct {
	Name                          string        `yaml:"name"`
	Pattern                       string        `yaml:"pattern"`
	EnforceAdmins                 bool          `yaml:"enforceAdmins,omitempty"`
	RequireConversationResolution bool          `yaml:"requireConversationResolution,omitempty"`
	RequireSignedCommits          bool          `yaml:"requireSignedCommits,omitempty"`
	RequiredStatusChecks          []StatusCheck `yaml:"requiredStatusChecks,omitempty"`
	RequiredPullRequestReviews    []PRReview    `yaml:"requiredPullRequestReviews,omitempty"`
}

type StatusCheck struct {
	Strict   bool     `yaml:"strict"`
	Contexts []string `yaml:"contexts"`
}

type PRReview struct {
	RequireCodeOwnerReviews      bool     `yaml:"requireCodeOwnerReviews"`
	DismissStaleReviews          bool     `yaml:"dismissStaleReviews"`
	RestrictDismissals           bool     `yaml:"restrictDismissals,omitempty"`
	RequiredApprovingReviewCount int      `yaml:"requiredApprovingReviewCount"`
	DismissalRestrictions        []string `yaml:"dismissalRestrictions,omitempty"`
}

type AutolinkRef struct {
	Name              string `yaml:"name"`
	IsAlphanumeric    *bool  `yaml:"isAlphanumeric,omitempty"`
	KeyPrefix         string `yaml:"keyPrefix"`
	TargetUrlTemplate string `yaml:"targetUrlTemplate"`
}
