import { readdir, readFile } from "node:fs/promises";
import { join } from "node:path";
import { assertSafePath } from "../guards/fileGuard.js";
import { assertNotSecretPath } from "../guards/secretGuard.js";

export async function searchRepo(root: string, query: string): Promise<string[]> {
  assertSafePath(root);
  const results: string[] = [];
  const lowered = query.toLowerCase();

  async function visit(dir: string): Promise<void> {
    const entries = await readdir(dir, { withFileTypes: true });
    for (const entry of entries) {
      const path = join(dir, entry.name);
      if (entry.name === "node_modules" || entry.name === ".git" || entry.name === "dist") {
        continue;
      }
      if (entry.isDirectory()) {
        await visit(path);
      } else {
        try {
          assertNotSecretPath(path);
          const content = await readFile(path, "utf8");
          if (content.toLowerCase().includes(lowered)) {
            results.push(path);
          }
        } catch {
          // Binary, unreadable, or sensitive files are skipped.
        }
      }
    }
  }

  await visit(root);
  return results;
}
