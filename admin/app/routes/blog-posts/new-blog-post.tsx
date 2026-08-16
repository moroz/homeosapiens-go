import React, { useCallback, useState } from "react";
import { AdminLayout } from "~/components/admin-layout";
import { type Path, useForm } from "react-hook-form";
import { type CreateBlogPostInput, useCreateBlogPostMutation } from "~/hooks";
import { ApiError, isValidationErrorBody } from "~/lib/api";
import { PageTitle } from "~/components/page-title";
import { Notification } from "~/components/notification";
import { Button } from "~/components/ui/button";
import { InputField, InputGroup } from "~/components/forms";
import { slugify } from "~/lib/slugify";
import { useNavigate } from "react-router";

interface Props {}

export const NewBlogPost: React.FC<Props> = () => {
  const form = useForm<CreateBlogPostInput>();
  const mutation = useCreateBlogPostMutation();
  const [formError, setFormError] = useState<string | null>(null);
  const navigate = useNavigate();

  const {
    formState: { errors },
    register,
    handleSubmit,
    setError,
    setValue,
    getValues,
  } = form;

  const onSubmit = useCallback(
    async (params: CreateBlogPostInput) => {
      setFormError(null);
      try {
        const blogPost = await mutation.mutateAsync(params);
        navigate(`/blog-posts/${blogPost.id}`);
      } catch (err) {
        if (err instanceof ApiError && err.status === 422 && isValidationErrorBody(err.body)) {
          for (const [field, message] of Object.entries(err.body.errors)) {
            setError(field as Path<CreateBlogPostInput>, { message });
          }
        }
      }
    },
    [mutation],
  );

  const onTitleBlur: React.ChangeEventHandler<HTMLInputElement> = useCallback(
    (e) => {
      const slug = getValues("slug");
      if (slug || !e.currentTarget.value) return;
      setValue("slug", slugify(e.currentTarget.value));
    },
    [getValues, setValue],
  );

  return (
    <AdminLayout title="New blog post">
      <form onSubmit={handleSubmit(onSubmit)} className="flex max-w-2xl flex-col gap-6">
        <PageTitle>Create a blog post</PageTitle>
        {formError ? (
          <Notification
            title="An error has prevented this blog post from being saved."
            variant="destructive"
          >
            {formError}
          </Notification>
        ) : null}

        <InputGroup>
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

        <div className="flex gap-2">
          <Button type="submit" disabled={mutation.isPending}>
            {mutation.isPending ? "Creating…" : "Create blog post"}
          </Button>
        </div>
      </form>
    </AdminLayout>
  );
};
