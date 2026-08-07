import { type DialogRoot } from "@base-ui/react";
import { Dialog, DialogTitle, DialogContent, DialogHeader } from "~/components/ui/dialog";
import React, { useCallback } from "react";
import { useNavigate, useParams } from "react-router";
import { useGetEventQuery } from "~/hooks";

interface Props {}

export const AddAttendantDialog: React.FC<Props> = () => {
  const navigate = useNavigate();

  const { id } = useParams();
  const { data: event, isPending } = useGetEventQuery(id!);

  const close = useCallback(() => navigate(".."), [navigate]);

  const onOpenChange = useCallback(
    (open: boolean, _eventDetails: DialogRoot.ChangeEventDetails) => {
      if (!open) close();
    },
    [navigate],
  );

  if (isPending) return null;

  return (
    <Dialog open onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Enroll student for event &ldquo;{event?.titleEn}&rdquo;</DialogTitle>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
};
