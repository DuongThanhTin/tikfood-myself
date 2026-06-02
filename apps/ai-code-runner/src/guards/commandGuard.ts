const allowedCommands = new Set(["git", "go", "npm", "pnpm", "yarn"]);
const blockedArgs = ["reset", "--hard", "clean", "-fdx", "push --force", "push origin main", "push origin master"];

export function assertAllowedCommand(command: string, args: string[]): void {
  if (!allowedCommands.has(command)) {
    throw new Error(`Command is not allowlisted: ${command}`);
  }

  const joined = args.join(" ");
  for (const blocked of blockedArgs) {
    if (joined.includes(blocked)) {
      throw new Error(`Command arguments are blocked: ${command} ${joined}`);
    }
  }
}
