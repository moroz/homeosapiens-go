import React, { useState } from "react";
import { useNavigate, useParams } from "react-router";
import Markdown from "react-markdown";

import { Button } from "~/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "~/components/ui/dialog";
import { Textarea } from "~/components/ui/textarea";
import { useGetEventQuery, useUpdateEventMutation } from "~/hooks";

/** `:locale` route param to the event field it edits. */
const FIELD_BY_LOCALE = { en: "descriptionEn", pl: "descriptionPl" } as const;
const LABEL_BY_LOCALE = { en: "English", pl: "Polish" } as const;

type Locale = keyof typeof FIELD_BY_LOCALE;

function isLocale(value: string | undefined): value is Locale {
  return value != null && value in FIELD_BY_LOCALE;
}

/**
 * Modal route (`events/:id/description/:locale`) editing one description of an
 * event in Markdown, opened from the detail view. Saving PATCHes just that
 * field; closing discards the draft.
 */
export const MarkdownEditorDialog: React.FC = () => {
  const { id, locale } = useParams();
  const { data: event, isPending } = useGetEventQuery(id);

  if (!isLocale(locale) || !id) return null;

  return (
    <Editor
      key={`${id}:${locale}`}
      id={id}
      locale={locale}
      initialValue={event?.[FIELD_BY_LOCALE[locale]] ?? ""}
      isPending={isPending}
    />
  );
};

interface EditorProps {
  id: string;
  locale: Locale;
  initialValue: string;
  isPending: boolean;
}

/**
 * Split out so `initialValue` seeds `useState` once the event has loaded — the
 * `key` on the parent remounts this when the loaded value arrives.
 */
const Editor: React.FC<EditorProps> = ({ id, locale, initialValue, isPending }) => {
  const navigate = useNavigate();
  const mutation = useUpdateEventMutation();
  const [draft, setDraft] = useState(initialValue);

  /** `..` resolves against the route hierarchy, i.e. back to the detail view. */
  const close = () => navigate("..");

  async function save() {
    try {
      await mutation.mutateAsync({
        id,
        params: { [FIELD_BY_LOCALE[locale]]: draft.trim() === "" ? null : draft },
      });
      close();
    } catch {
      // Surfaced below via `mutation.isError`.
    }
  }

  return (
    <Dialog open onOpenChange={(open) => !open && close()}>
      <DialogContent className="sm:max-w-7xl">
        <DialogHeader>
          <DialogTitle>Edit description &mdash; {LABEL_BY_LOCALE[locale]}</DialogTitle>
          <DialogDescription>
            Markdown. The preview approximates the published page.
          </DialogDescription>
        </DialogHeader>
        {isPending ? (
          <p className="text-muted-foreground">Loading&hellip;</p>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2">
            <Textarea
              aria-label="Markdown source"
              className="h-120 font-mono text-sm"
              value={draft}
              onChange={(event) => setDraft(event.target.value)}
              autoFocus
            />
            <div className="prose max-h-120 max-w-none overflow-y-auto rounded-md border p-4 text-foreground dark:prose-invert">
              <Markdown>{draft}</Markdown>
            </div>
          </div>
        )}
        <DialogFooter>
          {mutation.isError && (
            <p className="mr-auto self-center text-destructive">Could not save. Try again.</p>
          )}
          <Button type="button" variant="outline" onClick={close}>
            Cancel
          </Button>
          <Button type="button" onClick={save} disabled={isPending || mutation.isPending}>
            Save
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
