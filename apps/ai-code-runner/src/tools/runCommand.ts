import { execFile } from "node:child_process";
import { promisify } from "node:util";
import { assertAllowedCommand } from "../guards/commandGuard.js";

const execFileAsync = promisify(execFile);

export async function runCommand(command: string, args: string[], cwd: string): Promise<{ stdout: string; stderr: string }> {
  assertAllowedCommand(command, args);
  const result = await execFileAsync(command, args, { cwd, timeout: 120_000 });
  return {
    stdout: result.stdout,
    stderr: result.stderr
  };
}
