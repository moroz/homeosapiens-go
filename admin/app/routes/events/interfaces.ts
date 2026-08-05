import type { components } from "~/lib/api-types";
import { datetimeLocalValueToISO8601 } from "~/lib/time";
import type { PatchEventInput } from "~/hooks";

export interface EventFormValues {
  titleEn: string;
  titlePl: string;
  subtitleEn: string;
  subtitlePl: string;
  slug: string;
  eventType: string;
  startsAt: string;
  endsAt: string;
  isVirtual: boolean;
  venueNameEn: string;
  venueNamePl: string;
  venueStreet: string;
  venueCityEn: string;
  venueCityPl: string;
  venuePostalCode: string;
  venueCountryCode: string;
  isFree: boolean;
  price: string;
  currency: string;
  hostIds: string[];
}

/** Empty strings are sent as `null` for nullable fields, matching how the server treats "not provided". */
function blankToNull(value: string): string | null {
  return value.trim() === "" ? null : value;
}

export type EventInput = components["schemas"]["EventInput"];

export function toEventInput(values: EventFormValues): EventInput {
  return {
    titleEn: values.titleEn,
    titlePl: values.titlePl,
    subtitleEn: blankToNull(values.subtitleEn),
    subtitlePl: blankToNull(values.subtitlePl),
    slug: values.slug,
    eventType: values.eventType,
    startsAt: datetimeLocalValueToISO8601(values.startsAt),
    endsAt: datetimeLocalValueToISO8601(values.endsAt),
    isVirtual: values.isVirtual,
    venueNameEn: values.isVirtual ? null : blankToNull(values.venueNameEn),
    venueNamePl: values.isVirtual ? null : blankToNull(values.venueNamePl),
    venueStreet: values.isVirtual ? null : blankToNull(values.venueStreet),
    venueCityEn: values.isVirtual ? null : blankToNull(values.venueCityEn),
    venueCityPl: values.isVirtual ? null : blankToNull(values.venueCityPl),
    venuePostalCode: values.isVirtual ? null : blankToNull(values.venuePostalCode),
    venueCountryCode: values.isVirtual ? null : blankToNull(values.venueCountryCode),
    price: values.isFree ? null : blankToNull(values.price),
    currency: values.isFree ? null : values.currency,
    hostIds: values.hostIds,
  };
}

export function toPatchEventInput(values: EventFormValues): PatchEventInput {
  let partial: PatchEventInput = {
    titleEn: values.titleEn,
    titlePl: values.titlePl,
    subtitleEn: blankToNull(values.subtitleEn),
    subtitlePl: blankToNull(values.subtitlePl),
    slug: values.slug,
    startsAt: datetimeLocalValueToISO8601(values.startsAt),
    endsAt: datetimeLocalValueToISO8601(values.endsAt),
    isVirtual: values.isVirtual,
    hostIds: values.hostIds,
  };

  if (!values.isVirtual) {
    partial = {
      ...partial,
      venueNameEn: blankToNull(values.venueNameEn),
      venueNamePl: blankToNull(values.venueNamePl),
      venueStreet: blankToNull(values.venueStreet),
      venueCityEn: blankToNull(values.venueCityEn),
      venueCityPl: blankToNull(values.venueCityPl),
    };
  }

  if (!values.isFree) {
    partial = {
      ...partial,
      price: blankToNull(values.price),
      currency: values.currency,
    };
  }

  return partial;
}
