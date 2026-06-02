import { slugify } from "../utils/slugify.js";
import { prepareGitWorkspace } from "../tools/gitWorkspace.js";
import { readRepoContext } from "../tools/repoContext.js";

type FeatureRequest = {
  feature_id: string;
  repo: string;
  base_branch: string;
  title: string;
  description: string;
  acceptance_criteria: string[];
  mode: string;
};

type RunnerFailureResponse = {
  success: false;
  feature_id: string;
  repo: string;
  stage: string;
  error: string;
  logs: string;
  recommendation: string;
};

const antiGoalTerms = [
  "delivery",
  "cart",
  "order",
  "checkout",
  "payment",
  "booking",
  "reservation",
  "in-app chat",
  "follow graph",
  "creator monetization",
  "livestream"
];

function isObject(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function isFailureResponse(value: FeatureRequest | RunnerFailureResponse): value is RunnerFailureResponse {
  return "success" in value && value.success === false;
}

function validateFeatureRequest(value: unknown): FeatureRequest | RunnerFailureResponse {
  if (!isObject(value)) {
    return failure("unknown", "unknown", "validation", "Request body must be a JSON object.");
  }

  const required = ["feature_id", "repo", "base_branch", "title", "description", "mode"];
  for (const field of required) {
    if (typeof value[field] !== "string" || !String(value[field]).trim()) {
      return failure(
        String(value.feature_id || "unknown"),
        String(value.repo || "unknown"),
        "validation",
        `Missing or invalid required string field: ${field}.`
      );
    }
  }

  if (!Array.isArray(value.acceptance_criteria) || value.acceptance_criteria.some((item) => typeof item !== "string")) {
    return failure(
      String(value.feature_id),
      String(value.repo),
      "validation",
      "Missing or invalid required array field: acceptance_criteria."
    );
  }

  return value as FeatureRequest;
}

function detectsAntiGoal(request: FeatureRequest): string[] {
  const positiveCriteria = request.acceptance_criteria.filter((criterion) => {
    const normalized = criterion.trim().toLowerCase();
    return !normalized.startsWith("no ") && !normalized.startsWith("do not ");
  });

  const text = [
    request.title,
    request.description,
    ...positiveCriteria
  ].join(" ").toLowerCase();

  return antiGoalTerms.filter((term) => text.includes(term));
}

function failure(featureId: string, repo: string, stage: string, error: string): RunnerFailureResponse {
  return {
    success: false,
    feature_id: featureId,
    repo,
    stage,
    error,
    logs: "",
    recommendation: "Review the request or complete the TODO runner implementation before expecting automated PR creation."
  };
}

export async function runFeatureJob(body: unknown): Promise<RunnerFailureResponse> {
  const request = validateFeatureRequest(body);
  if (isFailureResponse(request)) {
    return request;
  }

  const detectedAntiGoals = detectsAntiGoal(request);
  const branch = `ai/${request.feature_id}-${slugify(request.title)}`;

  if (detectedAntiGoals.length > 0) {
    return {
      success: false,
      feature_id: request.feature_id,
      repo: request.repo,
      stage: "policy",
      error: `Request appears to include MVP anti-goals: ${detectedAntiGoals.join(", ")}.`,
      logs: `planned_branch=${branch}`,
      recommendation: "Ask a human to review the feature scope before running automation."
    };
  }

  try {
    const workspace = await prepareGitWorkspace({
      repo: request.repo,
      baseBranch: request.base_branch,
      branch,
      featureId: request.feature_id
    });

    const context = await readRepoContext(workspace.repoPath);

    // TODO: Implement OpenAI calls, guarded edits, verification, reviewer,
    // commit, push, and success response.
    return {
      success: false,
      feature_id: request.feature_id,
      repo: request.repo,
      stage: "model",
      error: "Git workspace and repo context prepared, but OpenAI planning and code editing are not implemented yet.",
      logs: [
        `workspace=${workspace.workspaceRoot}`,
        `repo_path=${workspace.repoPath}`,
        `remote=${workspace.remoteUrl}`,
        `branch=${workspace.branch}`,
        `base_commit=${workspace.commitSha}`,
        `status=${workspace.status || "clean"}`,
        `docs_read=${context.docsRead.join(",") || "none"}`,
        `safe_files_read=${context.safeFilesRead.join(",") || "none"}`,
        `missing_expected_files=${context.missingExpectedFiles.join(",") || "none"}`,
        `skipped_sensitive_paths=${context.skippedSensitivePaths.join(",") || "none"}`,
        `detected=${JSON.stringify(context.detected)}`
      ].join("\n"),
      recommendation: "Implement model-driven edit planning with the collected repo context next."
    };
  } catch (error) {
    return {
      success: false,
      feature_id: request.feature_id,
      repo: request.repo,
      stage: "git_workspace",
      error: error instanceof Error ? error.message : "Failed to prepare Git workspace.",
      logs: `planned_branch=${branch}`,
      recommendation: "Check repository URL, GitHub authentication, base_branch, and runner workspace permissions."
    };
  }
}
