export const DEFAULT_TIME_ZONE = "Europe/Warsaw";

export function formatInstant(iso: string) {
  return Temporal.Instant.from(iso).toLocaleString("en-GB", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  });
}

export function datetimeLocalValueToISO8601(formValue: string): string {
  return Temporal.PlainDateTime.from(formValue)
    .toZonedDateTime(DEFAULT_TIME_ZONE)
    .toInstant()
    .toJSON();
}

export function ISO8601ToDatetimeLocalValue(iso: string): string {
  return Temporal.Instant.from(iso).toZonedDateTimeISO(DEFAULT_TIME_ZONE).toJSON().slice(0, 16);
}
