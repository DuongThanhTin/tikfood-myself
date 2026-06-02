import { readFile as fsReadFile } from "node:fs/promises";
import { assertSafePath } from "../guards/fileGuard.js";
import { assertNotSecretPath } from "../guards/secretGuard.js";

export async function readFile(path: string): Promise<string> {
  assertSafePath(path);
  assertNotSecretPath(path);
  return fsReadFile(path, "utf8");
}
