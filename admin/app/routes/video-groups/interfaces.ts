import type { PatchVideoGroupInput, VideoGroupInput } from "~/hooks";

export interface VideoGroupFormValues {
  titleEn: string;
  titlePl: string;
  slug: string;
  /** A free group has no product at all, so pricing fields are hidden for one. */
  isFree: boolean;
  price: string;
  currency: string;
  /** Video IDs in playback order. Saved through its own endpoint, not with the metadata. */
  videoIds: string[];
}

/** Empty strings are sent as `null` for nullable fields, matching how the server treats "not provided". */
function blankToNull(value: string): string | null {
  return value.trim() === "" ? null : value;
}

export function toVideoGroupInput(values: VideoGroupFormValues): VideoGroupInput {
  return {
    titleEn: values.titleEn,
    titlePl: values.titlePl,
    slug: values.slug,
    price: values.isFree ? null : blankToNull(values.price),
    currency: values.isFree ? null : values.currency,
  };
}

export function toPatchVideoGroupInput(values: VideoGroupFormValues): PatchVideoGroupInput {
  const partial: PatchVideoGroupInput = {
    titleEn: values.titleEn,
    titlePl: values.titlePl,
    slug: values.slug,
  };

  // A group that already has a product keeps it even at price zero, so making a
  // group free is expressed as a zero price rather than as a cleared product.
  if (!values.isFree) {
    return { ...partial, price: blankToNull(values.price), currency: values.currency };
  }

  return partial;
}
