import React from "react";
import { Controller, useFormContext } from "react-hook-form";
import type { EventFormValues } from "./interfaces";
import { InputGroup, InputField, FieldError } from "~/components/forms";
import { Label } from "~/components/ui/label";
import { Select } from "~/components/forms";
import { Switch } from "~/components/ui/switch";
import { HostMultiSelect } from "./host-multi-select";

interface Props {
  onTitleEnBlur?: React.ChangeEventHandler<HTMLInputElement>;
}

const CURRENCY_OPTIONS = [
  { value: "PLN", label: "PLN" },
  { value: "EUR", label: "EUR" },
];

const EVENT_TYPE_OPTIONS = [
  { value: "webinar", label: "Webinar" },
  { value: "seminar", label: "Seminar" },
];

export const FormFields: React.FC<Props> = ({ onTitleEnBlur }) => {
  const form = useFormContext<EventFormValues>();
  const {
    register,
    formState: { errors },
    control,
    watch,
  } = form;

  const isVirtual = watch("isVirtual");
  const isFree = watch("isFree");

  return (
    <div className="flex max-w-2xl flex-col gap-6">
      <section className="flex flex-col gap-4">
        <h3 className="text-lg font-semibold">Basics</h3>

        <InputGroup>
          <InputField
            label="Title (EN)"
            errors={errors}
            {...register("titleEn", { required: "Required", onBlur: onTitleEnBlur })}
          />
          <InputField
            label="Title (PL)"
            errors={errors}
            {...register("titlePl", { required: "Required" })}
          />
          <InputField label="Subtitle (EN)" errors={errors} {...register("subtitleEn")} />
          <InputField label="Subtitle (PL)" errors={errors} {...register("subtitlePl")} />
        </InputGroup>

        <InputGroup>
          <InputField
            label="Slug"
            className="font-mono"
            placeholder="my-event-slug"
            errors={errors}
            {...register("slug", { required: "Required" })}
          />
          <Select
            name="eventType"
            label="Event type"
            errors={errors}
            options={EVENT_TYPE_OPTIONS}
            control={control}
          />
        </InputGroup>
      </section>

      <section className="flex flex-col gap-4">
        <header>
          <h3 className="text-lg font-semibold">Schedule</h3>
          <p className="text-sm text-muted-foreground">
            Dates and times are shown in Warsaw time zone.
          </p>
        </header>
        <InputGroup>
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
        </InputGroup>
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

        {isVirtual && (
          <InputGroup className="rounded-md border border-input p-4">
            <InputField
              label="Meeting link (Zoom)"
              containerClassName="col-span-2"
              placeholder="https://zoom.us/j/…"
              errors={errors}
              {...register("meetingUrl")}
            />
          </InputGroup>
        )}

        {!isVirtual && (
          <InputGroup className="rounded-md border border-input p-4">
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
          </InputGroup>
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
          <InputGroup className="rounded-md border border-input p-4">
            <InputField
              label="Price"
              placeholder="19.99"
              errors={errors}
              {...register("price", { required: "Required for paid events" })}
            />
            <Select label="Currency" name="currency" control={control} options={CURRENCY_OPTIONS} />
          </InputGroup>
        )}
      </section>

      <section className="flex flex-col gap-2">
        <h3 className="text-lg font-semibold">Hosts</h3>
        <Controller
          name="hostIds"
          control={control}
          render={({ field }) => <HostMultiSelect value={field.value} onChange={field.onChange} />}
        />
        <FieldError message={errors.hostIds?.message} />
      </section>
    </div>
  );
};
