import React from "react";
import { Controller, useFormContext } from "react-hook-form";

import { InputField, InputGroup } from "~/components/forms";
import { Select } from "~/components/forms/select";
import { Field } from "~/components/ui/field";
import { Label } from "~/components/ui/label";
import { Switch } from "~/components/ui/switch";
import { Textarea } from "~/components/ui/textarea";
import { useListHostsQuery } from "~/hooks";
import type { VideoFormValues } from "./interfaces";

export const FormFields: React.FC = () => {
  const {
    formState: { errors },
    register,
    control,
  } = useFormContext<VideoFormValues>();
  const { data: hosts } = useListHostsQuery();

  // An empty value clears the host, so it is offered as a first option rather
  // than leaving the select stuck on whoever happens to come first.
  const hostOptions = [
    { value: "", label: "(none)" },
    ...(hosts?.data ?? []).map((host) => ({
      value: host.id,
      label: `${host.givenName} ${host.familyName}`,
    })),
  ];

  return (
    <>
      <InputGroup className="max-w-2xl">
        <InputField
          label="Title (EN)"
          autoFocus
          errors={errors}
          {...register("titleEn", { required: "Required" })}
        />
        <InputField
          label="Title (PL)"
          errors={errors}
          {...register("titlePl", { required: "Required" })}
        />
      </InputGroup>

      <InputGroup className="max-w-2xl">
        <InputField
          label="Slug"
          className="font-mono"
          errors={errors}
          {...register("slug", { required: "Required" })}
        />
        <InputField label="Recorded on" type="date" errors={errors} {...register("recordedOn")} />
      </InputGroup>

      <InputGroup className="max-w-2xl">
        <Select
          label="Host"
          options={hostOptions}
          control={control}
          name="hostId"
          errors={errors}
        />
      </InputGroup>

      <div className="flex items-center gap-2.5">
        <Controller
          name="isPublic"
          control={control}
          render={({ field }) => (
            <Switch id="isPublic" checked={field.value} onCheckedChange={field.onChange} />
          )}
        />
        <Label htmlFor="isPublic">This video is watchable outside a paid playlist</Label>
      </div>

      <div className="flex max-w-4xl gap-6">
        <Field className="flex-1">
          <Label htmlFor="descriptionEn">Description (EN)</Label>
          <Textarea id="descriptionEn" {...register("descriptionEn")} className="h-48" />
        </Field>
        <Field className="flex-1">
          <Label htmlFor="descriptionPl">Description (PL)</Label>
          <Textarea id="descriptionPl" {...register("descriptionPl")} className="h-48" />
        </Field>
      </div>
    </>
  );
};
