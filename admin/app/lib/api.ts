import createClient, { type Middleware } from "openapi-fetch";

import type { paths } from "~/lib/api-types";

const API_BASE = "/api/admin";

export class ApiError extends Error {
  readonly status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

/**
 * Throw {@link ApiError} on non-2xx so React Query treats it as an error.
 */
const throwOnError: Middleware = {
  async onResponse({ response }) {
    if (!response.ok) {
      throw new ApiError(
        response.status,
        `Request to ${new URL(response.url).pathname} failed with status ${response.status}`,
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
