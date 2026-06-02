import { runCommand } from "./runCommand.js";

export async function gitDiff(cwd: string): Promise<string> {
  const result = await runCommand("git", ["diff"], cwd);
  return result.stdout;
}
