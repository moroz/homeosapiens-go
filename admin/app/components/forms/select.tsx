import React from "react";
import { Label } from "../ui/label";
import { FieldError } from "./field-error";
import {
  Controller,
  type Control,
  type FieldErrors,
  type FieldValues,
  type Path,
} from "react-hook-form";
import {
  Select as BaseSelect,
  SelectTrigger,
  SelectValue,
  SelectItem,
  SelectContent,
} from "~/components/ui/select";
import { NativeSelect, NativeSelectOption } from "~/components/ui/native-select";

interface SelectOption {
  value: string;
  label: string;
}

interface Props<T extends FieldValues> {
  options: SelectOption[];
  label: string;
  name: Path<T>;
  id?: string;
  control: Control<T>;
  errors?: FieldErrors;
}

export function Select<T extends FieldValues>({
  name,
  id = name,
  options,
  label,
  control,
  errors,
}: Props<T>) {
  const error = errors?.[id]?.message as string;

  return (
    <div className="flex flex-col gap-3">
      <Label htmlFor={id} className="leading-snug">
        {label}
      </Label>
      <Controller
        name={name}
        control={control}
        render={({ field }) => (
          <NativeSelect
            value={field.value}
            onChange={field.onChange}
            onBlur={field.onBlur}
            className="w-full"
          >
            {options.map((option) => (
              <NativeSelectOption key={option.value} value={option.value}>
                {option.label}
              </NativeSelectOption>
            ))}
          </NativeSelect>
        )}
      />
      <FieldError message={error} />
    </div>
  );
}
