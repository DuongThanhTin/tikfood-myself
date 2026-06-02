import { mkdir, rm } from "node:fs/promises";
import { join } from "node:path";
import { runCommand } from "./runCommand.js";

export type GitWorkspaceResult = {
  workspaceRoot: string;
  repoPath: string;
  branch: string;
  commitSha: string;
  remoteUrl: string;
  status: string;
};

type GitWorkspaceInput = {
  repo: string;
  baseBranch: string;
  branch: string;
  featureId: string;
};

function runnerWorkspaceRoot(): string {
  return process.env.RUNNER_WORKSPACE_DIR || "/runner_workspace";
}

function repoUrl(repo: string): string {
  if (repo.startsWith("https://") || repo.startsWith("git@")) {
    return repo;
  }
  return `https://github.com/${repo}.git`;
}

function safeWorkspaceName(featureId: string): string {
  return featureId.replace(/[^a-zA-Z0-9_.-]/g, "-");
}

export async function prepareGitWorkspace(input: GitWorkspaceInput): Promise<GitWorkspaceResult> {
  const workspaceRoot = runnerWorkspaceRoot();
  const jobRoot = join(workspaceRoot, safeWorkspaceName(input.featureId));
  const repoPath = join(jobRoot, "repo");
  const remoteUrl = repoUrl(input.repo);

  await rm(jobRoot, { recursive: true, force: true });
  await mkdir(jobRoot, { recursive: true });

  await runCommand("git", ["clone", "--depth", "1", "--branch", input.baseBranch, remoteUrl, repoPath], jobRoot);
  await runCommand("git", ["checkout", "-b", input.branch], repoPath);

  const commitSha = (await runCommand("git", ["rev-parse", "HEAD"], repoPath)).stdout.trim();
  const status = (await runCommand("git", ["status", "--short"], repoPath)).stdout.trim();

  return {
    workspaceRoot: jobRoot,
    repoPath,
    branch: input.branch,
    commitSha,
    remoteUrl,
    status
  };
}
