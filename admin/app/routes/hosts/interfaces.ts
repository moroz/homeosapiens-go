import type { HostInput } from "~/hooks";

export interface HostFormValues {
  salutation: string;
  givenName: string;
  familyName: string;
  country: string;
}

/** Empty strings are sent as `null` for nullable fields, matching how the server treats "not provided". */
function blankToNull(value: string): string | null {
  return value.trim() === "" ? null : value.trim();
}

export function toHostInput(values: HostFormValues): HostInput {
  return {
    salutation: blankToNull(values.salutation),
    givenName: values.givenName,
    familyName: values.familyName,
    country: blankToNull(values.country),
  };
}

/** Human-readable labels for the salutation message keys defined in `i18n/*.json`. */
export const SALUTATION_LABELS: Record<string, string> = {
  "common.hosts.salutation.dr": "Dr.",
};
