import createClient, { type Middleware } from "openapi-fetch";

import type { paths } from "~/lib/api-types";

const API_BASE = "/api/admin";

export class ApiError extends Error {
  readonly status: number;
  /** Parsed JSON body of the failed response, when it was JSON (e.g. `ValidationErrors` on a 422). */
  readonly body: unknown;

  constructor(status: number, message: string, body?: unknown) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.body = body;
  }
}

/**
 * Throw {@link ApiError} on non-2xx so React Query treats it as an error.
 * The response body is read (via a clone, since the original is consumed elsewhere)
 * and attached so callers can surface e.g. 422 validation messages.
 */
const throwOnError: Middleware = {
  async onResponse({ response }) {
    if (!response.ok) {
      let body: unknown;
      try {
        body = await response.clone().json();
      } catch {
        // Non-JSON error body; leave `body` undefined.
      }
      throw new ApiError(
        response.status,
        `Request to ${new URL(response.url).pathname} failed with status ${response.status}`,
        body,
      );
    }
    return response;
  },
};

/**
 * Typed client for the same-origin `/api/admin` JSON API.
 * Paths and payloads are checked against `api-types.ts` (regen with `pnpm gen:api`).
 */
export const api = createClient<paths>({
  baseUrl: API_BASE,
  headers: { Accept: "application/json" },
});

api.use(throwOnError);

/** Narrows an {@link ApiError} body to the `ValidationErrors` shape (`{ errors: { field: message } }`). */
export function isValidationErrorBody(
  body: unknown,
): body is { errors: Record<string, string> } {
  return (
    typeof body === "object" &&
    body !== null &&
    "errors" in body &&
    typeof (body as { errors: unknown }).errors === "object"
  );
}
