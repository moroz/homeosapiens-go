import React from "react";
import { Controller, useFormContext } from "react-hook-form";
import type { EventFormValues } from "./interfaces";
import { InputGroup, InputField, FieldError } from "~/components/forms";
import { Label } from "~/components/ui/label";
import { Textarea } from "~/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";

interface Props {}

export const FormFields: React.FC<Props> = () => {
  const form = useFormContext<EventFormValues>();
  const {
    register,
    formState: { errors },
    control,
  } = form;

  return (
    <>
      <section className="flex flex-col gap-4">
        <h3 className="text-lg font-semibold">Basics</h3>

        <InputGroup>
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
        </InputGroup>

        <InputGroup>
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
        </InputGroup>

        <InputGroup>
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="descriptionEn">Description (EN)</Label>
            <Textarea id="descriptionEn" {...register("descriptionEn", { required: "Required" })} />
            <FieldError message={errors.descriptionEn?.message} />
          </div>
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="descriptionPl">Description (PL)</Label>
            <Textarea id="descriptionPl" {...register("descriptionPl", { required: "Required" })} />
            <FieldError message={errors.descriptionPl?.message} />
          </div>
        </InputGroup>
      </section>
    </>
  );
};
