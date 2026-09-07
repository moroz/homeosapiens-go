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
  /** Only editable for youtube-provider videos; ignored (and left blank) otherwise. */
  youtubeId: string;
}

/**
 * Extracts a YouTube video ID from a pasted watch/share/embed/shorts URL, or
 * returns the input unchanged (trimmed) if it doesn't look like one of those,
 * so a bare ID still works.
 */
export function parseYoutubeId(input: string): string {
  const trimmed = input.trim();

  try {
    const url = new URL(trimmed);
    if (url.hostname === "youtu.be") {
      return url.pathname.slice(1) || trimmed;
    }
    if (url.hostname.endsWith("youtube.com")) {
      if (url.pathname === "/watch") {
        return url.searchParams.get("v") ?? trimmed;
      }
      const match = url.pathname.match(/^\/(?:embed|shorts)\/([^/]+)/);
      if (match) return match[1];
    }
  } catch {
    // Not a URL, fall through to returning the trimmed input as-is.
  }

  return trimmed;
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
    youtubeId: blankToNull(values.youtubeId),
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

/** A human label for a CloudFront video source's content type. */
export function sourceLabel(contentType: string): string {
  switch (contentType) {
    case "application/vnd.apple.mpegurl":
      return "HLS manifest (.m3u8)";
    case "video/mp4":
      return "MP4 fallback";
    default:
      return contentType;
  }
}
