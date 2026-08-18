import React from "react";
import { useFormContext } from "react-hook-form";
import type { BlogPostFormValues } from "./interfaces";
import { InputField, InputGroup } from "~/components/forms";
import { Select } from "~/components/forms/select";
import { Textarea } from "~/components/ui/textarea";
import Markdown from "react-markdown";
import { Field } from "~/components/ui/field";
import { Label } from "~/components/ui/label";

interface Props {
  onTitleBlur?: React.ChangeEventHandler<HTMLInputElement>;
}

const LanguageOptions = [
  { value: "en", label: "English" },
  { value: "pl", label: "Polish" },
];

export const FormFields: React.FC<Props> = ({ onTitleBlur }) => {
  const {
    formState: { errors },
    watch,
    register,
    control,
  } = useFormContext<BlogPostFormValues>();

  return (
    <>
      <InputGroup className="max-w-2xl">
        <InputField
          label="Title"
          autoFocus
          errors={errors}
          {...register("title", { required: "Required", onBlur: onTitleBlur })}
        />
        <InputField
          label="Slug"
          className="font-mono"
          errors={errors}
          {...register("slug", { required: "Required" })}
        />
      </InputGroup>

      <InputGroup className="max-w-2xl">
        <Select label="Language" options={LanguageOptions} control={control} name="language" />
      </InputGroup>

      <div className="flex max-w-4xl grid-cols-2 gap-6">
        <Field className="flex-1 ">
          <Label htmlFor="body">Body</Label>
          <Textarea id="body" {...register("body")} className="h-80 font-mono" />
        </Field>
        <Field className="flex h-full flex-1 flex-col">
          <Label>Preview</Label>
          <div className="prose h-80 overflow-y-auto border p-4 outline dark:prose-invert">
            <Markdown>{watch("body")}</Markdown>
          </div>
        </Field>
      </div>
    </>
  );
};
