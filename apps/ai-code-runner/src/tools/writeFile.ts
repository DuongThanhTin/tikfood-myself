import { mkdir, writeFile as fsWriteFile } from "node:fs/promises";
import { dirname } from "node:path";
import { assertSafePath } from "../guards/fileGuard.js";
import { assertNotSecretPath } from "../guards/secretGuard.js";

export async function writeFile(path: string, content: string): Promise<void> {
  assertSafePath(path);
  assertNotSecretPath(path);
  await mkdir(dirname(path), { recursive: true });
  await fsWriteFile(path, content, "utf8");
}
