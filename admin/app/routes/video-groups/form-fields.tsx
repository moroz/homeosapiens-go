import React from "react";
import { Controller, useFormContext } from "react-hook-form";

import { InputField, InputGroup, Select } from "~/components/forms";
import { Label } from "~/components/ui/label";
import { Switch } from "~/components/ui/switch";
import type { VideoGroupFormValues } from "./interfaces";
import { VideoPicker } from "./video-picker";

interface Props {
  onTitleEnBlur?: React.ChangeEventHandler<HTMLInputElement>;
  /** The video list is only offered once the group exists, so that its videos have somewhere to be saved. */
  showVideos?: boolean;
}

const CURRENCY_OPTIONS = [
  { value: "PLN", label: "PLN" },
  { value: "EUR", label: "EUR" },
];

export const FormFields: React.FC<Props> = ({ onTitleEnBlur, showVideos = false }) => {
  const {
    register,
    formState: { errors },
    control,
    watch,
  } = useFormContext<VideoGroupFormValues>();

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
          <InputField
            label="Slug"
            containerClassName="col-span-2"
            errors={errors}
            {...register("slug", { required: "Required" })}
          />
        </InputGroup>
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
          <Label htmlFor="isFree">This series is free to watch</Label>
        </div>

        {!isFree && (
          <InputGroup className="rounded-md border border-input p-4">
            <InputField
              label="Price"
              placeholder="199.00"
              errors={errors}
              {...register("price", { required: "Required for paid series" })}
            />
            <Select label="Currency" name="currency" control={control} options={CURRENCY_OPTIONS} />
          </InputGroup>
        )}
      </section>

      {showVideos && (
        <section className="flex flex-col gap-4">
          <h3 className="text-lg font-semibold">Videos</h3>
          <Controller
            name="videoIds"
            control={control}
            render={({ field }) => (
              <VideoPicker value={field.value ?? []} onChange={field.onChange} />
            )}
          />
        </section>
      )}
    </div>
  );
};
