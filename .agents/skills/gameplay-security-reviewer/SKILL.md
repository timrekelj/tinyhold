# Skill: gameplay-security-reviewer

## Purpose

Use this skill to review Tinyhold gameplay, networking, and persistence changes for security risks, cheat vectors, trust-boundary violations, and abuse cases.

## When To Use

- Reviewing authentication, connection setup, player identity, or session handling.
- Reviewing client commands, inventory actions, movement, combat, crafting, trading, or persistence.
- Evaluating whether a client can forge, replay, spam, race, or tamper with gameplay messages.
- Assessing WebSocket/API exposure and server-side validation.

## Approach

1. Treat the client as untrusted; all gameplay authority belongs on the server.
2. Verify every client command is authenticated, authorized, rate-limited where needed, and validated against server state.
3. Look for replay, duplication, race, desync, resource exhaustion, and persistence corruption risks.
4. Prefer concrete exploit scenarios and server-side mitigations over broad security advice.
5. Call out missing tests or observability when a risk is difficult to verify.

## References

- Godot multiplayer authentication: https://docs.godotengine.org/en/latest/tutorials/networking/high_level_multiplayer.html#authentication
- OWASP Cheat Sheet Series: https://cheatsheetseries.owasp.org/
- OWASP Web Security Testing Guide: https://owasp.org/www-project-web-security-testing-guide/
