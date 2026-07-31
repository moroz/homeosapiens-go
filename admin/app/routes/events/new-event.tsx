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
import { useCreateEventMutation, useListHostsQuery } from "~/hooks";
import { InputField } from "~/components/forms/input-field";
import { FieldError } from "~/components/forms/field-error";
import { type EventFormValues, toEventInput } from "./interfaces";

const defaultValues: Partial<EventFormValues> = {
  eventType: "webinar",
  isVirtual: true,
  isFree: true,
  currency: "PLN",
  hostIds: [],
};

/** Server field names line up 1:1 with `FormValues` keys, so validation errors map straight onto form fields. */
const FORM_FIELDS = new Set<string>(Object.keys(defaultValues));

function isFormField(field: string): field is keyof EventFormValues {
  return FORM_FIELDS.has(field);
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
  } = useForm<EventFormValues>({ defaultValues });

  const isVirtual = watch("isVirtual");
  const isFree = watch("isFree");

  async function onSubmit(values: EventFormValues) {
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
            <InputField
              label="Title (EN)"
              errors={errors}
              {...register("titleEn", { required: "Required" })}
            />
            <InputField
              label="Title (PL)"
              errors={errors}
              {...register("titlePl", { required: "Required" })}
            />
            <InputField label="Subtitle (EN)" errors={errors} {...register("subtitleEn")} />
            <InputField label="Subtitle (PL)" errors={errors} {...register("subtitlePl")} />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <InputField
              label="Slug"
              className="font-mono"
              placeholder="my-event-slug"
              errors={errors}
              {...register("slug", { required: "Required" })}
            />
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
            <InputField
              label="Starts at"
              type="datetime-local"
              errors={errors}
              {...register("startsAt", { required: "Required" })}
            />
            <InputField
              label="Ends at"
              type="datetime-local"
              errors={errors}
              {...register("endsAt", { required: "Required" })}
            />
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
              <InputField
                label="Venue name (EN)"
                errors={errors}
                {...register("venueNameEn", { required: "Required for in-person events" })}
              />
              <InputField
                label="Venue name (PL)"
                errors={errors}
                {...register("venueNamePl", { required: "Required for in-person events" })}
              />
              <InputField
                label="Street"
                containerClassName="col-span-2"
                errors={errors}
                {...register("venueStreet", { required: "Required for in-person events" })}
              />
              <InputField
                label="City (EN)"
                errors={errors}
                {...register("venueCityEn", { required: "Required for in-person events" })}
              />
              <InputField
                label="City (PL)"
                errors={errors}
                {...register("venueCityPl", { required: "Required for in-person events" })}
              />
              <InputField
                label="Postal code"
                errors={errors}
                {...register("venuePostalCode", { required: "Required for in-person events" })}
              />
              <InputField
                label="Country code"
                placeholder="PL"
                maxLength={2}
                errors={errors}
                {...register("venueCountryCode", {
                  required: "Required for in-person events",
                  minLength: { value: 2, message: "Must be 2 letters" },
                  maxLength: { value: 2, message: "Must be 2 letters" },
                })}
              />
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
              <InputField
                label="Price"
                placeholder="19.99"
                errors={errors}
                {...register("price", { required: "Required for paid events" })}
              />
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
