import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { useNavigate } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { Button } from "~/components/ui/button";
import { Checkbox } from "~/components/ui/checkbox";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";
import { Switch } from "~/components/ui/switch";
import { Textarea } from "~/components/ui/textarea";
import { ApiError, isValidationErrorBody } from "~/lib/api";
import type { components } from "~/lib/api-types";
import { useCreateEventMutation, useListHostsQuery } from "~/hooks";

type EventInput = components["schemas"]["EventInput"];

interface FormValues {
  titleEn: string;
  titlePl: string;
  subtitleEn: string;
  subtitlePl: string;
  slug: string;
  eventType: string;
  descriptionEn: string;
  descriptionPl: string;
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

const defaultValues: FormValues = {
  titleEn: "",
  titlePl: "",
  subtitleEn: "",
  subtitlePl: "",
  slug: "",
  eventType: "webinar",
  descriptionEn: "",
  descriptionPl: "",
  startsAt: "",
  endsAt: "",
  isVirtual: true,
  venueNameEn: "",
  venueNamePl: "",
  venueStreet: "",
  venueCityEn: "",
  venueCityPl: "",
  venuePostalCode: "",
  venueCountryCode: "",
  isFree: true,
  price: "",
  currency: "PLN",
  hostIds: [],
};

/** Empty strings are sent as `null` for nullable fields, matching how the server treats "not provided". */
function blankToNull(value: string): string | null {
  return value.trim() === "" ? null : value;
}

function toEventInput(values: FormValues): EventInput {
  return {
    titleEn: values.titleEn,
    titlePl: values.titlePl,
    subtitleEn: blankToNull(values.subtitleEn),
    subtitlePl: blankToNull(values.subtitlePl),
    slug: values.slug,
    eventType: values.eventType,
    descriptionEn: values.descriptionEn,
    descriptionPl: values.descriptionPl,
    startsAt: new Date(values.startsAt).toISOString(),
    endsAt: new Date(values.endsAt).toISOString(),
    isVirtual: values.isVirtual,
    venueNameEn: values.isVirtual ? null : blankToNull(values.venueNameEn),
    venueNamePl: values.isVirtual ? null : blankToNull(values.venueNamePl),
    venueStreet: values.isVirtual ? null : blankToNull(values.venueStreet),
    venueCityEn: values.isVirtual ? null : blankToNull(values.venueCityEn),
    venueCityPl: values.isVirtual ? null : blankToNull(values.venueCityPl),
    venuePostalCode: values.isVirtual ? null : blankToNull(values.venuePostalCode),
    venueCountryCode: values.isVirtual ? null : blankToNull(values.venueCountryCode),
    isFree: values.isFree,
    price: values.isFree ? null : blankToNull(values.price),
    currency: values.isFree ? null : values.currency,
    hostIds: values.hostIds,
  };
}

/** Server field names line up 1:1 with `FormValues` keys, so validation errors map straight onto form fields. */
const FORM_FIELDS = new Set<string>(Object.keys(defaultValues));

function isFormField(field: string): field is keyof FormValues {
  return FORM_FIELDS.has(field);
}

function FieldError({ message }: { message?: string }) {
  if (!message) return null;
  return <p className="text-sm text-destructive">{message}</p>;
}

function HostMultiSelect({
  value,
  onChange,
}: {
  value: string[];
  onChange: (ids: string[]) => void;
}) {
  const [search, setSearch] = useState("");
  const { data: hosts, isPending } = useListHostsQuery();

  const filtered = (hosts?.data ?? []).filter((host) =>
    `${host.givenName} ${host.familyName}`.toLowerCase().includes(search.toLowerCase()),
  );

  function toggle(id: string) {
    onChange(value.includes(id) ? value.filter((existing) => existing !== id) : [...value, id]);
  }

  return (
    <div className="flex flex-col gap-2">
      <Input
        placeholder="Search hosts…"
        value={search}
        onChange={(event) => setSearch(event.target.value)}
      />
      <div className="flex max-h-48 flex-col gap-1 overflow-y-auto rounded-md border border-input p-2">
        {isPending ? (
          <p className="text-sm text-muted-foreground">Loading hosts…</p>
        ) : filtered.length === 0 ? (
          <p className="text-sm text-muted-foreground">No hosts found.</p>
        ) : (
          filtered.map((host) => (
            <label key={host.id} className="flex items-center gap-2 rounded-sm p-1 text-sm">
              <Checkbox checked={value.includes(host.id)} onCheckedChange={() => toggle(host.id)} />
              {host.givenName} {host.familyName}
            </label>
          ))
        )}
      </div>
    </div>
  );
}

export default function NewEvent() {
  const navigate = useNavigate();
  const createEvent = useCreateEventMutation();
  const [formError, setFormError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    control,
    watch,
    setError,
    formState: { errors },
  } = useForm<FormValues>({ defaultValues });

  const isVirtual = watch("isVirtual");
  const isFree = watch("isFree");

  async function onSubmit(values: FormValues) {
    setFormError(null);
    try {
      const event = await createEvent.mutateAsync(toEventInput(values));
      navigate(`/events/${event.id}`);
    } catch (err) {
      if (err instanceof ApiError && err.status === 422 && isValidationErrorBody(err.body)) {
        for (const [field, message] of Object.entries(err.body.errors)) {
          if (isFormField(field)) {
            setError(field, { message });
          }
        }
        setFormError("Please fix the errors below.");
        return;
      }
      setFormError("Something went wrong creating the event. Please try again.");
    }
  }

  return (
    <AdminLayout title="New event">
      <form onSubmit={handleSubmit(onSubmit)} className="flex max-w-2xl flex-col gap-6">
        <FieldError message={formError ?? undefined} />

        <section className="flex flex-col gap-4">
          <h3 className="text-lg font-semibold">Basics</h3>

          <div className="grid grid-cols-2 gap-4">
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="titleEn">Title (EN)</Label>
              <Input id="titleEn" {...register("titleEn", { required: "Required" })} />
              <FieldError message={errors.titleEn?.message} />
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="titlePl">Title (PL)</Label>
              <Input id="titlePl" {...register("titlePl", { required: "Required" })} />
              <FieldError message={errors.titlePl?.message} />
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="subtitleEn">Subtitle (EN)</Label>
              <Input id="subtitleEn" {...register("subtitleEn")} />
              <FieldError message={errors.subtitleEn?.message} />
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="subtitlePl">Subtitle (PL)</Label>
              <Input id="subtitlePl" {...register("subtitlePl")} />
              <FieldError message={errors.subtitlePl?.message} />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="slug">Slug</Label>
              <Input
                id="slug"
                className="font-mono"
                placeholder="my-event-slug"
                {...register("slug", { required: "Required" })}
              />
              <FieldError message={errors.slug?.message} />
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="eventType">Event type</Label>
              <Controller
                name="eventType"
                control={control}
                render={({ field }) => (
                  <Select value={field.value} onValueChange={field.onChange}>
                    <SelectTrigger id="eventType" className="w-full">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="webinar">Webinar</SelectItem>
                      <SelectItem value="seminar">Seminar</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
              <FieldError message={errors.eventType?.message} />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="descriptionEn">Description (EN)</Label>
              <Textarea
                id="descriptionEn"
                {...register("descriptionEn", { required: "Required" })}
              />
              <FieldError message={errors.descriptionEn?.message} />
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="descriptionPl">Description (PL)</Label>
              <Textarea
                id="descriptionPl"
                {...register("descriptionPl", { required: "Required" })}
              />
              <FieldError message={errors.descriptionPl?.message} />
            </div>
          </div>
        </section>

        <section className="flex flex-col gap-4">
          <h3 className="text-lg font-semibold">Schedule</h3>
          <div className="grid grid-cols-2 gap-4">
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="startsAt">Starts at</Label>
              <Input
                id="startsAt"
                type="datetime-local"
                {...register("startsAt", { required: "Required" })}
              />
              <FieldError message={errors.startsAt?.message} />
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="endsAt">Ends at</Label>
              <Input
                id="endsAt"
                type="datetime-local"
                {...register("endsAt", { required: "Required" })}
              />
              <FieldError message={errors.endsAt?.message} />
            </div>
          </div>
        </section>

        <section className="flex flex-col gap-4">
          <div className="flex items-center gap-2.5">
            <Controller
              name="isVirtual"
              control={control}
              render={({ field }) => (
                <Switch id="isVirtual" checked={field.value} onCheckedChange={field.onChange} />
              )}
            />
            <Label htmlFor="isVirtual">This is a virtual event</Label>
          </div>

          {!isVirtual && (
            <div className="grid grid-cols-2 gap-4 rounded-md border border-input p-4">
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="venueNameEn">Venue name (EN)</Label>
                <Input
                  id="venueNameEn"
                  {...register("venueNameEn", { required: "Required for in-person events" })}
                />
                <FieldError message={errors.venueNameEn?.message} />
              </div>
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="venueNamePl">Venue name (PL)</Label>
                <Input
                  id="venueNamePl"
                  {...register("venueNamePl", { required: "Required for in-person events" })}
                />
                <FieldError message={errors.venueNamePl?.message} />
              </div>
              <div className="col-span-2 flex flex-col gap-1.5">
                <Label htmlFor="venueStreet">Street</Label>
                <Input
                  id="venueStreet"
                  {...register("venueStreet", { required: "Required for in-person events" })}
                />
                <FieldError message={errors.venueStreet?.message} />
              </div>
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="venueCityEn">City (EN)</Label>
                <Input
                  id="venueCityEn"
                  {...register("venueCityEn", { required: "Required for in-person events" })}
                />
                <FieldError message={errors.venueCityEn?.message} />
              </div>
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="venueCityPl">City (PL)</Label>
                <Input
                  id="venueCityPl"
                  {...register("venueCityPl", { required: "Required for in-person events" })}
                />
                <FieldError message={errors.venueCityPl?.message} />
              </div>
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="venuePostalCode">Postal code</Label>
                <Input
                  id="venuePostalCode"
                  {...register("venuePostalCode", { required: "Required for in-person events" })}
                />
                <FieldError message={errors.venuePostalCode?.message} />
              </div>
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="venueCountryCode">Country code</Label>
                <Input
                  id="venueCountryCode"
                  placeholder="PL"
                  maxLength={2}
                  {...register("venueCountryCode", {
                    required: "Required for in-person events",
                    minLength: { value: 2, message: "Must be 2 letters" },
                    maxLength: { value: 2, message: "Must be 2 letters" },
                  })}
                />
                <FieldError message={errors.venueCountryCode?.message} />
              </div>
            </div>
          )}
        </section>

        <section className="flex flex-col gap-4">
          <div className="flex items-center gap-2.5">
            <Controller
              name="isFree"
              control={control}
              render={({ field }) => (
                <Switch id="isFree" checked={field.value} onCheckedChange={field.onChange} />
              )}
            />
            <Label htmlFor="isFree">This event is free</Label>
          </div>

          {!isFree && (
            <div className="grid grid-cols-2 gap-4 rounded-md border border-input p-4">
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="price">Price</Label>
                <Input
                  id="price"
                  placeholder="19.99"
                  {...register("price", { required: "Required for paid events" })}
                />
                <FieldError message={errors.price?.message} />
              </div>
              <div className="flex flex-col gap-1.5">
                <Label htmlFor="currency">Currency</Label>
                <Controller
                  name="currency"
                  control={control}
                  render={({ field }) => (
                    <Select value={field.value} onValueChange={field.onChange}>
                      <SelectTrigger id="currency" className="w-full">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="PLN">PLN</SelectItem>
                        <SelectItem value="EUR">EUR</SelectItem>
                      </SelectContent>
                    </Select>
                  )}
                />
                <FieldError message={errors.currency?.message} />
              </div>
            </div>
          )}
        </section>

        <section className="flex flex-col gap-2">
          <h3 className="text-lg font-semibold">Hosts</h3>
          <Controller
            name="hostIds"
            control={control}
            render={({ field }) => (
              <HostMultiSelect value={field.value} onChange={field.onChange} />
            )}
          />
          <FieldError message={errors.hostIds?.message} />
        </section>

        <div className="flex gap-2">
          <Button type="submit" disabled={createEvent.isPending}>
            {createEvent.isPending ? "Creating…" : "Create event"}
          </Button>
          <Button type="button" variant="ghost" onClick={() => navigate("/events")}>
            Cancel
          </Button>
        </div>
      </form>
    </AdminLayout>
  );
}
