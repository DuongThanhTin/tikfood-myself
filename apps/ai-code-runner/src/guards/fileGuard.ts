import { isAbsolute, normalize } from "node:path";

const protectedFragments = [
  ".env",
  "secrets/",
  "credentials/",
  "id_rsa",
  "private_key",
  ".pem",
  ".key"
];

export function assertSafePath(path: string): void {
  const normalized = normalize(path);
  if (!normalized || normalized.includes("..")) {
    throw new Error(`Unsafe path: ${path}`);
  }
  if (isAbsolute(normalized) && !normalized.includes(process.env.RUNNER_WORKSPACE_DIR || "/runner_workspace")) {
    throw new Error(`Path is outside runner workspace: ${path}`);
  }
  for (const fragment of protectedFragments) {
    if (normalized.includes(fragment)) {
      throw new Error(`Protected path blocked: ${path}`);
    }
  }
}
