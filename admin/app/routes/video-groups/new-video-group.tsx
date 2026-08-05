import { useCallback, useState } from "react";
import { FormProvider, type Path, useForm } from "react-hook-form";
import { useNavigate } from "react-router";

import { AdminLayout } from "~/components/admin-layout";
import { Notification } from "~/components/notification";
import { PageTitle } from "~/components/page-title";
import { Button } from "~/components/ui/button";
import { useCreateVideoGroupMutation } from "~/hooks";
import { ApiError, isValidationErrorBody } from "~/lib/api";
import { slugify } from "~/lib/slugify";
import { FormFields } from "./form-fields";
import { type VideoGroupFormValues, toVideoGroupInput } from "./interfaces";

const defaultValues: Partial<VideoGroupFormValues> = {
  isFree: true,
  currency: "PLN",
};

export default function NewVideoGroup() {
  const navigate = useNavigate();
  const createVideoGroup = useCreateVideoGroupMutation();
  const [formError, setFormError] = useState<string | null>(null);

  const form = useForm<VideoGroupFormValues>({ defaultValues });
  const { handleSubmit, setError, setValue, getValues } = form;

  /** Set a default slug after the English title has been set. */
  const onTitleEnBlur = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const slug = getValues("slug");
      if (slug || !e.currentTarget.value) return;

      setValue("slug", slugify(e.currentTarget.value));
    },
    [getValues, setValue],
  );

  async function onSubmit(values: VideoGroupFormValues) {
    setFormError(null);
    try {
      await createVideoGroup.mutateAsync(toVideoGroupInput(values));
      navigate("/videos");
    } catch (err) {
      if (err instanceof ApiError && err.status === 422 && isValidationErrorBody(err.body)) {
        for (const [field, message] of Object.entries(err.body.errors)) {
          setError(field as Path<VideoGroupFormValues>, { message });
        }
        setFormError("Please fix the errors below.");
        return;
      }
      setFormError("Something went wrong creating the series. Please try again.");
    }
  }

  return (
    <AdminLayout title="New video series">
      <FormProvider {...form}>
        <form onSubmit={handleSubmit(onSubmit)} className="flex max-w-2xl flex-col gap-6">
          <PageTitle className="mb-0">Create a video series</PageTitle>
          {formError ? (
            <Notification
              title="An error has prevented this series from being saved."
              variant="destructive"
            >
              {formError}
            </Notification>
          ) : null}

          <FormFields onTitleEnBlur={onTitleEnBlur} />

          <div className="flex gap-2">
            <Button type="submit" disabled={createVideoGroup.isPending}>
              {createVideoGroup.isPending ? "Creating…" : "Create series"}
            </Button>
            <Button type="button" variant="ghost" onClick={() => navigate("/videos")}>
              Cancel
            </Button>
          </div>
        </form>
      </FormProvider>
    </AdminLayout>
  );
}
