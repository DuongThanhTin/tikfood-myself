import { readdir, stat } from "node:fs/promises";
import { join, relative } from "node:path";
import { readFile } from "./readFile.js";
import { assertNotSecretPath } from "../guards/secretGuard.js";

export type RepoContextSummary = {
  docsRead: string[];
  safeFilesRead: string[];
  skippedSensitivePaths: string[];
  missingExpectedFiles: string[];
  detected: {
    rootPackageJson: boolean;
    goMod: boolean;
    aiAgentYaml: boolean;
    docsDirectory: boolean;
    backendDirectory: boolean;
    frontendDirectory: boolean;
  };
  snippets: Array<{
    path: string;
    content: string;
    truncated: boolean;
  }>;
};

const maxFileBytes = 12_000;
const maxDocsFiles = 20;
const skipDirs = new Set([".git", "node_modules", "dist", "build", "coverage", ".next"]);
const expectedRootFiles = ["README.md", ".ai-agent.yaml", "package.json", "go.mod"];

function isMarkdownOrConfig(path: string): boolean {
  return /\.(md|markdown|yaml|yml|json|mod)$/i.test(path);
}

function isSensitivePath(path: string): boolean {
  try {
    assertNotSecretPath(path);
    return false;
  } catch {
    return true;
  }
}

async function exists(path: string): Promise<boolean> {
  try {
    await stat(path);
    return true;
  } catch {
    return false;
  }
}

async function readSnippet(root: string, path: string): Promise<{ content: string; truncated: boolean }> {
  const content = await readFile(join(root, path));
  if (Buffer.byteLength(content, "utf8") <= maxFileBytes) {
    return { content, truncated: false };
  }

  return {
    content: content.slice(0, maxFileBytes),
    truncated: true
  };
}

async function collectDocs(root: string): Promise<string[]> {
  const docsRoot = join(root, "docs");
  if (!(await exists(docsRoot))) {
    return [];
  }

  const results: string[] = [];

  async function visit(dir: string): Promise<void> {
    if (results.length >= maxDocsFiles) {
      return;
    }

    const entries = await readdir(dir, { withFileTypes: true });
    for (const entry of entries) {
      if (results.length >= maxDocsFiles) {
        return;
      }

      if (skipDirs.has(entry.name)) {
        continue;
      }

      const absolutePath = join(dir, entry.name);
      const relativePath = relative(root, absolutePath);

      if (entry.isDirectory()) {
        await visit(absolutePath);
      } else if (entry.isFile() && isMarkdownOrConfig(relativePath)) {
        results.push(relativePath);
      }
    }
  }

  await visit(docsRoot);
  return results.sort();
}

export async function readRepoContext(root: string): Promise<RepoContextSummary> {
  const skippedSensitivePaths: string[] = [];
  const missingExpectedFiles: string[] = [];
  const snippets: RepoContextSummary["snippets"] = [];

  const docsFiles = await collectDocs(root);
  const candidateFiles = [...expectedRootFiles, ...docsFiles];
  const uniqueFiles = [...new Set(candidateFiles)];

  for (const path of uniqueFiles) {
    if (isSensitivePath(path)) {
      skippedSensitivePaths.push(path);
      continue;
    }

    if (!(await exists(join(root, path)))) {
      if (expectedRootFiles.includes(path)) {
        missingExpectedFiles.push(path);
      }
      continue;
    }

    try {
      const snippet = await readSnippet(root, path);
      snippets.push({
        path,
        content: snippet.content,
        truncated: snippet.truncated
      });
    } catch {
      skippedSensitivePaths.push(path);
    }
  }

  return {
    docsRead: docsFiles,
    safeFilesRead: snippets.map((snippet) => snippet.path),
    skippedSensitivePaths,
    missingExpectedFiles,
    detected: {
      rootPackageJson: await exists(join(root, "package.json")),
      goMod: await exists(join(root, "go.mod")),
      aiAgentYaml: await exists(join(root, ".ai-agent.yaml")),
      docsDirectory: await exists(join(root, "docs")),
      backendDirectory: await exists(join(root, "backend")),
      frontendDirectory: await exists(join(root, "frontend"))
    },
    snippets
  };
}
