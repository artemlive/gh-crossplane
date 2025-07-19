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
	DeletionPolicy           string        `yaml:"deletionPolicy,omitempty"`
	ManagementPolicies       []string      `yaml:"managementPolicies,omitempty"`
	Repositories             []Repository  `yaml:"repositories"`
	Permissions              []Permission  `yaml:"permissions,omitempty"`
	Topics                   []string      `yaml:"topics,omitempty"`
	Protections              []Protection  `yaml:"protections,omitempty"`
	SecurityAndAnalysis      []SecAnalysis `yaml:"securityAndAnalysis,omitempty"`
	DefaultBranch            string        `yaml:"defaultBranch,omitempty"`
	Visibility               string        `yaml:"visibility,omitempty"`
	HasIssues                *bool         `yaml:"hasIssues,omitempty"`
	HasDownloads             *bool         `yaml:"hasDownloads,omitempty"`
	HasWiki                  *bool         `yaml:"hasWiki,omitempty"`
	AllowAutoMerge           *bool         `yaml:"allowAutoMerge,omitempty"`
	AllowSquashMerge         *bool         `yaml:"allowSquashMerge,omitempty"`
	DeleteBranchOnMerge      *bool         `yaml:"deleteBranchOnMerge,omitempty"`
	AutoInit                 *bool         `yaml:"autoInit,omitempty"`
	ArchiveOnDestroy         *bool         `yaml:"archiveOnDestroy,omitempty"`
	HasDiscussions           *bool         `yaml:"hasDiscussions,omitempty"`
	AllowUpdateBranch        *bool         `yaml:"allowUpdateBranch,omitempty"`
	AllowMergeCommit         *bool         `yaml:"allowMergeCommit,omitempty"`
	AllowRebaseMerge         *bool         `yaml:"allowRebaseMerge,omitempty"`
	IsTemplate               *bool         `yaml:"isTemplate,omitempty"`
	MergeCommitMessage       string        `yaml:"mergeCommitMessage,omitempty"`
	MergeCommitTitle         string        `yaml:"mergeCommitTitle,omitempty"`
	SquashMergeCommitMessage string        `yaml:"squashMergeCommitMessage,omitempty"`
	SquashMergeCommitTitle   string        `yaml:"squashMergeCommitTitle,omitempty"`
	VulnerabilityAlerts      *bool         `yaml:"vulnerabilityAlerts,omitempty"`
	AutolinkReferences       []AutolinkRef `yaml:"autolinkReferences,omitempty"`
}

type Repository struct {
	Name                string        `yaml:"name"`
	Description         string        `yaml:"description,omitempty"`
	Permissions         []Permission  `yaml:"permissions,omitempty"`
	Topics              []string      `yaml:"topics,omitempty"`
	Archived            *bool         `yaml:"archived,omitempty"`
	Visibility          string        `yaml:"visibility,omitempty"`
	DefaultBranch       string        `yaml:"defaultBranch,omitempty"`
	AllowAutoMerge      *bool         `yaml:"allowAutoMerge,omitempty"`
	DeleteBranchOnMerge *bool         `yaml:"deleteBranchOnMerge,omitempty"`
	SecurityAndAnalysis []SecAnalysis `yaml:"securityAndAnalysis,omitempty"`
	Protections         []Protection  `yaml:"protections,omitempty"`
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
