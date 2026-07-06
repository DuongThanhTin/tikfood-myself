import { describe, expect, it } from "vitest";

import { assertAllowedCommand } from "./commandGuard";

describe("assertAllowedCommand", () => {
  it("allows an allowlisted git subcommand", () => {
    expect(() => assertAllowedCommand("git", ["status"])).not.toThrow();
  });

  it("rejects a command that is not on the allowlist", () => {
    expect(() => assertAllowedCommand("rm", ["-rf", "/"])).toThrow(/not allowlisted/);
  });

  it("rejects a git subcommand that is not allowlisted (e.g. push to main)", () => {
    expect(() => assertAllowedCommand("git", ["push", "origin", "main"])).toThrow(
      /Git subcommand is not allowlisted/
    );
  });
});
