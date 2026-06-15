package fixtures

import "time"

// recentDate returns an RFC3339 timestamp offset from the current time,
// keeping fixture dates within the default 3-month --since window.
func recentDate(daysAgo int) string {
	return time.Now().AddDate(0, 0, -daysAgo).UTC().Format(time.RFC3339)
}

// OrgValidationResponse returns a successful org validation response.
func OrgValidationResponse(org string) string {
	return `{
		"data": {
			"organization": {"login": "` + org + `"}
		}
	}`
}

// RateLimitResponse returns a healthy rate limit response.
func RateLimitResponse() string {
	return `{
		"data": {
			"rateLimit": {"remaining": 4999, "resetAt": "2099-01-01T00:00:00Z"}
		}
	}`
}

// OrgReposResponse returns a single-page repos response.
func OrgReposResponse(repoNames ...string) string {
	nodes := ""
	for i, name := range repoNames {
		if i > 0 {
			nodes += ","
		}
		nodes += `{"name": "` + name + `"}`
	}
	return `{
		"data": {
			"organization": {
				"repositories": {
					"pageInfo": {"hasNextPage": false, "endCursor": null},
					"nodes": [` + nodes + `]
				}
			},
			"rateLimit": {"remaining": 4998, "resetAt": "2099-01-01T00:00:00Z"}
		}
	}`
}

// ReviewActivityResponse returns a response with one PR containing review comments.
// Dates are generated relative to now so they always fall within the default 3-month --since window.
func ReviewActivityResponse() string {
	d0 := recentDate(30)
	d1 := recentDate(29)
	d2 := recentDate(28)
	d3 := recentDate(27)
	return `{
		"data": {
			"repository": {
				"pullRequests": {
					"pageInfo": {"hasNextPage": false, "endCursor": null},
					"nodes": [{
						"number": 42,
						"title": "Add feature X",
						"url": "https://github.com/test-org/test-repo/pull/42",
						"state": "MERGED",
						"updatedAt": "` + d0 + `",
						"author": {"login": "pr-author"},
						"reviews": {
							"nodes": [
								{"author": {"login": "reviewer-alice"}, "state": "APPROVED", "createdAt": "` + d1 + `"},
								{"author": {"login": "reviewer-bob"}, "state": "CHANGES_REQUESTED", "createdAt": "` + d2 + `"}
							]
						},
						"reviewThreads": {
							"nodes": [
								{
									"comments": {
										"nodes": [
											{
												"author": {"login": "reviewer-alice"},
												"body": "Consider using early return here to reduce nesting",
												"diffHunk": "@@ -1,3 +1,5 @@\n func foo() {\n+  if err != nil {\n+    return err\n+  }",
												"createdAt": "` + d1 + `"
											}
										]
									}
								},
								{
									"comments": {
										"nodes": [
											{
												"author": {"login": "reviewer-bob"},
												"body": "This needs error handling for the nil case",
												"diffHunk": "@@ -10,3 +12,5 @@\n resp, err := client.Do(req)",
												"createdAt": "` + d2 + `"
											}
										]
									}
								},
								{
									"comments": {
										"nodes": [
											{
												"author": {"login": "reviewer-alice"},
												"body": "Prefer const over let for immutable bindings",
												"diffHunk": "@@ -20,1 +22,1 @@\n-let x = 5;\n+const x = 5;",
												"createdAt": "` + d3 + `"
											}
										]
									}
								}
							]
						}
					}]
				}
			},
			"rateLimit": {"remaining": 4997, "resetAt": "2099-01-01T00:00:00Z"}
		}
	}`
}

