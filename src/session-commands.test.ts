// The vectors are graded in each language a consumer reads them in: a parse that succeeds in one
// and fails in the other is a file two sides disagree about.
import { describe, expect, it } from "vitest";

import { SESSION_COMMANDS, sessionCommand } from "./index";

describe("the session commands", () => {
  it("names at least one command", () => {
    expect(SESSION_COMMANDS.length).toBeGreaterThan(0);
  });

  it("names each command once, with a purpose", () => {
    const seen = new Set<string>();
    for (const one of SESSION_COMMANDS) {
      expect(one.command).toMatch(/^session\.[a-z][a-zA-Z]*$/);
      expect(one.purpose.length).toBeGreaterThan(0);
      expect(seen.has(one.command)).toBe(false);
      seen.add(one.command);
    }
  });

  it("names each parameter once per command", () => {
    for (const one of SESSION_COMMANDS) {
      expect(new Set(one.params).size).toBe(one.params.length);
      for (const param of one.params) expect(param).toMatch(/^[a-z][a-zA-Z]*$/);
    }
  });

  it("answers a declared name and refuses an undeclared one", () => {
    expect(sessionCommand("session.attach")?.params).toEqual(["session", "owner", "view"]);
    expect(sessionCommand("session_attach")).toBeUndefined();
  });
});
