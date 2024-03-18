package structs

type OrganizationQuery struct {
	Organization struct {
		Id           string
		Login        string
		Url          string
		Repositories struct {
			Edges []struct {
				Node struct {
					Name string
					Url  string
				} `graphql:"node"`
			} `graphql:"edges"`
		} `graphql:"repositories(first: $first)"`
	} `graphql:"organization(login: $login)"`
}

// type OrganizationQuery struct {
// 	Organization struct {
// 		Id           string
// 		Login        string
// 		Url          string
// 		Repositories struct {
// 			Edges []struct {
// 				Node RepositoryQuery `graphql:"node"`
// 			} `graphql:"edges"`
// 		} `graphql:"repositories(first: $first)"`
// 	} `graphql:"organization(login: $login)"`
// }
