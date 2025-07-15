---
# Codacy Rules
Configuration for AI behavior when interacting with Codacy's MCP Server

## After ANY successful `edit_file` or `reapply` operation
- YOU MUST IMMEDIATELY run the `codacy_cli_analyze` tool from Codacy's MCP Server for each file that was edited, with:
 - `rootPath`: set to the workspace path
 - `file`: set to the path of the edited file
 - `tool`: leave empty or unset
- If any issues are found in the new edits, propose and apply fixes for them.
- > NOTE: Failure to follow this rule is considered a critical error.

## CRITICAL: Dependencies and Security Checks
- IMMEDIATELY after ANY of these actions:
 - Running npm/yarn/pnpm install
 - Adding dependencies to package.json
 - Adding requirements to requirements.txt
 - Adding dependencies to pom.xml
 - Adding dependencies to build.gradle
 - Any other package manager operations
- You MUST run the `codacy_cli_analyze` tool with:
 - `rootPath`: set to the workspace path
 - `tool`: set to "trivy"
 - `file`: leave empty or unset
- If any vulnerabilities are found because of the newly added packages:
 - Stop all other operations
 - Propose and apply fixes for the security issues
 - Only continue with the original task after security issues are resolved
- EXAMPLE:
 - After: npm install react-markdown
 - Do: Run codacy_cli_analyze with trivy
 - Before: Continuing with any other tasks

## General
- Repeat the relevant steps for each modified file.
- "Propose fixes" means to both suggest and, if possible, automatically apply the fixes.
- You MUST NOT wait for the user to ask for analysis or remind you to run the tool.
- Do not run `codacy_cli_analyze` looking for changes in duplicated code or code complexity metrics.
- Do not run `codacy_cli_analyze` looking for changes in code coverage.
- Do not try to manually install Codacy CLI using either brew, npm, npx, or any other package manager.
- If the Codacy CLI is not installed, just run the `codacy_cli_analyze` tool from Codacy's MCP Server.
- When calling `codacy_cli_analyze`, only send provider, organization and repository if the project is a git repository.

## Whenever a call to a Codacy tool that uses `repository` or `organization` as a parameter returns a 404 error
- Offer to run the `codacy_setup_repository` tool to add the repository to Codacy
- If the user accepts, run the `codacy_setup_repository` tool
- Do not ever try to run the `codacy_setup_repository` tool on your own
- After setup, immediately retry the action that failed (only retry once)
---