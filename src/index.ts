import vectors from "../session-command-vectors.json" with { type: "json" };

/** The commands a view calls to write and read the core's session index.
 *
 *  Named once here so both sides grade themselves against one file rather than against each other.
 *  A repository reading another's source is what the coupling law forbids, and a name that differs
 *  by one character produces an empty index and no error.
 *
 *  Measured 2026-09-04: a view called `session_attach` with `viewId` while the command it could run
 *  was `session.attach` with `view`. Either mismatch alone left the index empty, and neither raised:
 *  the runner answers `{ok:false}`, so a name nothing serves reads as a success to a caller that
 *  only catches. */
export interface SessionCommand {
  /** The exact name a caller runs. */
  readonly command: string;
  /** The exact parameter names that call takes. */
  readonly params: readonly string[];
  /** What the command is for, in one line. */
  readonly purpose: string;
}

export const SESSION_COMMANDS: readonly SessionCommand[] = Object.freeze(
  vectors.map((one) => Object.freeze({
    command: one.command,
    params: Object.freeze([...one.params]),
    purpose: one.purpose,
  })),
);

/** The declared command of that exact name, or undefined when the contract declares none. */
export function sessionCommand(name: string): SessionCommand | undefined {
  return SESSION_COMMANDS.find((one) => one.command === name);
}

export const CONTROL_CONTRACT = Object.freeze({
  id: "soksak-contract-control",
  version: "0.0.1",
} as const);
