package structs

import (
	"time"
)

type RepositoryQuery struct {
	Repository Repository `graphql:"repository(owner: $owner, name: $name)"`
}

type Repository struct {
	Name                          string
	Url                           string
	Id                            string
	AutoMergeAllowed              bool
	DeleteBranchOnMerge           bool
	RebaseMergeAllowed            bool
	MergeCommitAllowed            bool
	HasDiscussionsEnabled         bool
	HasIssuesEnabled              bool
	HasWikiEnabled                bool
	HasProjectsEnabled            bool
	HasVulnerabilityAlertsEnabled bool
	IsArchived                    bool
	IsDisabled                    bool
	IsFork                        bool
	IsLocked                      bool
	IsMirror                      bool
	IsPrivate                     bool
	IsTemplate                    bool
	StargazerCount                int
	SquashMergeAllowed            bool
	UpdatedAt                     time.Time
	DefaultBranchRef              struct {
		Name                 string
		BranchProtectionRule BranchProtectionRuleQuery `graphql:"branchProtectionRule"`
	} `graphql:"defaultBranchRef"`
	VulnerabilityAlerts struct {
		TotalCount int
	} `graphql:"vulnerabilityAlerts"`
	PullRequests struct {
		TotalCount int
	} `graphql:"pullRequests(states: OPEN)"`
}

type BranchProtectionRuleQuery struct {
	AllowsDeletions                bool
	AllowsForcePushes              bool
	DismissesStaleReviews          bool
	IsAdminEnforced                bool
	RequiredApprovingReviewCount   int
	RequiresApprovingReviews       bool
	RequiresCodeOwnerReviews       bool
	RequiresCommitSignatures       bool
	RequiresConversationResolution bool
	RequiresLinearHistory          bool
	RequiresStatusChecks           bool
	RequiresStrictStatusChecks     bool
	RequiresDeployments            bool
	LockBranch                     bool
	RestrictsPushes                bool
	RestrictsReviewDismissals      bool
	RequireLastPushApproval        bool
}
