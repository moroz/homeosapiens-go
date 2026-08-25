import React from "react";
import { useFormContext } from "react-hook-form";

import { InputField, InputGroup } from "~/components/forms";
import { Select } from "~/components/forms/select";
import { SALUTATION_LABELS, type HostFormValues } from "./interfaces";

/**
 * Salutations are i18n message keys rather than free text, so the public site
 * can render them in the visitor's language. Only the keys defined in
 * `i18n/*.json` are offered.
 */
const SalutationOptions = [
  { value: "", label: "(none)" },
  ...Object.entries(SALUTATION_LABELS).map(([value, label]) => ({ value, label })),
];

export const FormFields: React.FC = () => {
  const {
    formState: { errors },
    register,
    control,
  } = useFormContext<HostFormValues>();

  return (
    <>
      <InputGroup className="max-w-2xl">
        <InputField
          label="Given name"
          autoFocus
          errors={errors}
          {...register("givenName", { required: "Required" })}
        />
        <InputField
          label="Family name"
          errors={errors}
          {...register("familyName", { required: "Required" })}
        />
      </InputGroup>

      <InputGroup className="max-w-2xl">
        <Select
          label="Salutation"
          options={SalutationOptions}
          control={control}
          name="salutation"
          errors={errors}
        />
        <InputField
          label="Country"
          placeholder="ISO code, e.g. PL"
          errors={errors}
          {...register("country")}
        />
      </InputGroup>
    </>
  );
};
