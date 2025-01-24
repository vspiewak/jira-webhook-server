# JIRA WebHook Server


## Launch Tailscale
    
    tailscale funnel 3000


## Create Jira WebHook


    # List WebHooks

    curl -v -XGET \
    -H 'Accept: application/json' \
    -H 'Content-Type: application/json' \
    -H "Authorization: Basic $JIRA_TOKEN" \
    "$JIRA_URL/rest/webhooks/1.0/webhook"

    # Create WebHook

    curl -v \
    -H 'Accept: application/json' \
    -H 'Content-Type: application/json' \
    -H "Authorization: Basic $JIRA_TOKEN" \
    "$JIRA_URL/rest/webhooks/1.0/webhook" -d '
    {
        "name": "KANBAN Issues Webhook",
        "description": "KANBAN Issues Webhook",
        "url": "https://<change_me_for_the_tailscaledns.ts.net>/webhooks/jira/projects/{project.id}/issues/{issue.id}/on-event",
        "events": [
            "jira:issue_created",
            "jira:issue_updated",
            "jira:issue_deleted"
        ],
        "filters": {
            "issue-related-events-section": "Project = KANBAN"
        },
        "excludeBody": false,
        "secret": "my_secret_key"
    }
    '