// MultiReviewerResponse returns a response with 2 PRs, 3 human reviewers
// (alice: 3 comments, bob: 2, charlie: 1) plus dependabot[bot] in both
// reviews and comments. Supports testing --top N filtering, bot filtering
// on reviews, totalComments, and CSV ranking order.
// Dates are generated relative to now so they always fall within the default 3-month --since window.
func MultiReviewerResponse() string {
	d0 := recentDate(30)
	d1 := recentDate(29)
	d2 := recentDate(28)
	d3 := recentDate(27)
	d4 := recentDate(26)
	d5 := recentDate(25)
	d6 := recentDate(24)
	d7 := recentDate(23)
	return `{
		"data": {
			"repository": {
				"pullRequests": {
					"pageInfo": {"hasNextPage": false, "endCursor": null},
					"nodes": [
						{
							"number": 100,
							"title": "Big refactor",
							"url": "https://github.com/test-org/test-repo/pull/100",
							"state": "MERGED",
							"updatedAt": "` + d0 + `",
							"author": {"login": "pr-author"},
							"reviews": {
								"nodes": [
									{"author": {"login": "reviewer-alice"}, "state": "APPROVED", "createdAt": "` + d1 + `"},
									{"author": {"login": "reviewer-bob"}, "state": "CHANGES_REQUESTED", "createdAt": "` + d2 + `"},
									{"author": {"login": "reviewer-charlie"}, "state": "APPROVED", "createdAt": "` + d3 + `"},
									{"author": {"login": "dependabot[bot]"}, "state": "APPROVED", "createdAt": "` + d4 + `"}
								]
							},
							"reviewThreads": {
								"nodes": [
									{"comments": {"nodes": [{"author": {"login": "reviewer-alice"}, "body": "Use early return", "diffHunk": "@@ -1,3 +1,5 @@\n+if err != nil {", "createdAt": "` + d1 + `"}]}},
									{"comments": {"nodes": [{"author": {"login": "reviewer-alice"}, "body": "Prefer const", "diffHunk": "@@ -10,1 +10,1 @@\n-let x = 5;", "createdAt": "` + d2 + `"}]}},
									{"comments": {"nodes": [{"author": {"login": "reviewer-bob"}, "body": "Handle nil case", "diffHunk": "@@ -20,3 +22,5 @@\n resp, err := client.Do(req)", "createdAt": "` + d2 + `"}]}},
									{"comments": {"nodes": [{"author": {"login": "reviewer-charlie"}, "body": "Add docs", "diffHunk": "@@ -30,1 +32,1 @@\n+// TODO", "createdAt": "` + d3 + `"}]}},
									{"comments": {"nodes": [{"author": {"login": "dependabot[bot]"}, "body": "Dep update safe", "diffHunk": "@@ -1,1 +1,1 @@\n-v1.0\n+v1.1", "createdAt": "` + d4 + `"}]}}
								]
							}
						},
						{
							"number": 101,
							"title": "Small fix",
							"url": "https://github.com/test-org/test-repo/pull/101",
							"state": "MERGED",
							"updatedAt": "` + d5 + `",
							"author": {"login": "pr-author"},
							"reviews": {
								"nodes": [
									{"author": {"login": "reviewer-alice"}, "state": "APPROVED", "createdAt": "` + d6 + `"},
									{"author": {"login": "reviewer-bob"}, "state": "APPROVED", "createdAt": "` + d7 + `"}
								]
							},
							"reviewThreads": {
								"nodes": [
									{"comments": {"nodes": [{"author": {"login": "reviewer-alice"}, "body": "LGTM with nit", "diffHunk": "@@ -5,1 +5,1 @@\n-old\n+new", "createdAt": "` + d6 + `"}]}},
									{"comments": {"nodes": [{"author": {"login": "reviewer-bob"}, "body": "Typo here", "diffHunk": "@@ -8,1 +8,1 @@\n-colour\n+color", "createdAt": "` + d7 + `"}]}}
								]
							}
						}
					]
				}
			},
			"rateLimit": {"remaining": 4997, "resetAt": "2099-01-01T00:00:00Z"}
		}
	}`
}

// RepoNotFoundError returns a GraphQL error response for a repository
// that cannot be resolved, used to test graceful error handling.
func RepoNotFoundError(org, repo string) string {
	return `{
		"data": null,
		"errors": [{"type": "NOT_FOUND", "message": "Could not resolve to a Repository with the name '` + org + `/` + repo + `'."}]
	}`
}

// ReviewActivityWithBotResponse includes a bot reviewer and bot commenter
// alongside human reviewers, for testing bot filtering.
// Dates are generated relative to now so they always fall within the default 3-month --since window.
func ReviewActivityWithBotResponse() string {
	d0 := recentDate(30)
	d1 := recentDate(29)
	d2 := recentDate(28)
	d3 := recentDate(27)
	return `{
		"data": {
			"repository": {
				"pullRequests": {
					"pageInfo": {"hasNextPage": false, "endCursor": null},
					"nodes": [{
						"number": 42,
						"title": "Add feature X",
						"url": "https://github.com/test-org/test-repo/pull/42",
						"state": "MERGED",
						"updatedAt": "` + d0 + `",
						"author": {"login": "pr-author"},
						"reviews": {
							"nodes": [
								{"author": {"login": "reviewer-alice"}, "state": "APPROVED", "createdAt": "` + d1 + `"},
								{"author": {"login": "reviewer-bob"}, "state": "CHANGES_REQUESTED", "createdAt": "` + d2 + `"},
								{"author": {"login": "dependabot[bot]"}, "state": "APPROVED", "createdAt": "` + d3 + `"}
							]
						},
						"reviewThreads": {
							"nodes": [
								{
									"comments": {
										"nodes": [
											{
												"author": {"login": "reviewer-alice"},
												"body": "Consider using early return here to reduce nesting",
												"diffHunk": "@@ -1,3 +1,5 @@\n func foo() {\n+  if err != nil {\n+    return err\n+  }",
												"createdAt": "` + d1 + `"
											}
										]
									}
								},
								{
									"comments": {
										"nodes": [
											{
												"author": {"login": "reviewer-bob"},
												"body": "This needs error handling for the nil case",
												"diffHunk": "@@ -10,3 +12,5 @@\n resp, err := client.Do(req)",
												"createdAt": "` + d2 + `"
											}
										]
									}
								},
								{
									"comments": {
										"nodes": [
											{
												"author": {"login": "dependabot[bot]"},
												"body": "This dependency update is safe to merge",
												"diffHunk": "@@ -1,1 +1,1 @@\n-dep: v1.0.0\n+dep: v1.1.0",
												"createdAt": "` + d3 + `"
											}
										]
									}
								},
								{
									"comments": {
										"nodes": [
											{
												"author": {"login": "reviewer-alice"},
												"body": "Prefer const over let for immutable bindings",
												"diffHunk": "@@ -20,1 +22,1 @@\n-let x = 5;\n+const x = 5;",
												"createdAt": "` + d3 + `"
											}
										]
									}
								}
							]
						}
					}]
				}
			},
			"rateLimit": {"remaining": 4997, "resetAt": "2099-01-01T00:00:00Z"}
		}
	}`
}
