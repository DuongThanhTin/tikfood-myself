const allowedCommands = new Set(["git", "go", "npm", "pnpm", "yarn"]);
const blockedArgs = ["reset", "--hard", "clean", "-fdx", "push --force", "push origin main", "push origin master"];
const allowedGitSubcommands = new Set([
  "clone",
  "checkout",
  "diff",
  "status",
  "rev-parse",
  "remote",
  "branch"
]);

export function assertAllowedCommand(command: string, args: string[]): void {
  if (!allowedCommands.has(command)) {
    throw new Error(`Command is not allowlisted: ${command}`);
  }

  if (command === "git" && !allowedGitSubcommands.has(args[0] || "")) {
    throw new Error(`Git subcommand is not allowlisted: ${args[0] || "(empty)"}`);
  }

  const joined = args.join(" ");
  for (const blocked of blockedArgs) {
    if (joined.includes(blocked)) {
      throw new Error(`Command arguments are blocked: ${command} ${joined}`);
    }
  }
}
