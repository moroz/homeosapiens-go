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
import { Select } from "~/components/forms/select";
import { Textarea } from "~/components/ui/textarea";
import Markdown from "react-markdown";
import { Field } from "~/components/ui/field";
import { Label } from "~/components/ui/label";
import { BackButton } from "~/components/back-button";

interface Props {}

const LanguageOptions = [
  { value: "en", label: "English" },
  { value: "pl", label: "Polish" },
];

export const NewBlogPost: React.FC<Props> = () => {
  const form = useForm<CreateBlogPostInput>({ defaultValues: { language: "en" } });
  const mutation = useCreateBlogPostMutation();
  const [formError, setFormError] = useState<string | null>(null);
  const navigate = useNavigate();

  const {
    formState: { errors },
    watch,
    register,
    handleSubmit,
    setError,
    setValue,
    getValues,
    control,
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
      <BackButton href="/blog-posts">Back to list</BackButton>
      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
        <PageTitle>Create a blog post</PageTitle>
        {formError ? (
          <Notification
            title="An error has prevented this blog post from being saved."
            variant="destructive"
          >
            {formError}
          </Notification>
        ) : null}

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

        <div className="flex gap-2">
          <Button type="submit" disabled={mutation.isPending}>
            {mutation.isPending ? "Creating…" : "Create blog post"}
          </Button>
        </div>
      </form>
    </AdminLayout>
  );
};
