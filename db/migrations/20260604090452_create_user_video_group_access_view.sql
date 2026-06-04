-- +goose Up
create view user_video_group_access as
select u.id user_id, vg.id video_group_id, (vg.product_id is null or u.user_role = 'Administrator' or upa.id is not null) has_access
from users u
cross join video_groups vg
left join user_product_access upa on u.id = upa.user_id and vg.product_id = upa.product_id;

-- +goose Down
drop view user_video_group_access;
