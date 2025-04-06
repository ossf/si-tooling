package security_insights_spec

// place the yaml input here with "-l"
input: _

// also take a project name as a tag input "-t project=foo"
project_name: string @tag(project)

// validate the input against the v1 schema
input: #v1

// helpful defaults
default_administrators: [
    {
        name: "TODO: replace with actual administrators"
        affiliation: "Foo"
        email: "joe.bob@email.com"
        social: "https://bsky.com/joebob"
        primary: true
    }
]
default_repositories: [
    {
        name: "Foo"
        url: "https://my.vcs/foobar/foo"
        comment: "TODO: replace with actual repository(ies)"
    }
]
default_repository: {
    url: "https://TODO/your/project"
    status: "active"
    "accepts-change-request": *(input."contribution-policy"."accepts-pull-requests") | true
    "accepts-automated-change-request": *(input."contribution-policy"."accepts-automated-pull-requests") | true
    "core-team": [
        {
            name: "TODO: replace with actual core team members"
            affiliation: "Foo Bar"
            email: "alicewhite@email.com"
            social: "https://bsky.com/alicewhite"
            primary: true
        }
    ]
    license: {
        url: *(input.header.license) | "https://foo.bar/LICENSE"
        expression: "TODO: replace with actual license SPDX expression"
    }
    security: {
        assessments: {
            self: {
                comment: "Self assessment has not yet been completed."
            }
        }
    }
}

security_contact: {
    for k,v in input."security-contacts" {
        if v.primary && v.type == "email" {
            name: v.value
            primary: v.primary
            email: v.value
        }
    }
}

// convert the input to the v2 schema
output: header: {
    "schema-version": "2.0.0"
    "last-updated": input.header."last-updated"
    "last-reviewed": input.header."last-reviewed"
    url: input.header."project-url"
}
output: {
    project: {
        name: project_name
        homepage: input.header."project-url"
        roadmap: input."project-lifecycle".roadmap
        "vulnerability-reporting": {
            "reports-accepted": input."vulnerability-reporting"."accepts-vulnerability-reports"
            "bug-bounty-available": *(input."vulnerability-reporting"."bug-bounty-available") | false
            if input."vulnerability-reporting"."security-policy" != "" {
                "security-policy": input."vulnerability-reporting"."security-policy"
            }
            contact: security_contact
        },
        administrators: default_administrators,
        repositories: default_repositories,
        documentation: {
            if input."contribution-policy"."code-of-conduct" != "" {
                "code-of-conduct": input."contribution-policy"."code-of-conduct"
            }
        },
    },
    repository: default_repository
}

output: #v2