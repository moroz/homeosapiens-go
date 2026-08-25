import type { UpdateVideoInput } from "~/hooks";

export interface VideoFormValues {
  titleEn: string;
  titlePl: string;
  slug: string;
  descriptionEn: string;
  descriptionPl: string;
  /** `yyyy-mm-dd`, as produced by an `<input type="date">`. */
  recordedOn: string;
  isPublic: boolean;
  hostId: string;
}

/** Empty strings are sent as `null` for nullable fields, matching how the server treats "not provided". */
function blankToNull(value: string): string | null {
  return value.trim() === "" ? null : value;
}

export function toUpdateVideoInput(values: VideoFormValues): UpdateVideoInput {
  return {
    titleEn: values.titleEn,
    titlePl: values.titlePl,
    slug: values.slug,
    descriptionEn: blankToNull(values.descriptionEn),
    descriptionPl: blankToNull(values.descriptionPl),
    recordedOn: values.recordedOn ? `${values.recordedOn}T00:00:00Z` : null,
    isPublic: values.isPublic,
    hostId: blankToNull(values.hostId),
  };
}

/** The date part of an ISO instant, i.e. what an `<input type="date">` expects. */
export function toDateInputValue(iso: string | null | undefined): string {
  return iso ? iso.slice(0, 10) : "";
}

/** Renders a duration in seconds as `h:mm:ss`, or `mm:ss` for a short video. */
export function formatDuration(seconds: number | null | undefined): string {
  if (seconds == null) return "";

  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;
  const pad = (n: number) => String(n).padStart(2, "0");

  return hours > 0 ? `${hours}:${pad(minutes)}:${pad(secs)}` : `${minutes}:${pad(secs)}`;
}
